package exit

import (
	"context"
	"io"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/variadico/noti/cmd/noti/stats"
)

const FlagKey = "exit"

type Trigger struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer

	stats stats.Info
	ctx   context.Context
}

func NewTrigger(ctx context.Context, s stats.Info) *Trigger {
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

func (t *Trigger) Run(cmdErr chan error, stats chan stats.Info) {
	start := time.Now()

	if t.stats.Cmd == "" {
		// User executed something like, "noti" or "noti banner", meaning
		// without a utility argument to run.
		out <- stats.Info{Cmd: "noti"}
		return
	}

	select {
	case err := <-cmdErr:
		if err != nil {
			t.stats.Err = err
			t.stats.ExitStatus = exitStatus(err)
		}
		t.stats.Duration = time.Since(start)
		out <- t.stats
	case <-t.ctx.Done():
		return
	}
}

func exitStatus(err error) int {
	eerr, is := err.(*exec.ExitError)
	if !is {
		return noExitStatus
	}

	if status, is := eerr.Sys().(syscall.WaitStatus); is {
		return status.ExitStatus()
	}

	return noExitStatus
}
