package exit

import (
	"context"
	"os/exec"
	"syscall"
	"time"

	"github.com/variadico/noti/runstat"
)

const FlagKey = "exit"

type Trigger struct {
	stats runstat.Result
	ctx   context.Context
}

func NewTrigger(ctx context.Context, s runstat.Result) *Trigger {
	return &Trigger{
		stats: s,
		ctx:   ctx,
	}
}

func (t *Trigger) Run(cmdErr chan error, stats chan runstat.Result) {
	start := time.Now()

	if t.stats.Cmd == "" {
		// User executed something like, "noti" or "noti banner", meaning
		// without a utility argument to run.
		stats <- runstat.Result{Cmd: "noti"}
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
		return runstat.NoExitStatus
	}

	if status, is := eerr.Sys().(syscall.WaitStatus); is {
		return status.ExitStatus()
	}

	return runstat.NoExitStatus
}
