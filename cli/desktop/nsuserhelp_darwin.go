package desktop

import (
	"github.com/variadico/noti/cli"
	"github.com/variadico/noti/triggers"
)

var helpText = `noti desktop [options] [command]

OPTIONS
    -t, -title
        Notification title. Default is utility name.
    -m, -message
        Notification message. Default is "Done!"

    -subtitle
        Notification subtitle.
    -icon
        Notification icon.
    -sound
        Notification sound.
` + triggers.Usage + cli.GlobalUsage
