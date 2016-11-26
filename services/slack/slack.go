package slack

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

const (
	// ParseNone autoformats the text in messages less.
	ParseNone = "none"
	// ParseFull autoformats message text more, like creating hyperlinks
	// automatically.
	ParseFull = "full"

	// LinkNamesOn enables making usernames hyperlinks.
	LinkNamesOn = 1
	// LinkNamesOn disables making usernames hyperlinks.
	LinkNamesOff = 0

	postAPI = "https://slack.com/api/chat.postMessage"
)

type apiResponse struct {
	OK        bool   `json:"ok"`
	Channel   string `json:"channel"`
	Timestamp string `json:"ts"`
	Message   struct {
		Text     string `json:"text"`
		Username string `json:"username"`
		Icons    struct {
			Emoji   string `json:"emoji"`
			Image64 string `json:"image_64"`
		} `json:"icons"`
		Type      string `json:"type"`
		Subtype   string `json:"subtype"`
		Timestamp string `json:"ts"`
	} `json:"message"`
	Error string `json:"error"`
}

type Notification struct {
	// Token is a user's authentication token.
	Token string
	// Channel is a notification's destination. It can be a channel, private
	// group, or username.
	Channel string
	// Text is the notification's message.
	Text string
	// Parse is the mode used to parse text.
	Parse string
	// LinkNames converts usernames into links.
	LinkNames int
	// Attachments are rich text snippets.
	Attachments map[string]string
	// UnfurlLinks attempts to expand a link to show a preview. Success depends
	// on the webpage having the right markdown.
	UnfurlLinks bool
	// UnfurlMedia attempts to expand a link to show a preview. Success depends
	// on the webpage having the right markdown.
	UnfurlMedia bool
	// Username given to bot. If AsUser is true, then message will try to be
	// sent from the given user.
	Username string
	// AsUser attempt to send a message as the user in Username.
	AsUser bool
	// IconURL is a URL to set as the user icon.
	IconURL string
	// IconEmoji is an emoji to set as the user icon.
	IconEmoji string

	client *http.Client
}

func (n *Notification) SetClient(c *http.Client) {
	n.client = c
}

func (n *Notification) Client() *http.Client {
	return n.client
}

// Notify sends a message request to the Slack API.
func (n *Notification) Send() error {
	if n.Token == "" {
		return errors.New("missing authentication token")
	}
	if n.Channel == "" {
		return errors.New("missing channel, group, or username destination")
	}
	if n.Text == "" {
		return errors.New("missing message text")
	}

	attach, err := json.Marshal(n.Attachments)
	if err != nil {
		return err
	}

	vals := make(url.Values)
	vals.Set("token", n.Token)
	vals.Set("channel", n.Channel)
	vals.Set("text", n.Text)
	vals.Set("parse", n.Parse)
	vals.Set("link_names", fmt.Sprint(n.LinkNames))
	vals.Set("attachments", string(attach))
	vals.Set("unfurl_links", fmt.Sprintf("%t", n.UnfurlLinks))
	vals.Set("unfurl_media", fmt.Sprintf("%t", n.UnfurlMedia))
	vals.Set("username", n.Username)
	vals.Set("as_user", fmt.Sprintf("%t", n.AsUser))
	vals.Set("icon_url", n.IconURL)
	vals.Set("icon_emoji", n.IconEmoji)

	// fmt.Println(vals)

	resp, err := n.client.PostForm(postAPI, vals)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}

	if !r.OK {
		return errors.New(r.Error)
	}

	return nil
}

func (n *Notification) SetMessage(m string) {
	n.Text = m
}
