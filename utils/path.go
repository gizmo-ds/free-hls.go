package utils

import (
	"os"
	"path/filepath"
)

var AppDir string

func init() {
	var err error
	AppDir, err = os.Executable()
	if err != nil {
		panic(err)
	}
	AppDir = filepath.Dir(AppDir)
}
