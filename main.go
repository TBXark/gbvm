package main

import (
	"github.com/TBXark/gbvm/cmd"
	"github.com/TBXark/gbvm/internal/command"
)

func main() {
	commands := map[string]*command.Command{
		"list":    cmd.NewListCommand(),
		"upgrade": cmd.NewUpgradeCommand(),
		"install": cmd.NewInstallCommand(),
	}
	printDefaults := func() {
		command.PrintCommandsUsage(commands, "gbvm", "A command line tool to manage Go binaries")
	}
	command.Execute(commands, printDefaults)
}
