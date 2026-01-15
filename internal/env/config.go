package env

import (
	"os"
	"path/filepath"
	"strings"
)

var (
	GoPaths   []string
	GoProxy   string
	GoProxies []string
)

const OfficialProxy = "https://proxy.golang.org"

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
	paths := filepath.SplitList(goPath)
	for _, path := range paths {
		trimmed := strings.TrimSpace(path)
		if trimmed != "" {
			GoPaths = append(GoPaths, trimmed)
		}
	}
	if len(GoPaths) == 0 {
		GoPaths = []string{goPath}
	}
	goProxy := os.Getenv("GOPROXY")
	if goProxy == "" {
		goProxy = "https://proxy.golang.org"
	}
	GoProxy = goProxy
	GoProxies = SplitGoProxy(goProxy)
}

func SplitGoProxy(proxy string) []string {
	parts := strings.FieldsFunc(proxy, func(r rune) bool { return r == ',' || r == '|' })
	var proxies []string
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			proxies = append(proxies, trimmed)
		}
	}
	if len(proxies) == 0 && proxy != "" {
		return []string{proxy}
	}
	return proxies
}
