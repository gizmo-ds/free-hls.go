package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"free-hls.go/utils"
	"github.com/spf13/cobra"
)

var (
	configFile string
)

var (
	rootCmd = &cobra.Command{
		Use:   "freeHLS",
		Short: "一个免费的 HLS 解决方案",
	}
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of FreeHLS.go",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf(
				"FreeHLS.go %s\nBuild Commit: %s\nBuild Date: %s\n",
				Version, Commit, BuildDate,
			)
		},
	}
)

func Execute(version, commit, buildDate string) {
	Version, Commit, BuildDate = version, commit, buildDate
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "当前目录下config.toml文件", "自定义配置文件路径")
	rootCmd.AddCommand(versionCmd)

}

func initConfig() {
	if configFile != "" && configFile != "当前目录下config.toml文件" {
		utils.LoadConfig(configFile)
	} else {
		configFile = filepath.Join(utils.AppDir, "./config.toml")
		utils.LoadConfig(configFile)
	}
}
