// +build !darwin
// +build !windows

package banner

var helpText = `noti banner [options] [command]

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
