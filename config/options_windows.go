package config

import (
	"github.com/variadico/noti/services/notifyicon"
	"github.com/variadico/noti/services/slack"
	"github.com/variadico/noti/services/speechsynthesizer"
)

type Options struct {
	DefaultNotifications []string
	DefaultTriggers      []string
	Desktop              *notifyicon.Notification
	Speech               *speechsynthesizer.Notification
	Slack                *slack.Notification
}

func NewOptions() Options {
	return Options{
		DefaultNotifications: make([]string, 0),
		DefaultTriggers:      make([]string, 0),
		Desktop:              new(notifyicon.Notification),
		Speech:               new(speechsynthesizer.Notification),
		Slack:                new(slack.Notification),
	}
}
