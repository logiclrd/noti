package timeout

import (
	"context"
	"errors"
	"time"

	"github.com/variadico/noti/cmd/noti/run"
)

const FlagKey = "timeout"

type Trigger struct {
	stats run.Stats
	ctx   context.Context
	dur   time.Duration
}

func NewTrigger(ctx context.Context, s run.Stats, d time.Duration) *Trigger {
	return &Trigger{
		stats: s,
		ctx:   ctx,
		dur:   d,
	}
}

func (t *Trigger) Run(cmdErr chan error, stats chan run.Stats) {
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
