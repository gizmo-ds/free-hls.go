package client

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	sp "service-provider"

	"free-hls.go/utils"
	"github.com/grafov/m3u8"
)

type (
	UploadI interface {
		Upload([]byte) (string, error)
		SetProxy(string)
		Info() (string, string, int64)
	}
)

func uploadFile(i UploadI, data []byte) (string, error) {
	i.SetProxy(utils.Config.Client.Proxy.Upload)
	dataUrl, err := i.Upload(data)
	if err != nil {
		return "", err
	}
	return dataUrl, nil
}

func UploadHLS(title, spName string) (id string, err error) {
	listFileName := filepath.Join(utils.AppDir, "./tmp/out.m3u8")
	f, err := os.Open(listFileName)
	if err != nil {
		panic(err)
	}
	p, listType, err := m3u8.DecodeFrom(bufio.NewReader(f), true)
	if err != nil {
		panic(err)
	}
	switch listType {
	case m3u8.MEDIA:
		mediaPlaylist := p.(*m3u8.MediaPlaylist)
		for _, v := range mediaPlaylist.Segments {
			if v != nil {
				localTsName := filepath.Join(filepath.Dir(listFileName), v.URI)
				tsData, err := ioutil.ReadFile(localTsName)
				if err != nil {
					log.Println(err)
					return "", err
				}
				var tsUrl string
				upload, err := sp.Load(spName)
				if err != nil {
					log.Println(err)
					return "", err
				}
				name, _, maxSize := upload.Info()
				if int64(len(tsData)) < maxSize {
					tsUrl, err = uploadFile(upload, tsData)
				} else {
					fmt.Printf("文件(%s)过大, 服务(%s)无法上传", v.URI, name)
				}
				fmt.Println(name, v.URI, "-->", tsUrl)
				v.URI = tsUrl
			}
		}

		id, err = pushToServer(p.Encode().Bytes(), title)
		if err != nil {
			panic(err)
		}
	default:
		return
	}
	return
}

func pushToServer(data []byte, title string) (id string, err error) {
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)

	if fw, err1 := w.CreateFormFile("file", "data.m3u8"); err1 == nil && fw != nil {
		if _, err = fw.Write(data); err != nil {
			return
		}
	}
	if fw, err1 := w.CreateFormField("title"); err1 == nil && fw != nil {
		if _, err = fw.Write([]byte(title)); err != nil {
			return
		}
	}
	_ = w.Close()

	cl := http.Client{
		Timeout: time.Second * 10,
	}
	if proxyUrl := utils.Config.Client.Proxy.Push; proxyUrl != "" {
		proxy, err := url.Parse(proxyUrl)
		if err == nil {
			cl.Transport = &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				Proxy:           http.ProxyURL(proxy),
			}
		}
	}
	_url, err := url.Parse(utils.Config.Client.Server)
	if err != nil {
		return
	}
	_url.Path = "/file"
	req, err := http.NewRequest("POST", _url.String(), buf)
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", w.FormDataContentType())
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", utils.Config.Client.UploadKey))

	resp, err := cl.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("StatusCode: %v", resp.StatusCode))
		return
	}

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var rj struct {
		Id string `json:"id"`
	}
	if err = json.Unmarshal(data, &rj); err != nil {
		return
	}
	id = rj.Id
	return
}
