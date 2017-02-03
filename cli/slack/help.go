package slack

import (
	"github.com/variadico/noti/cli"
	"github.com/variadico/noti/triggers"
)

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
` + triggers.Usage + cli.GlobalUsage
