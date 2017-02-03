// +build !darwin
// +build !windows

package speech

import (
	"github.com/variadico/noti/cli"
	"github.com/variadico/noti/triggers"
)

var helpText = `noti speech [options] [command]

OPTIONS
    -m, -message
        Notification message. Default is "Done!"

    -rate
        Rate.
    -voice
        Voice name.
    -word-gap
        Word gap.
    -pitch
        Pitch.
` + triggers.Usage + cli.GlobalUsage
