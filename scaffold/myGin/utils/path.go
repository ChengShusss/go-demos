package utils

import (
	"os"
	"path"
)

var exePath string
var workDir string

func init() {
	var err error
	exePath, err = os.Executable()
	if err != nil {
		panic(err)
	}
	workDir = path.Dir(exePath)
}

func GetExePath() string {
	return exePath
}

func GetExeDir() string {
	return workDir
}
