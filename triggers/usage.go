package triggers

const Usage = `
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
        pid=<process id>
            Trigger a notification when a pid disappears. If the pid doesn't
            exist or this trigger isn't supported, then the notification will
            trigger immediately.
`
