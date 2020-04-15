package server

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/xid"

	"free-hls.go/utils"

	"github.com/labstack/echo/v4"
)

func Start(addr string) {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	e.Static("/", filepath.Join(utils.AppDir, "server/web"))
	e.Renderer = &Template{
		templates: template.Must(template.ParseGlob(filepath.Join(utils.AppDir, "server/web/*.html"))),
	}

	e.GET("/file/:id", func(c echo.Context) error {
		fileId := c.Param("id")
		data, err := ioutil.ReadFile(filepath.Join(utils.Config.Server.DataDir, fileId))
		if err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"message": "file not found"})
		}
		var fileInfo DataInfo
		if err = json.Unmarshal(data, &fileInfo); err != nil {
			return err
		}
		data, err = base64.StdEncoding.DecodeString(fileInfo.Data)
		if err != nil {
			return err
		}
		return c.Blob(http.StatusOK, "application/vnd.apple.mpegurl", data)
	})
	e.GET("/info/:id", func(c echo.Context) error {
		fileId := c.Param("id")
		data, err := ioutil.ReadFile(filepath.Join(utils.Config.Server.DataDir, fileId))
		if err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"message": "file not found"})
		}
		var fileInfo DataInfo
		if err = json.Unmarshal(data, &fileInfo); err != nil {
			return err
		}
		_id, _ := xid.FromString(fileId)
		fileInfo.CreatedAt = _id.Time().Unix()
		return c.JSON(http.StatusOK, fileInfo)
	})
	e.POST("/file", func(c echo.Context) error {
		uploadFile, err := c.FormFile("file")
		if err != nil {
			return err
		}
		if utils.Config.Server.FileSizeLimit > 0 &&
			uploadFile.Size > utils.Config.Server.FileSizeLimit {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": "Bad Request"})
		}
		f, err := uploadFile.Open()
		if err != nil {
			return err
		}
		defer f.Close()

		data, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}
		dataB64 := base64.StdEncoding.EncodeToString(data)
		metaMap := echo.Map{}
		title := c.FormValue("title")
		if title != "" {
			metaMap["title"] = title
		}
		saveFileData, _ := json.Marshal(Data{
			UserKey:     fmt.Sprint(c.Get("userKey")),
			Data:        dataB64,
			ContentType: "application/vnd.apple.mpegurl",
			Meta:        metaMap,
		})
		id := xid.New().String()
		if err = ioutil.WriteFile(filepath.Join(utils.Config.Server.DataDir, id), saveFileData, 0666); err != nil {
			return err
		}
		return c.JSON(http.StatusOK, echo.Map{"message": "success", "id": id})
	}, func() []echo.MiddlewareFunc {
		if utils.Config.Server.UseUploadKey {
			return []echo.MiddlewareFunc{
				middleware.KeyAuth(func(key string, c echo.Context) (bool, error) {
					for _, userKey := range utils.Config.Server.UploadKeys {
						if key == userKey {
							c.Set("userKey", userKey)
							return true, nil
						}
					}
					return false, nil
				}),
			}
		}
		return nil
	}()...)
	e.GET("/play/:id", func(c echo.Context) error {
		fileId := c.Param("id")
		data, err := ioutil.ReadFile(filepath.Join(utils.Config.Server.DataDir, fileId))
		if err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"message": "file not found"})
		}
		var fileInfo Data
		if err = json.Unmarshal(data, &fileInfo); err != nil {
			return err
		}
		return c.Render(http.StatusOK, "player.html", echo.Map{
			"title": fileInfo.Meta["title"],
			"data":  fileInfo.Data,
		})
	})
	e.Logger.Fatal(e.Start(addr))
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
