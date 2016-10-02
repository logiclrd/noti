package slack

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/variadico/noti/cmd/noti/cli"
	"github.com/variadico/noti/cmd/noti/config"
	"github.com/variadico/noti/cmd/noti/run"
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

func ptrs(n *slack.Notification) []interface{} {
	if n == nil {
		return nil
	}

	return []interface{}{
		&n.Token,
		&n.Channel,
		&n.Text,
		&n.Parse,
		&n.LinkNames,
		&n.UnfurlLinks,
		&n.UnfurlMedia,
		&n.Username,
		&n.AsUser,
		&n.IconURL,
		&n.IconEmoji,
	}
}

type Command struct {
	flag cli.Flags
	v    vbs.Printer
	help bool

	n *slack.Notification
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

	fromFlags := new(slack.Notification)
	if c.flag.Set("message", "m") {
		fromFlags.Text = c.n.Text
	}
	if c.flag.Set("token") {
		fromFlags.Token = c.n.Token
	}
	if c.flag.Set("channel") {
		fromFlags.Channel = c.n.Channel
	}
	if c.flag.Set("parse") {
		fromFlags.Parse = c.n.Parse
	}
	if c.flag.Set("link-names") {
		fromFlags.LinkNames = c.n.LinkNames
	}
	if c.flag.Set("unfurl-links") {
		fromFlags.UnfurlLinks = c.n.UnfurlLinks
	}
	if c.flag.Set("unfurl-media") {
		fromFlags.UnfurlMedia = c.n.UnfurlMedia
	}
	if c.flag.Set("username") {
		fromFlags.Username = c.n.Username
	}
	if c.flag.Set("as-user") {
		fromFlags.AsUser = c.n.AsUser
	}
	if c.flag.Set("icon-url") {
		fromFlags.IconURL = c.n.IconURL
	}
	if c.flag.Set("icon-emoji") {
		fromFlags.IconEmoji = c.n.IconEmoji
	}

	c.v.Println("Evaluating")
	c.v.Printf("Default: %+v\n", cmdDefault)
	c.v.Printf("Config: %+v\n", conf.Slack)
	c.v.Printf("Flags: %+v\n", fromFlags)

	config.EvalFields(ptrs(cmdDefault), stats)
	config.EvalFields(ptrs(conf.Slack), stats)
	config.EvalFields(ptrs(fromFlags), stats)

	c.v.Println("Merging")
	merged := new(slack.Notification)
	err = config.MergeFields(
		ptrs(merged),
		ptrs(cmdDefault),
		ptrs(conf.Slack),
		ptrs(fromFlags),
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
	c.v.Println("Executing command")
	stats := run.Exec(c.flag.Args()...)
	c.v.Println("Executed command")
	c.v.Printf("Run stats: %+v\n", stats)

	return c.Notify(stats)
}

func NewCommand() cli.NotifyCmd {
	cmd := &Command{
		flag: cli.Flags{flag.NewFlagSet("slack", flag.ExitOnError)},
		v:    vbs.New(),
		n:    new(slack.Notification),
	}

	cmd.flag.StringVar(&cmd.n.Token, "token", cmdDefault.Token, "Token (Required)")
	cmd.flag.StringVar(&cmd.n.Channel, "channel", cmdDefault.Channel, "Channel (Required)")
	cmd.flag.StringVar(&cmd.n.Text, "message", cmdDefault.Text, "Message")
	cmd.flag.StringVar(&cmd.n.Text, "m", cmdDefault.Text, "Message")
	cmd.flag.StringVar(&cmd.n.Parse, "parse", cmdDefault.Parse, "Parse")
	cmd.flag.IntVar(&cmd.n.LinkNames, "link-names", cmdDefault.LinkNames, "LinkNames")
	cmd.flag.BoolVar(&cmd.n.UnfurlLinks, "unfurl-links", cmdDefault.UnfurlLinks, "UnfurlLinks")
	cmd.flag.BoolVar(&cmd.n.UnfurlMedia, "unfurl-media", cmdDefault.UnfurlMedia, "UnfurlMedia")
	cmd.flag.StringVar(&cmd.n.Username, "username", cmdDefault.Username, "Username")
	cmd.flag.BoolVar(&cmd.n.AsUser, "as-user", cmdDefault.AsUser, "AsUser")
	cmd.flag.StringVar(&cmd.n.IconURL, "icon-url", cmdDefault.IconURL, "Username")
	cmd.flag.StringVar(&cmd.n.IconEmoji, "icon-emoji", cmdDefault.IconEmoji, "IconEmoji")

	cmd.flag.BoolVar(&cmd.v.Verbose, "verbose", false, "Enable verbose mode")
	cmd.flag.BoolVar(&cmd.help, "h", false, "Show help")
	cmd.flag.BoolVar(&cmd.help, "help", false, "Show help")

	return cmd
}
