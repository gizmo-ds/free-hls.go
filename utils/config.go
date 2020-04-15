package utils

import (
	"github.com/BurntSushi/toml"
)

var Config struct {
	Server struct {
		DataDir       string   `yaml:"dataDir"`
		FileSizeLimit int64    `yaml:"fileSizeLimit"`
		UseUploadKey  bool     `yaml:"useUploadKey"`
		UploadKeys    []string `yaml:"uploadKeys"`
	} `yaml:"server"`
	Client struct {
		Server    string `yaml:"server"`
		UploadKey string `yaml:"uploadKey"`
		TmpDir    string `yaml:"tmpDir"`
		Proxy     struct {
			Upload string `yaml:"cdn"`
			Push   string `yaml:"server"`
		} `yaml:"proxy"`
		FFmpegRoot string `yaml:"ffmpegRoot"`
	} `yaml:"client"`
}

func LoadConfig(filename string) {
	if _, err := toml.DecodeFile(filename, &Config); err != nil {
		panic(err)
	}
}
