package version

import (
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
		flag: cli.NewFlags("version"),
		v:    vbs.New(),
	}

	cmd.flag.SetBool(&cmd.v.Verbose, "verbose", false)
	cmd.flag.SetBools(&cmd.help, "h", "help", false)

	return cmd
}
