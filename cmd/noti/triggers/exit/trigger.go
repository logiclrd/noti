package exit

import (
	"context"
	"io"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/variadico/noti/cmd/noti/run"
)

const FlagKey = "exit"

type Trigger struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer

	stats run.Stats
	ctx   context.Context
}

func NewTrigger(ctx context.Context, s run.Stats) *Trigger {
	return &Trigger{
		stdin:  os.Stdin,
		stdout: os.Stdout,
		stderr: os.Stderr,
		stats:  s,
		ctx:    ctx,
	}
}

func (t *Trigger) Streams() (io.Reader, io.Writer, io.Writer) {
	return t.stdin, t.stdout, t.stderr
}

func (t *Trigger) Run(cmdErr chan error, stats chan run.Stats) {
	start := time.Now()

	if t.stats.Cmd == "" {
		// User executed something like, "noti" or "noti banner", meaning
		// without a utility argument to run.
		stats <- run.Stats{Cmd: "noti"}
		return
	}

	select {
	case err := <-cmdErr:
		if err != nil {
			t.stats.Err = err
			t.stats.ExitStatus = exitStatus(err)
		}
		t.stats.Duration = time.Since(start)
		stats <- t.stats
	case <-t.ctx.Done():
		return
	}
}

func exitStatus(err error) int {
	eerr, is := err.(*exec.ExitError)
	if !is {
		return run.NoExitStatus
	}

	if status, is := eerr.Sys().(syscall.WaitStatus); is {
		return status.ExitStatus()
	}

	return run.NoExitStatus
}
