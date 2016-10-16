package cli

import "github.com/variadico/noti/cmd/noti/triggers"

type Cmd interface {
	Run() error
	Parse(args []string) error
}

type NotifyCmd interface {
	Cmd
	Notify(triggers.Stats) error
}
