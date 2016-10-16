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

type trigger interface {
	streams() (stdin io.Reader, stdout io.Writer, stderr io.Writer)
	run(chan error, chan Stats)
}

type onExitTrigger struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer

	stats Stats
	ctx   context.Context
}

func newOnExitTrigger(ctx context.Context, s Stats) *onExitTrigger {
	return &onExitTrigger{
		stdin:  os.Stdin,
		stdout: os.Stdout,
		stderr: os.Stderr,
		stats:  s,
		ctx:    ctx,
	}
}

func (t *onExitTrigger) streams() (io.Reader, io.Writer, io.Writer) {
	return t.stdin, t.stdout, t.stderr
}

func (t *onExitTrigger) run(cmdErr chan error, out chan Stats) {
	fmt.Println(">>> onExit")
	defer fmt.Println(">>> end onExit")

	start := time.Now()

	if t.stats.Cmd == "" {
		// User executed something like, "noti" or "noti banner".
		out <- Stats{Cmd: "noti"}
		return
	}

	fmt.Println("SELECT!!!!")
	select {
	case err := <-cmdErr:
		fmt.Println("PULLED ERR!!!!")
		if err != nil {
			t.stats.Err = err
			t.stats.ExitCode = exitCode(err)
		}
		t.stats.Duration = time.Since(start)
		out <- t.stats
		fmt.Println("SENT STATS!!!!")
	case <-t.ctx.Done():
		fmt.Println("context cancelled!!!!")
		return
	}
}

type onTimeoutTrigger struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer

	stats Stats
	ctx   context.Context
	dur   time.Duration
}

func newOnTimeoutTrigger(ctx context.Context, s Stats, wait string) (*onTimeoutTrigger, error) {
	d, err := time.ParseDuration(wait)
	if err != nil {
		return nil, err
	}

	fmt.Println(">>>>>> CALLED TIMEOUT TRIGGER")

	return &onTimeoutTrigger{
		stdin:  os.Stdin,
		stdout: os.Stdout,
		stderr: os.Stderr,
		stats:  s,
		ctx:    ctx,
		dur:    d,
	}, nil
}

func (t *onTimeoutTrigger) streams() (io.Reader, io.Writer, io.Writer) {
	fmt.Println(">>>>>> CALLED TIMEOUT STREAMS")
	return t.stdin, t.stdout, t.stderr
}

func (t *onTimeoutTrigger) run(cmdErr chan error, out chan Stats) {
	fmt.Println(">>> CALLED TIMEOUT TRIGGER RUN")
	defer fmt.Println(">>> end onTimeout")
	start := time.Now()

	select {
	case <-cmdErr:
		return
	case <-time.After(t.dur):
		t.stats.Err = errors.New("command timeout exceeded")
		t.stats.Duration = time.Since(start)
		out <- t.stats
	case <-t.ctx.Done():
		return
	}
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

type onContainsTrigger struct {
	stdin  io.Reader
	stdout *scanWriter
	stderr *scanWriter

	stats  Stats
	ctx    context.Context
	target string
}

func newOnContainsTrigger(ctx context.Context, s Stats, t string) *onContainsTrigger {
	scanStdout := &scanWriter{out: os.Stdout, target: []byte(t)}
	scanStderr := &scanWriter{out: os.Stderr, target: []byte(t)}

	return &onContainsTrigger{
		stdin:  os.Stdin,
		stdout: scanStdout,
		stderr: scanStderr,
		stats:  s,
		ctx:    ctx,
		target: t,
	}
}

func (t *onContainsTrigger) streams() (io.Reader, io.Writer, io.Writer) {
	return t.stdin, t.stdout, t.stderr
}

func (t *onContainsTrigger) run(cmdErr chan error, out chan Stats) {
	fmt.Println(">>> onContains")
	defer fmt.Println(">>> end onContains")
	start := time.Now()

	for {
		select {
		case <-t.ctx.Done():
			return
		case <-cmdErr:
			return
		default:
			t.stats.Duration = time.Since(start)
			t.stats.State = "running"

			if t.stdout.found || t.stderr.found {
				out <- t.stats
				t.stdout.found = false
				t.stderr.found = false
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
