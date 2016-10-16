package slack

import (
	"fmt"
	"net/http"
	"time"

	"github.com/variadico/noti/cmd/noti/cli"
	"github.com/variadico/noti/cmd/noti/config"
	"github.com/variadico/noti/cmd/noti/triggers"
	"github.com/variadico/noti/slack"
	"github.com/variadico/vbs"
)

var cmdDefault = &slack.Notification{
	Text:        "{{.Cmd}} done!",
	Parse:       slack.ParseFull,
	LinkNames:   slack.LinkNamesOn,
	UnfurlLinks: true,
	UnfurlMedia: true,
	Username:    "Noti",
}

type Command struct {
	flag *cli.Flags
	v    vbs.Printer

	n *slack.Notification
}

func (c *Command) Parse(args []string) error {
	if err := c.flag.Parse(args); err != nil {
		return err
	}

	c.v.Verbose = c.flag.Verbose
	return nil
}

func (c *Command) Notify(stats triggers.Stats) error {
	if c.flag.Help {
		fmt.Println(helpText)
		return nil
	}

	conf, err := config.File()
	if err != nil {
		c.v.Println(err)
	} else {
		c.v.Println("Found config file")
	}

	fromFlags := new(slack.Notification)
	if c.flag.Passed("message", "m") {
		fromFlags.Text = c.n.Text
	}
	if c.flag.Passed("token") {
		fromFlags.Token = c.n.Token
	}
	if c.flag.Passed("channel") {
		fromFlags.Channel = c.n.Channel
	}
	if c.flag.Passed("parse") {
		fromFlags.Parse = c.n.Parse
	}
	if c.flag.Passed("link-names") {
		fromFlags.LinkNames = c.n.LinkNames
	}
	if c.flag.Passed("unfurl-links") {
		fromFlags.UnfurlLinks = c.n.UnfurlLinks
	}
	if c.flag.Passed("unfurl-media") {
		fromFlags.UnfurlMedia = c.n.UnfurlMedia
	}
	if c.flag.Passed("username") {
		fromFlags.Username = c.n.Username
	}
	if c.flag.Passed("as-user") {
		fromFlags.AsUser = c.n.AsUser
	}
	if c.flag.Passed("icon-url") {
		fromFlags.IconURL = c.n.IconURL
	}
	if c.flag.Passed("icon-emoji") {
		fromFlags.IconEmoji = c.n.IconEmoji
	}

	c.v.Println("Evaluating")
	c.v.Printf("Default: %+v\n", cmdDefault)
	c.v.Printf("Config: %+v\n", conf.Slack)
	c.v.Printf("Flags: %+v\n", fromFlags)

	config.EvalStringFields(cmdDefault, stats)
	config.EvalStringFields(conf.Slack, stats)
	config.EvalStringFields(fromFlags, stats)

	c.v.Println("Merging")
	merged := new(slack.Notification)
	err = config.MergeFields(
		merged,
		cmdDefault,
		conf.Slack,
		fromFlags,
	)
	if err != nil {
		return err
	}
	c.v.Printf("Merge result: %+v\n", merged)

	merged.SetClient(&http.Client{Timeout: 5 * time.Second})

	c.v.Println("Sending notification")
	err = merged.Send()
	c.v.Println("Sent notification")
	return err
}

func (c *Command) Run() error {
	return nil
}

func NewCommand() cli.NotifyCmd {
	cmd := &Command{
		flag: cli.NewFlags("slack"),
		v:    vbs.New(),
		n:    new(slack.Notification),
	}

	cmd.flag.SetStrings(&cmd.n.Text, "m", "message", cmdDefault.Text)

	cmd.flag.SetString(&cmd.n.Token, "token", cmdDefault.Token)
	cmd.flag.SetString(&cmd.n.Channel, "channel", cmdDefault.Channel)
	cmd.flag.SetString(&cmd.n.Parse, "parse", cmdDefault.Parse)
	cmd.flag.SetInt(&cmd.n.LinkNames, "link-names", cmdDefault.LinkNames)
	cmd.flag.SetBool(&cmd.n.UnfurlLinks, "unfurl-links", cmdDefault.UnfurlLinks)
	cmd.flag.SetBool(&cmd.n.UnfurlMedia, "unfurl-media", cmdDefault.UnfurlMedia)
	cmd.flag.SetString(&cmd.n.Username, "username", cmdDefault.Username)
	cmd.flag.SetBool(&cmd.n.AsUser, "as-user", cmdDefault.AsUser)
	cmd.flag.SetString(&cmd.n.IconURL, "icon-url", cmdDefault.IconURL)
	cmd.flag.SetString(&cmd.n.IconEmoji, "icon-emoji", cmdDefault.IconEmoji)

	return cmd
}
