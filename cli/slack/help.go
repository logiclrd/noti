package slack

var helpText = `noti slack [options] [command]

OPTIONS
    -m, -message
        Notification title. Default is utility name.

    -token
        Authentication token.
    -channel
        Channel, private group, or IM channel to send message to.
    -parse
        Change how messages are treated.
    -link-names
        Find and link channel names and usernames.
    -unfurl-links
        Pass true to enable unfurling of primarily text-based content.
    -unfurl-media
        Pass false to disable unfurling of media content.
    -username
        Set your bot's user name.
    -as-user
        Pass true to post the message as the authed user, instead of as a bot.
    -icon-url
        URL to an image to use as the icon for this message.
    -icon-emoji
        Emoji to use as the icon for this message.

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
