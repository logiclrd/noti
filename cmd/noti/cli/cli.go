package cli

import "github.com/variadico/noti/cmd/noti/runstat"

type Cmd interface {
	Run() error
	Parse(args []string) error
}

type NotifyCmd interface {
	Cmd
	Notify(runstat.Result) error
}
