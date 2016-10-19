package cli

import (
	"flag"
	"io/ioutil"

	"github.com/variadico/noti/triggers"
)

type Flags struct {
	*flag.FlagSet

	Help    bool
	Verbose bool

	Triggers triggers.Flag
}

func NewFlags(name string) *Flags {
	set := flag.NewFlagSet(name, flag.ContinueOnError)
	set.Usage = func() {} // We handle this ourselves.
	set.SetOutput(ioutil.Discard)

	fs := new(Flags)
	fs.FlagSet = set

	fs.SetBools(&fs.Help, "h", "help", false)
	fs.SetBools(&fs.Verbose, "v", "verbose", false)

	fs.Var(&fs.Triggers, "trigger", "")

	return fs
}

func (fs *Flags) SetString(v *string, name, defaultVal string) {
	fs.StringVar(v, name, defaultVal, "")
}

func (fs *Flags) SetStrings(v *string, short, long, defaultVal string) {
	fs.StringVar(v, short, defaultVal, "")
	fs.StringVar(v, long, defaultVal, "")
}

func (fs *Flags) SetBool(v *bool, name string, defaultVal bool) {
	fs.BoolVar(v, name, defaultVal, "")
}

func (fs *Flags) SetBools(v *bool, short, long string, defaultVal bool) {
	fs.BoolVar(v, short, defaultVal, "")
	fs.BoolVar(v, long, defaultVal, "")
}

func (fs *Flags) SetInt(v *int, name string, defaultVal int) {
	fs.IntVar(v, name, defaultVal, "")
}

func (fs *Flags) SetInts(v *int, short, long string, defaultVal int) {
	fs.IntVar(v, short, defaultVal, "")
	fs.IntVar(v, long, defaultVal, "")
}

// Passed returns true if any of the given flags were passed by the user.
func (fs *Flags) Passed(names ...string) bool {
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
