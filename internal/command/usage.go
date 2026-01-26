package command

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"text/tabwriter"
)

func PrintUsage(fs *flag.FlagSet, usageLine, description string) {
	output := fs.Output()
	if output == nil {
		output = os.Stderr
	}
	_, _ = fmt.Fprintf(output, "Usage: %s\n", usageLine)
	if description != "" {
		_, _ = fmt.Fprintf(output, "\n%s\n", description)
	}
	_, _ = fmt.Fprintln(output)
	_, _ = fmt.Fprintln(output, "Options:")
	if !hasFlags(fs) {
		_, _ = fmt.Fprintln(output, "  (none)")
		return
	}
	PrintFlagDefaults(fs, output)
}

func PrintCommandsUsage(commands map[string]*Command, program, description string) {
	fmt.Printf("Usage: %s <command> [options]\n", program)
	if description != "" {
		fmt.Printf("\n%s\n", description)
	}
	fmt.Printf("\nCommands:\n")
	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	for _, name := range sortedCommandNames(commands) {
		cmd := commands[name]
		summary := cmd.Usage
		if summary == "" {
			summary = cmd.FlagSet.Name()
		}
		_, _ = fmt.Fprintf(w, "  %s\t%s\n", name, summary)
	}
	_ = w.Flush()
	fmt.Printf("\nRun \"%s <command> -help\" for details.\n", program)
}

func PrintFlagDefaults(fs *flag.FlagSet, output io.Writer) {
	w := tabwriter.NewWriter(output, 0, 4, 2, ' ', 0)
	fs.VisitAll(func(f *flag.Flag) {
		option := fmt.Sprintf("-%s", f.Name)
		if !isBoolFlag(f) {
			option = fmt.Sprintf("%s <value>", option)
		}
		usage := f.Usage
		if f.DefValue != "" {
			usage = fmt.Sprintf("%s (default: %s)", usage, f.DefValue)
		}
		_, _ = fmt.Fprintf(w, "  %s\t%s\n", option, usage)
	})
	_ = w.Flush()
}

func sortedCommandNames(commands map[string]*Command) []string {
	names := make([]string, 0, len(commands))
	for name := range commands {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func hasFlags(fs *flag.FlagSet) bool {
	hasAny := false
	fs.VisitAll(func(*flag.Flag) {
		hasAny = true
	})
	return hasAny
}

func isBoolFlag(f *flag.Flag) bool {
	if flagValue, ok := f.Value.(interface{ IsBoolFlag() bool }); ok {
		return flagValue.IsBoolFlag()
	}
	return false
}
