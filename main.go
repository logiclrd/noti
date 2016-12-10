package main

import (
	"log"
	"os"

	"github.com/variadico/noti/cli"
	"github.com/variadico/noti/cli/desktop"
	"github.com/variadico/noti/cli/root"
	"github.com/variadico/noti/cli/slack"
	"github.com/variadico/noti/cli/speech"
	"github.com/variadico/noti/cli/version"
	"github.com/variadico/noti/config"
	"github.com/variadico/yaml"
)

func main() {
	log.SetFlags(0)

	// Don't want user to run command, but fail to notify because of a syntax
	// error. Tell users now, so they don't wait.
	if _, err := config.File(); err != nil {
		if yerr, is := err.(*yaml.TypeError); is {
			log.Fatalln("Config file error:", yerr)
		}
	}

	noti := root.NewCommand().(*root.Command)
	if err := noti.Parse(os.Args[1:]); err != nil {
		log.Println("Error:", err)
		log.Fatalln("Try 'noti -help' for more information.")
	}

	noti.Cmds = map[string]cli.Cmd{
		"version": version.NewCommand(),
		"desktop": desktop.NewCommand().(cli.Cmd),
		"speech":  speech.NewCommand().(cli.Cmd),
		"slack":   slack.NewCommand().(cli.Cmd),
	}

	if len(noti.Args()) == 0 {
		// noti was called by itself.
		if err := noti.Run(); err != nil {
			log.Fatalln("Error:", err)
		}
		return
	}

	var cmd cli.Cmd
	var found bool
	cmd, found = noti.Cmds[noti.Args()[0]]

	if found {
		// Command is something like: "noti foo ls"
		if err := cmd.Parse(noti.Args()[1:]); err != nil {
			log.Println("Error:", err)
			log.Fatalf("Try 'noti %s -help' for more information.", noti.Args()[0])
		}
	} else {
		// Command is something like: noti ls
		cmd = noti
	}

	if err := cmd.Run(); err != nil {
		log.Fatalln("Error:", err)
	}
}
