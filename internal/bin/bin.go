package bin

import (
	"debug/buildinfo"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/TBXark/gbvm/internal/env"
	"github.com/TBXark/gbvm/internal/log"
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

func LoadAllBinVersions(verbose bool) ([]*Info, error) {
	var result []*Info
	seen := make(map[string]struct{})
	for _, goPath := range env.GoPaths {
		binPath := filepath.Join(goPath, "bin")
		files, err := os.ReadDir(binPath)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if _, exists := seen[file.Name()]; exists {
				continue
			}
			fullPath := filepath.Join(binPath, file.Name())
			info, e := LoadBinInfo(fullPath)
			if e != nil {
				if verbose {
					log.Errorf("skip %s: %v", fullPath, e)
				}
				continue
			}
			seen[file.Name()] = struct{}{}
			result = append(result, info)
		}
	}
	return result, nil
}

func FindBinInfo(binName string) (*Info, error) {
	for _, goPath := range env.GoPaths {
		binPath := filepath.Join(goPath, "bin", binName)
		info, err := LoadBinInfo(binPath)
		if err == nil {
			return info, nil
		}
		if os.IsNotExist(err) {
			continue
		}
		return nil, err
	}
	return nil, os.ErrNotExist
}

func InstallBinByVersion(cmdPath, version, proxy string) error {
	uri := fmt.Sprintf("%s@%s", cmdPath, version)
	cmd := exec.Command("go", "install", uri)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if proxy != "" {
		cmd.Env = append(os.Environ(), fmt.Sprintf("GOPROXY=%s", proxy))
	}
	return cmd.Run()
}
