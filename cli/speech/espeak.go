// +build !darwin
// +build !windows

package speech

import (
	"fmt"

	"github.com/variadico/noti/cli"
	"github.com/variadico/noti/config"
	"github.com/variadico/noti/ntypes/espeak"
	"github.com/variadico/noti/runstat"
	"github.com/variadico/noti/triggers"
	"github.com/variadico/vbs"
)

var cmdDefault = &espeak.Notification{
	Text:      "{{.Cmd}} done!",
	VoiceName: "english-us",
}

type Command struct {
	flag *cli.Flags
	v    vbs.Printer
	n    *espeak.Notification
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

	fromFlags := new(espeak.Notification)

	if c.flag.Passed("message", "m") {
		fromFlags.Text = c.n.Text
	}
	if c.flag.Passed("word-gap") {
		fromFlags.WordGap = c.n.WordGap
	}
	if c.flag.Passed("pitch") {
		fromFlags.PitchAdjustment = c.n.PitchAdjustment
	}
	if c.flag.Passed("rate") {
		fromFlags.Rate = c.n.Rate
	}
	if c.flag.Passed("voice-name") {
		fromFlags.VoiceName = c.n.VoiceName
	}

	c.v.Println("Evaluating")
	c.v.Printf("Default: %+v\n", cmdDefault)
	c.v.Printf("Config: %+v\n", conf.Speech)
	c.v.Printf("Flags: %+v\n", fromFlags)

	config.EvalStringFields(cmdDefault, stats)
	config.EvalStringFields(conf.Speech, stats)
	config.EvalStringFields(fromFlags, stats)

	c.v.Println("Merging")
	merged := new(espeak.Notification)
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
		n:    new(espeak.Notification),
	}

	cmd.flag.SetStrings(&cmd.n.Text, "m", "message", cmdDefault.Text)

	cmd.flag.SetString(&cmd.n.VoiceName, "voice", cmdDefault.VoiceName)
	cmd.flag.SetInt(&cmd.n.Rate, "rate", cmdDefault.Rate)
	cmd.flag.SetInt(&cmd.n.PitchAdjustment, "pitch", cmdDefault.PitchAdjustment)
	cmd.flag.SetInt(&cmd.n.WordGap, "word-gap", cmdDefault.WordGap)

	return cmd
}
