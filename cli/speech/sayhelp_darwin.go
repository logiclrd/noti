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
        Rate of speech.
    -voice
        Voice used to speak. Default is "Alex".
` + triggers.Usage + cli.GlobalUsage
