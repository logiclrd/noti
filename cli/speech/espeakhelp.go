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

    -trigger 'mode[=value]'
        Set notification trigger. You can set this flag multiple times to set
        multiple triggers.

        exit
            Trigger a notification when the process exits. Default.
        match=<string>
            Trigger a notification when the running command prints a string to
            stdout or stderr.
        timeout=<duration>
            Trigger a notification and kill the running command after a certain
            duration.

    -v, -verbose
        Enable verbose mode.
    -h, -help
        Show help.
`
