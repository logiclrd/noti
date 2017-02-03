// +build !darwin
// +build !windows

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

    -app-name
        App name.
    -replaces-id
        Replaces ID.
    -icon
        Icon.
    -timeout
        Timeout.
` + triggers.Usage + cli.GlobalUsage
