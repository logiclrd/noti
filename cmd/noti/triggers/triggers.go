package triggers

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

var notiStats = Stats{Cmd: "noti"}

func onExit(ctx context.Context, args []string, out chan Stats) {
	fmt.Println(">>> onExit")

	sts := statsFromArgs(args)
	if sts.Cmd == "" {
		// User executed something like, "noti" or "noti banner".
		out <- notiStats
		return
	}

	var cmd *exec.Cmd
	if len(sts.ExpandedAlias) == 0 {
		cmd = exec.CommandContext(ctx, sts.Cmd, sts.Args...)
	} else {
		cmd = exec.CommandContext(ctx, sts.ExpandedAlias[0], sts.ExpandedAlias[1:]...)
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	start := time.Now()
	err := cmd.Run()
	sts.Duration = time.Since(start)

	if err != nil {
		sts.Err = err
		sts.ExitCode = exitCode(err)
	}

	out <- sts
}

func onTimeout(ctx context.Context, wait string, args []string, out chan Stats) {
	fmt.Println(">>> onTimeout")

	sts := statsFromArgs(args)
	if sts.Cmd == "" {
		return
	}

	dur, err := time.ParseDuration(wait)
	if err != nil {
		sts.Err = err
		out <- sts
		return
	}

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, dur)
	defer cancel()

	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	start := time.Now()
	err = cmd.Run()
	sts.Duration = time.Since(start)

	if err := ctx.Err(); err == context.DeadlineExceeded {
		sts.Err = errors.New("command timeout exceeded")
		out <- sts
		return
	}

	if err != nil {
		sts.Err = err
		sts.ExitCode = exitCode(err)
	}

	out <- sts
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

func onContains(ctx context.Context, target string, args []string, out chan Stats) {
	sts := statsFromArgs(args)
	if sts.Cmd == "" {
		return
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

	for {
		select {
		case <-ctx.Done():
			fmt.Println(">>>>>>>> CONTEXT CANCELLED!")
			sts.Duration = time.Since(start)
			out <- sts
			return
		case <-errc:
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

func exitCode(err error) int {
	eerr, is := err.(*exec.ExitError)
	if !is {
		return noExitCode
	}

	if status, is := eerr.Sys().(syscall.WaitStatus); is {
		return status.ExitStatus()
	}

	return noExitCode
}
