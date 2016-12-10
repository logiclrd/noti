package speech

import (
	"github.com/variadico/noti/cli"
	"github.com/variadico/noti/runstat"
)

type Command struct {
	//
}

func (c *Command) Parse(args []string) error {
	return nil
}

func (c *Command) Notify(stats runstat.Result) error {
	return nil
}

func (c *Command) Run() error {
	return nil
}

func NewCommand() cli.NotifyCmd {
	return &Command{}
}
