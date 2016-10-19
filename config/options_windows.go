package config

import (
	"github.com/variadico/noti/notifyicon"
	"github.com/variadico/noti/slack"
)

type Options struct {
	DefaultSet []string
	Banner     *notifyicon.Notification
	Slack      *slack.Notification
}
