package cmd

import (
	"flag"
	"fmt"

	"github.com/TBXark/gbvm/internal/bin"
	"github.com/TBXark/gbvm/internal/command"
	"github.com/TBXark/gbvm/internal/env"
	"github.com/TBXark/gbvm/internal/log"
	"github.com/TBXark/gbvm/internal/version"
)

func NewUpgradeCommand() *command.Command {
	fs := flag.NewFlagSet("upgrade", flag.ExitOnError)
	fs.Usage = func() {
		command.PrintUsage(fs, "gbvm upgrade [options] [bin1 bin2 ...]", "Upgrade Go binaries")
	}
	skipDev := fs.Bool("skip-dev", false, "skip dev version")
	verbose := fs.Bool("verbose", false, "show scan errors")
	cmd := command.NewCommand(fs, func() error {
		if fs.NArg() == 0 {
			return upgradeAllBins(*skipDev, *verbose)
		}
		var failures int
		for _, binName := range fs.Args() {
			if err := upgradeBin(binName); err != nil {
				failures++
				log.Errorf("failed to upgrade %s: %v", binName, err)
			}
		}
		if failures > 0 {
			return fmt.Errorf("%d upgrade(s) failed", failures)
		}
		return nil
	})
	cmd.Usage = "Upgrade Go binaries"
	return cmd
}

func upgradeAllBins(skipDev, verbose bool) error {
	versions, err := bin.LoadAllBinVersions(verbose)
	if err != nil {
		return err
	}
	for _, v := range versions {
		if skipDev && v.Version == env.DevelVersion {
			continue
		}
		if e := tryUpgradeBin(v); e != nil {
			log.Errorf("failed to upgrade %s: %v", v.Name, e)
		}
	}
	return nil
}

func upgradeBin(binName string) error {
	info, err := bin.FindBinInfo(binName)
	if err != nil {
		return err
	}
	return tryUpgradeBin(info)
}

func tryUpgradeBin(info *bin.Info) error {
	log.Debugf("checking for updates for %s", info.Name)
	latestVersion, proxy, err := version.FetchLatest(info.Mod)
	if err != nil {
		return fmt.Errorf("failed to fetch latest version: %v", err)
	}
	shouldUpgrade, message := determineUpgrade(info, latestVersion)
	log.Debug(message)
	if shouldUpgrade {
		return bin.InstallBinByVersion(info.Path, latestVersion, proxy)
	}
	return nil
}

func determineUpgrade(info *bin.Info, latestVersion string) (bool, string) {
	if version.Compare(info.Version, latestVersion) < 0 {
		return true, fmt.Sprintf("upgrading %s from %s to %s", info.Name, info.Version, latestVersion)
	}
	return false, fmt.Sprintf("%s is already the latest version", info.Name)
}
