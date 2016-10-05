package cli

import "github.com/variadico/noti/cmd/noti/run"

type Cmd interface {
	Run() error
	Parse(args []string) error
}

type NotifyCmd interface {
	Cmd
	Notify(run.Stats) error
}
