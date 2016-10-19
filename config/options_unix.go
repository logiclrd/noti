// +build !darwin
// +build !windows

package config

import (
	"github.com/variadico/noti/ntypes/espeak"
	"github.com/variadico/noti/ntypes/slack"
)

type Options struct {
	DefaultSet []string
	Banner     *freedesktop.Notification
	Speech     *espeak.Notification
	Slack      *slack.Notification
}
