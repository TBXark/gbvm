package cmd

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/TBXark/gbvm/internal/bin"
	"github.com/TBXark/gbvm/internal/command"
	"github.com/TBXark/gbvm/internal/env"
)

func NewInstallCommand() *command.Command {
	fs := flag.NewFlagSet("install", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Printf("Usage: gbvm install [options] <backup file>\n\n")
		fmt.Printf("Install Go binaries from backup file\n\n")
		fs.PrintDefaults()
	}
	return command.NewCommand(fs, func() error {
		if fs.NArg() == 0 {
			return fmt.Errorf("missing backup file")
		}
		return handleInstall(fs.Arg(0))
	})
}

func handleInstall(backupPath string) error {
	file, err := os.ReadFile(backupPath)
	if err != nil {
		return err
	}
	var versions []*bin.Info
	err = json.Unmarshal(file, &versions)
	if err != nil {
		return err
	}
	for _, v := range versions {
		binPath := filepath.Join(env.GoPath, "bin", v.Name)
		info, e := bin.LoadBinInfo(binPath)
		if e != nil {
			if !errors.Is(e, os.ErrNotExist) {
				fmt.Printf("failed to load %s: %v\n", v.Name, e)
				continue
			}
		} else {
			if info.Version == v.Version {
				fmt.Printf("skip %s\n", v.Name)
				continue
			}
		}
		fmt.Printf("installing %s@%s\n", v.Name, v.Version)
		e = bin.InstallBinByVersion(v.Path, v.Version)
		if e != nil {
			fmt.Printf("failed to install %s: %v\n", v.Name, e)
		}
	}
	return nil
}
