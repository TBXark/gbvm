package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/TBXark/gbvm/internal/bin"
	"github.com/TBXark/gbvm/internal/command"
)

func NewListCommand() *command.Command {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Printf("Usage: gbvm list [options]\n\n")
		fmt.Printf("List all installed Go binaries\n\n")
		fs.PrintDefaults()
	}
	showVersion := fs.Bool("versions", false, "show version")
	jsonMode := fs.Bool("json", false, "json mode")
	cmd := command.NewCommand(fs, func() error {
		return handleList(*jsonMode, *showVersion)
	})
	return cmd
}

func handleList(jsonMode, showVersion bool) error {
	versions, err := bin.LoadAllBinVersions()
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
