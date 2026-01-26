package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/tbxark/gbvm/internal/bin"
	"github.com/tbxark/gbvm/internal/command"
)

func NewListCommand() *command.Command {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	fs.Usage = func() {
		command.PrintUsage(fs, "gbvm list [options]", "List all installed Go binaries")
	}
	showVersion := fs.Bool("versions", false, "show version")
	jsonMode := fs.Bool("json", false, "json mode")
	verbose := fs.Bool("verbose", false, "show scan errors")
	cmd := command.NewCommand(fs, func() error {
		return handleList(*jsonMode, *showVersion, *verbose)
	})
	cmd.Usage = "List all installed Go binaries"
	return cmd
}

func handleList(jsonMode, showVersion, verbose bool) error {
	versions, err := bin.LoadAllBinVersions(verbose)
	if err != nil {
		return err
	}
	if jsonMode {
		encoded, e := json.MarshalIndent(versions, "", "  ")
		if e != nil {
			return e
		}
		fmt.Println(string(encoded))
		return nil
	}
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	for _, v := range versions {
		if showVersion {
			_, _ = fmt.Fprintf(writer, "%s\t%s\n", v.Name, v.Version)
		} else {
			_, _ = fmt.Fprintln(writer, v.Name)
		}
	}
	return writer.Flush()
}
