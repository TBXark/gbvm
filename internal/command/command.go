package command

import (
	"flag"
	"os"

	"github.com/TBXark/gbvm/internal/log"
)

type Command struct {
	Name       string
	Usage      string
	FlagSet    *flag.FlagSet
	HandleFunc func() error
}

func NewCommand(fs *flag.FlagSet, handleFunc func() error) *Command {
	cmd := &Command{
		Name:    fs.Name(),
		FlagSet: fs,
	}
	help := cmd.FlagSet.Bool("help", false, "show help")
	cmd.HandleFunc = func() error {
		if *help {
			fs.Usage()
			return nil
		}
		return handleFunc()
	}
	return cmd
}

func Execute(commands map[string]*Command, printDefaults func()) {
	if len(os.Args) < 2 {
		printDefaults()
		return
	}

	cmd, exists := commands[os.Args[1]]
	if !exists {
		printDefaults()
		return
	}

	if err := cmd.FlagSet.Parse(os.Args[2:]); err != nil {
		log.Error(err)
		return
	}

	if err := cmd.HandleFunc(); err != nil {
		log.Error(err)
	}
}
