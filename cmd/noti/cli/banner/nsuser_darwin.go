package banner

import (
	"flag"
	"fmt"

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

func ptrs(n *nsuser.Notification) []interface{} {
	return []interface{}{
		&n.Title,
		&n.Subtitle,
		&n.InformativeText,
		&n.ContentImage,
		&n.SoundName,
	}
}

type Command struct {
	flag cli.Flags
	v    vbs.Printer
	help bool

	n *nsuser.Notification
}

func (c *Command) Parse(args []string) error {
	return c.flag.Parse(args)
}

func (c *Command) Notify(stats run.Stats) error {
	if c.help {
		fmt.Println(helpText)
		return nil
	}

	conf, err := config.File()
	if err != nil {
		c.v.Println(err)
	} else {
		c.v.Println("Found config file")
	}

	fromFlags := new(nsuser.Notification)

	if c.flag.Set("title", "t") {
		fromFlags.Title = c.n.Title
	}
	if c.flag.Set("subtitle") {
		fromFlags.Subtitle = c.n.Subtitle
	}
	if c.flag.Set("message", "m") {
		fromFlags.InformativeText = c.n.InformativeText
	}
	if c.flag.Set("icon") {
		fromFlags.ContentImage = c.n.ContentImage
	}
	if c.flag.Set("sound") {
		fromFlags.SoundName = c.n.SoundName
	}

	c.v.Println("Evaluating")
	c.v.Printf("Default: %+v\n", cmdDefault)
	c.v.Printf("Config: %+v\n", conf.Banner)
	c.v.Printf("Flags: %+v\n", fromFlags)

	config.EvalFields(ptrs(cmdDefault), stats)
	config.EvalFields(ptrs(conf.Banner), stats)
	config.EvalFields(ptrs(fromFlags), stats)

	c.v.Println("Merging")
	merged := new(nsuser.Notification)
	err = config.MergeFields(
		ptrs(merged),
		ptrs(cmdDefault),
		ptrs(conf.Banner),
		ptrs(fromFlags),
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
	c.v.Println("Executing command")
	stats := run.Exec(c.flag.Args()...)
	c.v.Println("Executed command")
	c.v.Printf("Run stats: %+v\n", stats)

	return c.Notify(stats)
}

func NewCommand() cli.NotifyCmd {
	cmd := &Command{
		flag: cli.Flags{flag.NewFlagSet("banner", flag.ExitOnError)},
		v:    vbs.New(),
		n:    new(nsuser.Notification),
	}

	cmd.flag.StringVar(&cmd.n.Title, "title", cmdDefault.Title, "Title")
	cmd.flag.StringVar(&cmd.n.Title, "t", cmdDefault.Title, "Title")
	cmd.flag.StringVar(&cmd.n.Subtitle, "subtitle", cmdDefault.Subtitle, "Subtitle")
	cmd.flag.StringVar(&cmd.n.InformativeText, "message", cmdDefault.InformativeText, "Message")
	cmd.flag.StringVar(&cmd.n.InformativeText, "m", cmdDefault.InformativeText, "Message")
	cmd.flag.StringVar(&cmd.n.ContentImage, "icon", cmdDefault.ContentImage, "Icon")
	cmd.flag.StringVar(&cmd.n.SoundName, "sound", cmdDefault.SoundName, "Sound")

	cmd.flag.BoolVar(&cmd.v.Verbose, "verbose", false, "Enable verbose mode")
	cmd.flag.BoolVar(&cmd.help, "h", false, "Show help")
	cmd.flag.BoolVar(&cmd.help, "help", false, "Show help")

	return cmd
}
