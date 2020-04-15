package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"free-hls.go/client"
	"free-hls.go/utils"
	"github.com/spf13/cobra"
)

var (
	uploadTitle string
	spName      string
	videoS      int
)

var (
	uploadCmd = &cobra.Command{
		Use:   "upload",
		Short: "上传视频",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("未找到必须提供的文件名")
			}
			uploadFile := args[0]
			f, err := os.OpenFile(uploadFile, os.O_RDONLY, 0666)
			if err != nil {
				return errors.New("文件读取失败, 请检测文件是否存在并可读.")
			}
			_ = f.Close()
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			uploadFile := args[0]
			if !cmd.Flag("title").Changed { // 如果未提供自定义标题, 则使用文件名作为标题
				uploadTitle = filepath.Base(uploadFile)
			}
			if err := client.FF(uploadFile, videoS); err != nil {
				log.Println(err)
				return
			}
			id, err := client.UploadHLS(uploadTitle, spName)
			if err != nil {
				log.Println(err)
				return
			}
			fmt.Printf("在线播放地址: %s/play/%s\n", utils.Config.Client.Server, id)
			fmt.Printf("M3U8清单地址: %s/file/%s\n", utils.Config.Client.Server, id)
		},
	}
)

func init() {
	uploadCmd.Flags().StringVarP(&uploadTitle, "title", "t", "视频文件名", "上传后视频的标题, 仅在线预览有效")
	uploadCmd.Flags().StringVarP(&spName, "provider", "p", "语雀", "服务提供者的名称")
	uploadCmd.Flags().IntVarP(&videoS, "second", "s", 5, "视频分片时间")

	rootCmd.AddCommand(uploadCmd)
}
