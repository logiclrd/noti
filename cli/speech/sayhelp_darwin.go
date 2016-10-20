package speech

var helpText = `noti speech [options] [command]

OPTIONS
    -m, -message
        Notification message. Default is "Done!"

    -rate
        Rate of speech.
    -voice
        Voice used to speak. Default is "Alex".

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
