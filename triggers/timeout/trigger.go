package timeout

import (
	"context"
	"errors"
	"time"

	"github.com/variadico/noti/cmd/noti/runstat"
)

const FlagKey = "timeout"

type Trigger struct {
	stats runstat.Result
	ctx   context.Context
	dur   time.Duration
}

func NewTrigger(ctx context.Context, s runstat.Result, d time.Duration) *Trigger {
	return &Trigger{
		stats: s,
		ctx:   ctx,
		dur:   d,
	}
}

func (t *Trigger) Run(cmdErr chan error, stats chan runstat.Result) {
	start := time.Now()

	select {
	case <-cmdErr:
		return
	case <-t.ctx.Done():
		return
	case <-time.After(t.dur):
		t.stats.Err = errors.New("command timeout exceeded")
		t.stats.Duration = time.Since(start)
		stats <- t.stats
	}
}
