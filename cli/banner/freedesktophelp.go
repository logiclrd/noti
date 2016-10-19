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

    -v, -verbose
        Enable verbose mode.
    -h, -help
        Show help.
`
