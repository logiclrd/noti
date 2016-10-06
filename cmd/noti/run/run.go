package run

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
	ExitCode int
	Err      error
	Duration time.Duration
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
	e, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	exp := strings.TrimSpace(string(e))
	trimLen := fmt.Sprintf("%s: aliased to ", a)
	exp = exp[len(trimLen):]

	return strings.Split(exp, " "), nil
}
