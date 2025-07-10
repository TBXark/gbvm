package env

import (
	"os"
	"path/filepath"
)

var (
	GoPath  string
	GoProxy string
)

const DevelVersion = "(devel)"

func init() {
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		home, err := os.UserHomeDir()
		if err == nil {
			goPath = filepath.Join(home, "go")
		} else {
			goPath = "go"
		}
	}
	goProxy := os.Getenv("GOPROXY")
	if goProxy == "" {
		goProxy = "https://proxy.golang.org"
	}
	GoPath = goPath
	GoProxy = goProxy
}
