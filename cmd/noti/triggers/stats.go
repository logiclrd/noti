package triggers

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	cmdNotFound = 127
	noExitCode  = -1
)

type Stats struct {
	Cmd  string
	Args []string
	// Stdout   string
	// Stderr   string
	ExitCode      int
	Err           error
	Duration      time.Duration
	State         string
	ExpandedAlias []string
}

func statsFromArgs(a []string) Stats {
	if len(a) == 0 {
		return Stats{}
	}

	sts := Stats{
		Cmd:      a[0],
		Args:     a[1:],
		ExitCode: noExitCode,
	}

	if _, err := exec.LookPath(a[0]); err != nil {
		// Before we run anything, we're going to check if we can find the
		// command. If we can't find a command, then we'll assume it might be
		// an aliased command.
		expanded, expErr := expandAlias(a[0])
		if expErr != nil {
			sts.ExitCode = cmdNotFound
			sts.Err = err
			return sts
		}

		// The user command could have been something like:
		// gss --foo
		// Put the expanded form first, then the args.
		sts.ExpandedAlias = append(expanded, a[1:]...)
	}

	return sts
}

// expandAlias attempts to expand an alias and return back the real command.
// Another way of executing an alias might be to directly execute the alias in
// the subshell, instead of expanding it and returning back to the supershell.
// Currently, that requires the user to do more escaping, which we want to
// avoid. That's why we're doing it this way instead.
// This has only been tested on ZSH and Bash.
func expandAlias(a string) ([]string, error) {
	shell := os.Getenv("SHELL")

	cmd := exec.Command(shell, "-l", "-i", "-c", "which "+a)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	exp := parseExpansion(string(out), a)
	return strings.Split(exp, " "), nil
}

func parseExpansion(s, alias string) string {
	s = strings.TrimSpace(s)

	prefix := fmt.Sprintf("%s: aliased to ", alias)
	start := strings.Index(s, prefix)

	if start == -1 {
		return ""
	}

	s = s[start:]
	s = s[len(prefix):]

	return s
}
