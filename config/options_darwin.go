package config

import (
	"github.com/variadico/noti/ntypes/nsuser"
	"github.com/variadico/noti/ntypes/say"
	"github.com/variadico/noti/ntypes/slack"
)

type Options struct {
	DefaultNotifications []string
	DefaultTriggers []string
	Banner     *nsuser.Notification
	Speech     *say.Notification
	Slack      *slack.Notification
}

func NewOptions() Options {
	return Options{
		DefaultNotifications: make([]string, 0),
		DefaultTriggers: make([]string, 0),
		Banner:     new(nsuser.Notification),
		Speech:     new(say.Notification),
		Slack:      new(slack.Notification),
	}
}
