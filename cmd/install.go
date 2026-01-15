package cmd

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/TBXark/gbvm/internal/bin"
	"github.com/TBXark/gbvm/internal/command"
	"github.com/TBXark/gbvm/internal/env"
	"github.com/TBXark/gbvm/internal/log"
)

func NewInstallCommand() *command.Command {
	fs := flag.NewFlagSet("install", flag.ExitOnError)
	fs.Usage = func() {
		command.PrintUsage(fs, "gbvm install [options] <backup file>", "Install Go binaries from backup file")
	}
	cmd := command.NewCommand(fs, func() error {
		if fs.NArg() == 0 {
			return fmt.Errorf("missing backup file")
		}
		return handleInstall(fs.Arg(0))
	})
	cmd.Usage = "Install Go binaries from backup file"
	return cmd
}

func handleInstall(backupPath string) error {
	file, err := os.ReadFile(backupPath)
	if err != nil {
		return err
	}
	versions, err := parseBackupFile(file)
	if err != nil {
		return err
	}
	for _, v := range versions {
		info, loadErr := bin.FindBinInfo(v.Name)
		decision := decideInstall(v, info, loadErr)
		if decision.message != "" {
			log.Debug(decision.message)
		}
		if decision.shouldInstall {
			if insErr := bin.InstallBinByVersion(v.Path, v.Version, env.GoProxy); insErr != nil {
				log.Errorf("failed to install %s: %v", v.Name, insErr)
			}
		}
	}
	return nil
}

func parseBackupFile(data []byte) ([]*bin.Info, error) {
	var versions []*bin.Info
	if err := json.Unmarshal(data, &versions); err != nil {
		return nil, err
	}
	return versions, nil
}

type installDecision struct {
	shouldInstall bool
	message       string
}

func decideInstall(desired *bin.Info, existing *bin.Info, err error) installDecision {
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return installDecision{
			shouldInstall: false,
			message:       fmt.Sprintf("failed to load %s: %v", desired.Name, err),
		}
	}
	if err == nil && existing.Version == desired.Version {
		return installDecision{
			shouldInstall: false,
			message:       fmt.Sprintf("skip %s", desired.Name),
		}
	}
	return installDecision{
		shouldInstall: true,
		message:       fmt.Sprintf("installing %s@%s", desired.Name, desired.Version),
	}
}
