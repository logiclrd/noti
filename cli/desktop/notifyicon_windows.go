package desktop

import (
	"fmt"

	"github.com/variadico/noti/cli"
	"github.com/variadico/noti/config"
	"github.com/variadico/noti/runstat"
	"github.com/variadico/noti/services/notifyicon"
	"github.com/variadico/noti/triggers"
	"github.com/variadico/vbs"
)

var cmdDefault = &notifyicon.Notification{
	BalloonTipTitle: "{{.Cmd}}",
	BalloonTipText:  "Done!",
	BalloonTipIcon:  "Info",
}

type Command struct {
	flag *cli.Flags
	v    vbs.Printer
	n    *notifyicon.Notification
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

	fromFlags := new(notifyicon.Notification)

	if c.flag.Passed("title", "t") {
		fromFlags.BalloonTipTitle = c.n.BalloonTipTitle
	}
	if c.flag.Passed("message", "m") {
		fromFlags.BalloonTipText = c.n.BalloonTipText
	}

	if c.flag.Passed("icon", "i") {
		fromFlags.BalloonTipIcon = c.n.BalloonTipIcon
	}

	c.v.Println("Evaluating")
	c.v.Printf("Default: %+v\n", cmdDefault)
	c.v.Printf("Config: %+v\n", conf.Desktop)
	c.v.Printf("Flags: %+v\n", fromFlags)

	config.EvalStringFields(cmdDefault, stats)
	config.EvalStringFields(conf.Desktop, stats)
	config.EvalStringFields(fromFlags, stats)

	c.v.Println("Merging")
	merged := new(notifyicon.Notification)
	err = config.MergeFields(
		merged,
		cmdDefault,
		conf.Desktop,
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
		flag: cli.NewFlags("desktop"),
		v:    vbs.New(),
		n:    new(notifyicon.Notification),
	}

	cmd.flag.SetStrings(&cmd.n.BalloonTipTitle, "t", "title", cmdDefault.BalloonTipTitle)
	cmd.flag.SetStrings(&cmd.n.BalloonTipText, "m", "message", cmdDefault.BalloonTipText)

	cmd.flag.SetStrings(&cmd.n.BalloonTipIcon, "i", "icon", cmdDefault.BalloonTipIcon)

	return cmd
}
