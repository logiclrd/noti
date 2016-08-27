package version

import (
	"flag"
	"fmt"

	"github.com/variadico/noti/cmd/noti/cli"
	"github.com/variadico/vbs"
)

type Command struct {
	flag cli.Flags
	v    vbs.Printer
	help bool
}

func (c *Command) Parse(args []string) error {
	return c.flag.Parse(args)
}

func (c *Command) Run() error {
	c.v.Println("Running version command")

	if c.help {
		fmt.Println("noti version [-verbose -h -help]")
		return nil
	}

	fmt.Println("noti v3.0.0")
	c.v.Println("Looking up latest version on GitHub")

	return nil
}

func NewCommand() cli.Cmd {
	cmd := &Command{
		flag: cli.Flags{flag.NewFlagSet("version", flag.ExitOnError)},
		v:    vbs.New(),
	}

	cmd.flag.BoolVar(&cmd.v.Verbose, "verbose", false, "Enable verbose mode")
	cmd.flag.BoolVar(&cmd.help, "h", false, "Show help")
	cmd.flag.BoolVar(&cmd.help, "help", false, "Show help")

	return cmd
}
