package say

import (
	"fmt"
	"os"
	"os/exec"
)

type Notification struct {
	Text string
	// Voice is the name of the voice to speak the notification.
	Voice string
	// Rate controls how fast the voice will speak. It's measured in words per
	// minute.
	Rate int
}

func (n *Notification) Send() error {
	r := fmt.Sprint(n.Rate)
	cmd := exec.Command("say", "-v", n.Voice, "-r", r, n.Text)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (n *Notification) SetMessage(m string) {
	n.Text = m
}

func (n *Notification) Message() string {
	return n.Text
}
