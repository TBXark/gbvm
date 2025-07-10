package cmd

import (
	"flag"
	"fmt"
	"path/filepath"

	"github.com/TBXark/gbvm/internal/bin"
	"github.com/TBXark/gbvm/internal/command"
	"github.com/TBXark/gbvm/internal/env"
	"github.com/TBXark/gbvm/internal/version"
)

func NewUpgradeCommand() *command.Command {
	fs := flag.NewFlagSet("upgrade", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Printf("Usage: gbvm upgrade [options] [bin1 bin2 ...]\n\n")
		fmt.Printf("Upgrade Go binaries\n\n")
		fs.PrintDefaults()
	}
	skipDev := fs.Bool("skip-dev", false, "skip dev version")
	return command.NewCommand(fs, func() error {
		if fs.NArg() == 0 {
			return upgradeAllBins(*skipDev)
		} else {
			for _, binName := range fs.Args() {
				if err := upgradeBin(binName); err != nil {
					return err
				}
			}
			return nil
		}
	})
}

func upgradeAllBins(skipDev bool) error {
	versions, err := bin.LoadAllBinVersions()
	if err != nil {
		return err
	}
	for _, v := range versions {
		if skipDev && v.Version == env.DevelVersion {
			continue
		}
		if e := tryUpgradeBin(v); e != nil {
			fmt.Printf("failed to upgrade %s: %v\n", v.Name, e)
		}
	}
	return nil
}

func upgradeBin(binName string) error {
	binPath := filepath.Join(env.GoPath, "bin", binName)
	info, err := bin.LoadBinInfo(binPath)
	if err != nil {
		return err
	}
	return tryUpgradeBin(info)
}

func tryUpgradeBin(info *bin.Info) error {
	fmt.Printf("checking for updates for %s\n", info.Name)
	latestVersion, err := version.FetchLatest(info.Mod)
	if err != nil {
		return fmt.Errorf("failed to fetch latest version: %v", err)
	}
	if version.Compare(info.Version, latestVersion) < 0 {
		fmt.Printf("upgrading %s from %s to %s\n", info.Name, info.Version, latestVersion)
		return bin.InstallBinByVersion(info.Path, latestVersion)
	}
	fmt.Printf("%s is already the latest version\n", info.Name)
	return nil
}
