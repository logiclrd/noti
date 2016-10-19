// +build !darwin
// +build !windows

package config

import (
	"github.com/variadico/noti/espeak"
	"github.com/variadico/noti/freedesktop"
	"github.com/variadico/noti/slack"
)

type Options struct {
	DefaultSet []string
	Banner     *freedesktop.Notification
	Speech     *espeak.Notification
	Slack      *slack.Notification
}
