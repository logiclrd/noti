package config

import (
	"github.com/variadico/noti/services/nsuser"
	"github.com/variadico/noti/services/say"
	"github.com/variadico/noti/services/slack"
)

type Options struct {
	DefaultNotifications []string
	DefaultTriggers      []string
	Desktop              *nsuser.Notification
	Speech               *say.Notification
	Slack                *slack.Notification
}

func NewOptions() Options {
	return Options{
		DefaultNotifications: make([]string, 0),
		DefaultTriggers:      make([]string, 0),
		Desktop:              new(nsuser.Notification),
		Speech:               new(say.Notification),
		Slack:                new(slack.Notification),
	}
}
