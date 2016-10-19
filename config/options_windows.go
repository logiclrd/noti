package config

import "github.com/variadico/noti/ntypes/slack"

type Options struct {
	DefaultSet []string
	Banner     *notifyicon.Notification
	Slack      *slack.Notification
}
