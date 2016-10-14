package run

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"
	"time"
)

// Exec executes a command and waits for it to finish. When it does, it'll
// return statistics about the run.
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

// ExecWithTimeout executes a command. If the process takes longer than d, then
// it kills the process and returns statistics up to that point. Otherwise, the
// process runs like normal.
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

// ExecNotify executes a command. The function will return statistics if a process takes longer
// than d. However, it will continue executing the process.
func ExecNotify(ctx context.Context, args ...string) chan Stats {
	out := make(chan Stats)
	go execNotify(ctx, out, args)
	return out
}

func execNotify(ctx context.Context, out chan Stats, args []string) {
	defer close(out)

	if len(args) == 0 {
		out <- Stats{
			Cmd: "noti",
		}
		return
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
			out <- sts
			return
		}

		args = append(expanded, args[1:]...)
	}

	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	start := time.Now()
	errc := make(chan error)
	go func() { errc <- cmd.Run() }()

	t := time.NewTicker(2 * time.Second)

	fmt.Println(">>>>>>>> WAIT LOOP!")
	for {
		select {
		case <-ctx.Done():
			fmt.Println(">>>>>>>> CONTEXT CANCELLED!")
			sts.Duration = time.Since(start)
			out <- sts
			return
		case err := <-errc:
			fmt.Println(">>>>>>>> COMMAND FINISHED!")
			sts.Duration = time.Since(start)
			sts.Err = err
			if eerr, is := err.(*exec.ExitError); is {
				if status, is := eerr.Sys().(syscall.WaitStatus); is {
					sts.ExitCode = status.ExitStatus()
				}
			}
			sts.State = "done"
			out <- sts
			return
		case <-t.C:
			sts.Duration = time.Since(start)
			sts.State = "running"
			out <- sts
		}
	}
}

// ExecContains
func ExecContains(ctx context.Context, target string, args ...string) chan Stats {
	out := make(chan Stats)
	go execContains(ctx, out, target, args)
	return out
}

type scanWriter struct {
	target []byte
	found  bool
	out    io.Writer
}

func (s *scanWriter) Write(p []byte) (int, error) {
	s.found = bytes.Contains(p, s.target)
	return fmt.Fprint(s.out, string(p))
}

func execContains(ctx context.Context, out chan Stats, target string, args []string) {
	defer close(out)

	if len(args) == 0 {
		out <- Stats{
			Cmd: "noti",
		}
		return
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
			out <- sts
			return
		}

		args = append(expanded, args[1:]...)
	}

	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	cmd.Stdin = os.Stdin

	scanStdout := &scanWriter{out: os.Stdout, target: []byte(target)}
	scanStderr := &scanWriter{out: os.Stderr, target: []byte(target)}
	cmd.Stdout = scanStdout
	cmd.Stderr = scanStderr

	start := time.Now()
	errc := make(chan error)
	go func() { errc <- cmd.Run() }()

	fmt.Println(">>>>>>>> WAIT LOOP!")
	for {
		select {
		case <-ctx.Done():
			fmt.Println(">>>>>>>> CONTEXT CANCELLED!")
			sts.Duration = time.Since(start)
			out <- sts
			return
		case err := <-errc:
			fmt.Println(">>>>>>>> COMMAND FINISHED!")
			sts.Duration = time.Since(start)
			sts.Err = err
			if eerr, is := err.(*exec.ExitError); is {
				if status, is := eerr.Sys().(syscall.WaitStatus); is {
					sts.ExitCode = status.ExitStatus()
				}
			}
			sts.State = "done"
			out <- sts
			return
		default:
			sts.Duration = time.Since(start)
			sts.State = "running"

			if scanStdout.found || scanStderr.found {
				out <- sts
				scanStdout.found = false
				scanStderr.found = false
			}
		}
	}
}
