package bin

import (
	"debug/buildinfo"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/TBXark/gbvm/internal/env"
)

type Info struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Mod     string `json:"mod"`
	Path    string `json:"path"`
}

func LoadBinInfo(binPath string) (*Info, error) {
	info, err := buildinfo.ReadFile(binPath)
	if err != nil {
		return nil, err
	}
	fineName := filepath.Base(binPath)
	return &Info{
		Name:    fineName,
		Version: info.Main.Version,
		Mod:     info.Main.Path,
		Path:    info.Path,
	}, nil
}

func LoadAllBinVersions() ([]*Info, error) {
	binPath := filepath.Join(env.GoPath, "bin")
	files, err := os.ReadDir(binPath)
	if err != nil {
		return nil, err
	}
	var result []*Info
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fullPath := filepath.Join(binPath, file.Name())
		info, e := LoadBinInfo(fullPath)
		if e != nil {
			continue
		}
		result = append(result, info)
	}
	return result, nil
}

func InstallBinByVersion(cmdPath, version string) error {
	uri := fmt.Sprintf("%s@%s", cmdPath, version)
	cmd := exec.Command("go", "install", uri)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
