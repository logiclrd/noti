package banner

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/variadico/noti/cmd/noti/cli"
	"github.com/variadico/noti/cmd/noti/config"
	"github.com/variadico/noti/cmd/noti/run"
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

	ktimeout string
	timeout  string
	contains string
}

func (c *Command) Parse(args []string) error {
	return c.flag.Parse(args)
}

func (c *Command) Notify(stats run.Stats) error {
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

	if c.ktimeout != "" {
		d, err := time.ParseDuration(c.ktimeout)
		if err != nil {
			return err
		}

		fmt.Println(">>>>>>>> EXEC TIMEOUT!")
		c.v.Println("Executing command with timeout")
		stats := run.ExecWithTimeout(d, c.flag.Args()...)
		return c.Notify(stats)
	}

	if c.timeout != "" {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		fmt.Println(">>>>>>>> EXEC NOTIFY!")
		stats := run.ExecNotify(ctx, c.flag.Args()...)
		for s := range stats {
			fmt.Println(">>>>>>>> SENDING NOTI!")
			err := c.Notify(s)
			if err != nil {
				return err
			}
		}
		return nil
	}

	if c.contains != "" {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		fmt.Println(">>>>>>>> EXEC CONTAINS!")
		stats := run.ExecContains(ctx, c.flag.Args()...)
		for s := range stats {
			fmt.Println(">>>>>>>> SENDING NOTI!")
			err := c.Notify(s)
			if err != nil {
				return err
			}
		}
		return nil
	}

	c.v.Println("Executing command")
	fmt.Println(">>>>>>>>  EXEC!")
	return c.Notify(run.Exec(c.flag.Args()...))
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

	cmd.flag.SetBool(&cmd.v.Verbose, "verbose", false)

	cmd.flag.SetString(&cmd.ktimeout, "ktimeout", "")
	cmd.flag.SetString(&cmd.timeout, "timeout", "")
	cmd.flag.SetString(&cmd.contains, "contains", "")

	return cmd
}
