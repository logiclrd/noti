package root

import (
	"fmt"
	"log"

	"github.com/variadico/noti/cmd/noti/cli"
	"github.com/variadico/noti/cmd/noti/config"
	"github.com/variadico/noti/cmd/noti/triggers"
	"github.com/variadico/vbs"
)

type Command struct {
	flag *cli.Flags
	v    vbs.Printer

	Cmds map[string]cli.Cmd
}

func (c *Command) Args() []string {
	return c.flag.Args()
}

func (c *Command) Parse(args []string) error {
	if err := c.flag.Parse(args); err != nil {
		return err
	}

	c.v.Verbose = c.flag.Verbose
	return nil
}

func (c *Command) Run() error {
	c.v.Println("Running noti command")

	if c.flag.Help {
		fmt.Println(helpText)
		return nil
	}

	return nil
}

func (c *Command) Notify(stats triggers.Stats) error {
	c.v.Println("Notifying")

	conf, err := config.File()
	if err != nil {
		c.v.Println(err)
	} else {
		c.v.Println("Found config file")
	}

	// Read default set of notification types.
	if len(conf.DefaultSet) == 0 {
		conf.DefaultSet = append(conf.DefaultSet, "banner")
	}

	for _, sub := range conf.DefaultSet {
		subCmd, found := c.Cmds[sub]
		if !found {
			log.Println("Unknown subcommand:", sub)
			continue
		}

		ncmd, is := subCmd.(cli.NotifyCmd)
		if !is {
			continue
		}

		if err := ncmd.Notify(stats); err != nil {
			log.Printf("Failed to run %s: %s", sub, err)
		}
	}

	return nil
}

func NewCommand() cli.NotifyCmd {
	cmd := &Command{
		flag: cli.NewFlags("noti"),
		v:    vbs.New(),
	}

	return cmd
}
