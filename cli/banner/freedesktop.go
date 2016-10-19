// +build !darwin
// +build !windows

package banner

import (
	"fmt"

	"github.com/variadico/noti/cli"
	"github.com/variadico/noti/config"
	"github.com/variadico/noti/ntypes/freedesktop"
	"github.com/variadico/noti/runstat"
	"github.com/variadico/noti/triggers"
	"github.com/variadico/vbs"
)

var cmdDefault = &freedesktop.Notification{
	Summary:       "{{.Cmd}}",
	Body:          "Done!",
	ExpireTimeout: 500,
}

type Command struct {
	flag *cli.Flags
	v    vbs.Printer
	n    *freedesktop.Notification
}

func (c *Command) Parse(args []string) error {
	if err := c.flag.Parse(args); err != nil {
		return err
	}

	c.v.Verbose = c.flag.Verbose
	return nil
}

func (c *Command) Notify(stats runstat.Result) error {
	conf, err := config.File()
	if err != nil {
		c.v.Println(err)
	} else {
		c.v.Println("Found config file")
	}

	fromFlags := new(freedesktop.Notification)

	if c.flag.Passed("title", "t") {
		fromFlags.Summary = c.n.Summary
	}
	if c.flag.Passed("message", "m") {
		fromFlags.Body = c.n.Body
	}
	if c.flag.Passed("app-name") {
		fromFlags.AppName = c.n.AppName
	}
	if c.flag.Passed("replaces-id") {
		fromFlags.ReplacesID = c.n.ReplacesID
	}
	if c.flag.Passed("icon") {
		fromFlags.AppIcon = c.n.AppIcon
	}
	if c.flag.Passed("timeout") {
		fromFlags.ExpireTimeout = c.n.ExpireTimeout
	}

	c.v.Println("Evaluating")
	c.v.Printf("Default: %+v\n", cmdDefault)
	c.v.Printf("Config: %+v\n", conf.Banner)
	c.v.Printf("Flags: %+v\n", fromFlags)

	config.EvalStringFields(cmdDefault, stats)
	config.EvalStringFields(conf.Banner, stats)
	config.EvalStringFields(fromFlags, stats)

	c.v.Println("Merging")
	merged := new(freedesktop.Notification)
	err = config.MergeFields(
		merged,
		cmdDefault,
		conf.Banner,
		fromFlags,
	)
	if err != nil {
		return err
	}
	c.v.Printf("Merge result: %+v\n", merged)

	c.v.Println("Sending notification")
	err = merged.Send()
	c.v.Println("Sent notification")
	return err
}

func (c *Command) Run() error {
	if c.flag.Help {
		fmt.Println(helpText)
		return nil
	}

	return triggers.Run([]string(c.flag.Triggers), c.flag.Args(), c.Notify)
}

func NewCommand() cli.NotifyCmd {
	cmd := &Command{
		flag: cli.NewFlags("banner"),
		v:    vbs.New(),
		n:    new(freedesktop.Notification),
	}

	cmd.flag.SetStrings(&cmd.n.Summary, "t", "title", cmdDefault.Summary)
	cmd.flag.SetStrings(&cmd.n.Body, "m", "message", cmdDefault.Body)

	cmd.flag.SetString(&cmd.n.AppName, "app-name", cmdDefault.AppName)
	cmd.flag.SetUint(&cmd.n.ReplacesID, "replaces-id", cmdDefault.ReplacesID)
	cmd.flag.SetString(&cmd.n.AppIcon, "icon", cmdDefault.AppIcon)
	cmd.flag.SetInt(&cmd.n.ExpireTimeout, "timeout", cmdDefault.ExpireTimeout)

	return cmd
}
