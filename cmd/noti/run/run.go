package run

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

const (
	cmdNotFound = 127
)

type Stats struct {
	Cmd      string
	Args     []string
	Stdout   string
	Stderr   string
	ExitCode int
	ExecErr  error
	Duration time.Duration
}

func Exec(args ...string) Stats {
	if len(args) == 0 {
		return Stats{
			Cmd: "noti",
		}
	}

	st := Stats{
		Cmd:  args[0],
		Args: args[1:],
	}

	if _, err := exec.LookPath(args[0]); err != nil {
		// Before we run anything, we're going to check if we can find the
		// command. If we can't find a command, then we'll assume it might be
		// an aliased command.
		expanded, expErr := expandAlias(args[0])
		if expErr != nil {
			st.ExitCode = cmdNotFound
			st.ExecErr = err
			return st
		}

		args = append(expanded, args[1:]...)
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr

	// Should this be conditional? This could grow large if noti is watching a
	// long running process or one that outputs a lot.
	out2 := new(bytes.Buffer)
	err2 := new(bytes.Buffer)
	cmd.Stdout = io.MultiWriter(os.Stdout, out2)
	cmd.Stderr = io.MultiWriter(os.Stderr, err2)

	start := time.Now()
	if err := cmd.Run(); err != nil {
		st.Duration = time.Since(start)
		st.ExecErr = err

		if eerr, is := err.(*exec.ExitError); is {
			if status, is := eerr.Sys().(syscall.WaitStatus); is {
				st.ExitCode = status.ExitStatus()
			}
		}

		return st
	}
	st.Duration = time.Since(start)
	st.Stdout = out2.String()
	st.Stderr = err2.String()

	return st
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
