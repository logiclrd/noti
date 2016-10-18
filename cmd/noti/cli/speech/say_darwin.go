package speech

import (
	"fmt"

	"github.com/variadico/noti/cmd/noti/cli"
	"github.com/variadico/noti/cmd/noti/config"
	"github.com/variadico/noti/cmd/noti/runstat"
	"github.com/variadico/noti/cmd/noti/triggers"
	"github.com/variadico/noti/say"
	"github.com/variadico/vbs"
)

var cmdDefault = &say.Notification{
	Voice: "Alex",
	Text:  "{{.Cmd}} done!",
	Rate:  200,
}

type Command struct {
	flag *cli.Flags
	v    vbs.Printer
	n    *say.Notification
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

	fromFlags := new(say.Notification)

	if c.flag.Passed("rate") {
		fromFlags.Rate = c.n.Rate
	}
	if c.flag.Passed("message", "m") {
		fromFlags.Text = c.n.Text
	}
	if c.flag.Passed("rate") {
		fromFlags.Rate = c.n.Rate
	}

	c.v.Println("Evaluating")
	c.v.Printf("Default: %+v\n", cmdDefault)
	c.v.Printf("Config: %+v\n", conf.Speech)
	c.v.Printf("Flags: %+v\n", fromFlags)

	config.EvalStringFields(cmdDefault, stats)
	config.EvalStringFields(conf.Speech, stats)
	config.EvalStringFields(fromFlags, stats)

	c.v.Println("Merging")
	merged := new(say.Notification)
	err = config.MergeFields(
		merged,
		cmdDefault,
		conf.Speech,
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
		flag: cli.NewFlags("speech"),
		v:    vbs.New(),
		n:    new(say.Notification),
	}

	cmd.flag.SetStrings(&cmd.n.Text, "m", "message", cmdDefault.Text)

	cmd.flag.SetString(&cmd.n.Voice, "voice", cmdDefault.Voice)
	cmd.flag.SetInt(&cmd.n.Rate, "rate", cmdDefault.Rate)

	return cmd
}
