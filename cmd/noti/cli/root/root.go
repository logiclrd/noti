package root

import (
	"flag"
	"fmt"
	"log"

	"github.com/variadico/noti/cmd/noti/cli"
	"github.com/variadico/noti/cmd/noti/config"
	"github.com/variadico/noti/cmd/noti/run"
	"github.com/variadico/vbs"
)

type Command struct {
	flag cli.Flags
	v    vbs.Printer
	help bool

	Cmds map[string]cli.Cmd
}

func (c *Command) Args() []string {
	return c.flag.Args()
}

func (c *Command) Parse(args []string) error {
	return c.flag.Parse(args)
}

func (c *Command) Run() error {
	c.v.Println("Running noti command")

	if c.help {
		fmt.Println("noti [-verbose -h -help] [notification type] [command]")
		return nil
	}

	c.v.Println("Executing command")
	stats := run.Exec(c.flag.Args()...)
	c.v.Println("Executed command")

	return c.Notify(stats)
}

func (c *Command) Notify(stats run.Stats) error {
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
		flag: cli.Flags{flag.NewFlagSet("noti", flag.ContinueOnError)},
		v:    vbs.New(),
	}

	cmd.flag.BoolVar(&cmd.v.Verbose, "verbose", false, "Enable verbose mode")
	cmd.flag.BoolVar(&cmd.help, "h", false, "Show help")
	cmd.flag.BoolVar(&cmd.help, "help", false, "Show help")

	return cmd
}
