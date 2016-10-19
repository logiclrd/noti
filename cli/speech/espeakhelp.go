// +build !darwin
// +build !windows

package speech

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

    -v, -verbose
        Enable verbose mode.
    -h, -help
        Show help.
`
