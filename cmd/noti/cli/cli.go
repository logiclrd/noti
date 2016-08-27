package cli

import (
	"flag"

	"github.com/variadico/noti/cmd/noti/run"
)

type Cmd interface {
	Run() error
	Parse(args []string) error
}

type NotifyCmd interface {
	Cmd
	Notify(run.Stats) error
}

type Flags struct {
	*flag.FlagSet
}

// Set returns true if any of the given flags were passed by the user.
func (fs Flags) Set(names ...string) bool {
	var wasPassed bool

	for _, n := range names {
		fs.Visit(func(f *flag.Flag) {
			if f.Name == n {
				wasPassed = true
			}
		})
		if wasPassed {
			return true
		}
	}

	return false
}
