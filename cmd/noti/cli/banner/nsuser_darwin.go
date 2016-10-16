package banner

import (
	"fmt"
	"runtime"

	"github.com/variadico/noti/cmd/noti/cli"
	"github.com/variadico/noti/cmd/noti/config"
	"github.com/variadico/noti/cmd/noti/triggers"
	"github.com/variadico/noti/nsuser"
	"github.com/variadico/vbs"
)

var cmdDefault = &nsuser.Notification{
	Title:           "{{.Cmd}}",
	InformativeText: "Done!",
	SoundName:       "Ping",
}

type Command struct {
	flag *cli.Flags
	v    vbs.Printer
	n    *nsuser.Notification
}

func (c *Command) Parse(args []string) error {
	if err := c.flag.Parse(args); err != nil {
		return err
	}

	c.v.Verbose = c.flag.Verbose
	return nil
}

func (c *Command) Notify(stats triggers.Stats) error {
	conf, err := config.File()
	if err != nil {
		c.v.Println(err)
	} else {
		c.v.Println("Found config file")
	}

	fromFlags := new(nsuser.Notification)

	if c.flag.Passed("title", "t") {
		fromFlags.Title = c.n.Title
	}
	if c.flag.Passed("subtitle") {
		fromFlags.Subtitle = c.n.Subtitle
	}
	if c.flag.Passed("message", "m") {
		fromFlags.InformativeText = c.n.InformativeText
	}
	if c.flag.Passed("icon") {
		fromFlags.ContentImage = c.n.ContentImage
	}
	if c.flag.Passed("sound") {
		fromFlags.SoundName = c.n.SoundName
	}

	c.v.Println("Evaluating")
	c.v.Printf("Default: %+v\n", cmdDefault)
	c.v.Printf("Config: %+v\n", conf.Banner)
	c.v.Printf("Flags: %+v\n", fromFlags)

	config.EvalStringFields(cmdDefault, stats)
	config.EvalStringFields(conf.Banner, stats)
	config.EvalStringFields(fromFlags, stats)

	c.v.Println("Merging")
	merged := new(nsuser.Notification)
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

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ts := []string(c.flag.Triggers)
	return triggers.Run(ts, c.flag.Args(), c.Notify)
}

func NewCommand() cli.NotifyCmd {
	cmd := &Command{
		flag: cli.NewFlags("banner"),
		v:    vbs.New(),
		n:    new(nsuser.Notification),
	}

	cmd.flag.SetStrings(&cmd.n.Title, "t", "title", cmdDefault.Title)
	cmd.flag.SetStrings(&cmd.n.InformativeText, "m", "message", cmdDefault.InformativeText)

	cmd.flag.SetString(&cmd.n.Subtitle, "subtitle", cmdDefault.Subtitle)
	cmd.flag.SetString(&cmd.n.ContentImage, "icon", cmdDefault.ContentImage)
	cmd.flag.SetString(&cmd.n.SoundName, "sound", cmdDefault.SoundName)

	return cmd
}
