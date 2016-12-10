package config

import "github.com/variadico/noti/services/slack"

type Options struct {
	DefaultSet []string
	Desktop     *notifyicon.Notification
	Slack      *slack.Notification
}
