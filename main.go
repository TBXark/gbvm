package main

import (
	"fmt"

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
		fmt.Printf("Usage: gbvm <command> [options]\n\n")
		fmt.Printf("A command line tool to manage Go binaries\n\n")
		for name, sub := range commands {
			fmt.Printf("gbvm %s:\n", name)
			sub.FlagSet.PrintDefaults()
			fmt.Println()
		}
	}
	command.Execute(commands, printDefaults)
}
