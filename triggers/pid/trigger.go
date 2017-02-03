package pid

import (
	"context"

	"github.com/variadico/noti/runstat"
)

const FlagKey = "pid"

type Trigger struct {
	stats runstat.Result
	ctx   context.Context
	pid   int
}

func NewTrigger(ctx context.Context, s runstat.Result, pid int) *Trigger {
	return &Trigger{
		stats: s,
		ctx:   ctx,
		pid:   pid,
	}
}
