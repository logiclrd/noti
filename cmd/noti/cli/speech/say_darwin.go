package speech

import (
	"flag"
	"fmt"

	"github.com/variadico/noti/cmd/noti/cli"
	"github.com/variadico/noti/cmd/noti/config"
	"github.com/variadico/noti/cmd/noti/run"
	"github.com/variadico/noti/say"
	"github.com/variadico/vbs"
)

var cmdDefault = &say.Notification{
	Voice: "Alex",
	Text:  "{{.Cmd}} done!",
	Rate:  200,
}

func ptrs(n *say.Notification) []interface{} {
	if n == nil {
		return nil
	}

	return []interface{}{
		&n.Text,
		&n.Voice,
		&n.Rate,
	}
}

type Command struct {
	flag cli.Flags
	v    vbs.Printer
	help bool

	n *say.Notification
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

	fromFlags := new(say.Notification)

	if c.flag.Set("rate") {
		fromFlags.Rate = c.n.Rate
	}
	if c.flag.Set("message", "m") {
		fromFlags.Text = c.n.Text
	}
	if c.flag.Set("rate") {
		fromFlags.Rate = c.n.Rate
	}

	c.v.Println("Evaluating")
	c.v.Printf("Default: %+v\n", cmdDefault)
	c.v.Printf("Config: %+v\n", conf.Speech)
	c.v.Printf("Flags: %+v\n", fromFlags)

	config.EvalFields(ptrs(cmdDefault), stats)
	config.EvalFields(ptrs(conf.Speech), stats)
	config.EvalFields(ptrs(fromFlags), stats)

	c.v.Println("Merging")
	merged := new(say.Notification)
	err = config.MergeFields(
		ptrs(merged),
		ptrs(cmdDefault),
		ptrs(conf.Speech),
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
		flag: cli.Flags{flag.NewFlagSet("speech", flag.ExitOnError)},
		v:    vbs.New(),
		n:    new(say.Notification),
	}

	cmd.flag.StringVar(&cmd.n.Text, "message", cmdDefault.Text, "Message")
	cmd.flag.StringVar(&cmd.n.Text, "m", cmdDefault.Text, "Message")
	cmd.flag.StringVar(&cmd.n.Voice, "voice", cmdDefault.Voice, "Voice")
	cmd.flag.IntVar(&cmd.n.Rate, "rate", cmdDefault.Rate, "Rate")

	cmd.flag.BoolVar(&cmd.v.Verbose, "verbose", false, "Enable verbose mode")
	cmd.flag.BoolVar(&cmd.help, "h", false, "Show help")
	cmd.flag.BoolVar(&cmd.help, "help", false, "Show help")

	return cmd
}
