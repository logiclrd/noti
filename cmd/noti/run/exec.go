package run

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"syscall"
	"time"
)

func Exec(args ...string) Stats {
	if len(args) == 0 {
		return Stats{
			Cmd: "noti",
		}
	}

	sts := Stats{
		Cmd:      args[0],
		Args:     args[1:],
		ExitCode: noExitCode,
	}

	if _, err := exec.LookPath(args[0]); err != nil {
		// Before we run anything, we're going to check if we can find the
		// command. If we can't find a command, then we'll assume it might be
		// an aliased command.
		expanded, expErr := expandAlias(args[0])
		if expErr != nil {
			sts.ExitCode = cmdNotFound
			sts.Err = err
			return sts
		}

		args = append(expanded, args[1:]...)
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	start := time.Now()
	err := cmd.Run()
	sts.Duration = time.Since(start)

	if err == nil {
		return sts
	}

	sts.Err = err

	if eerr, is := err.(*exec.ExitError); is {
		if status, is := eerr.Sys().(syscall.WaitStatus); is {
			sts.ExitCode = status.ExitStatus()
		}
	}

	return sts
}

func ExecWithTimeout(d time.Duration, args ...string) Stats {
	if len(args) == 0 {
		return Stats{
			Cmd: "noti",
		}
	}

	sts := Stats{
		Cmd:      args[0],
		Args:     args[1:],
		ExitCode: -1,
	}

	if _, err := exec.LookPath(args[0]); err != nil {
		// Before we run anything, we're going to check if we can find the
		// command. If we can't find a command, then we'll assume it might be
		// an aliased command.
		expanded, expErr := expandAlias(args[0])
		if expErr != nil {
			sts.ExitCode = cmdNotFound
			sts.Err = err
			return sts
		}

		args = append(expanded, args[1:]...)
	}

	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()

	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	start := time.Now()
	err := cmd.Run()
	sts.Duration = time.Since(start)

	if cerr := ctx.Err(); cerr == context.DeadlineExceeded {
		sts.Err = errors.New("command timeout exceeded")
		return sts
	}

	if err == nil {
		return sts
	}

	sts.Err = err

	if eerr, is := err.(*exec.ExitError); is {
		if status, is := eerr.Sys().(syscall.WaitStatus); is {
			sts.ExitCode = status.ExitStatus()
		}
	}

	return sts
}
