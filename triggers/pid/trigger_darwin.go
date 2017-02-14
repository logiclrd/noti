package pid

import (
	"os"
	"syscall"
	"time"

	"github.com/variadico/noti/runstat"
)

func (t *Trigger) Run(cmdErr chan error, stats chan runstat.Result) {
	proc, err := os.FindProcess(t.pid)
	if err != nil {
		t.stats.Err = err
		stats <- t.stats
		return
	}

	for {
		select {
		case <-t.ctx.Done():
			return
		default:
			// Check if PID exists.

			// If sending sig 0 works, then process is still alive.
			err = proc.Signal(syscall.Signal(0))
			if err != nil {
				t.stats.Err = err
				stats <- t.stats
				return
			}

			time.Sleep(2 * time.Second)
		}
	}
}
