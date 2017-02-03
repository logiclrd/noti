package root

import (
	"github.com/variadico/noti/cli"
	"github.com/variadico/noti/triggers"
)

var helpText = `noti [options] [notification [options]] [command]

OPTIONS` + triggers.Usage + cli.GlobalUsage + `
NOTIFICATIONS
    desktop
    slack
    speech
`
