package espeak

import (
	"os"
	"os/exec"
)

type Notification struct {
	// -g
	WordGap int
	// -p
	PitchAdjustment int
	// -s
	Rate int
	// -v
	VoiceName string

	Text string
}

func (n *Notification) Send() error {
	cmd := exec.Command("espeak", "-v", n.VoiceName, "--", n.Text)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (n *Notification) SetMessage(m string) {
	n.Text = m
}

func (n *Notification) Message() string {
	return n.Text
}
