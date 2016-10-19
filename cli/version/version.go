package version

import (
	"fmt"

	"github.com/variadico/noti/cmd/noti/cli"
	"github.com/variadico/vbs"
)

type Command struct {
	flag *cli.Flags
	v    vbs.Printer
}

func (c *Command) Parse(args []string) error {
	if err := c.flag.Parse(args); err != nil {
		return err
	}

	c.v.Verbose = c.flag.Verbose
	return nil
}

func (c *Command) Run() error {
	c.v.Println("Running version command")

	if c.flag.Help {
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

	return cmd
}
