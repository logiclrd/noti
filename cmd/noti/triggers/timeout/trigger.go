package timeout

import (
	"context"
	"errors"
	"io"
	"os"
	"time"

	"github.com/variadico/noti/cmd/noti/stats"
)

const FlagKey = "timeout"

type Trigger struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer

	stats stats.Info
	ctx   context.Context
	dur   time.Duration
}

func NewTrigger(ctx context.Context, s stats.Info, wait string) (*Trigger, error) {
	d, err := time.ParseDuration(wait)
	if err != nil {
		return nil, err
	}

	return &Trigger{
		stdin:  os.Stdin,
		stdout: os.Stdout,
		stderr: os.Stderr,
		stats:  s,
		ctx:    ctx,
		dur:    d,
	}, nil
}

func (t *Trigger) Streams() (io.Reader, io.Writer, io.Writer) {
	return t.stdin, t.stdout, t.stderr
}

func (t *Trigger) Run(cmdErr chan error, stats chan stats.Info) {
	start := time.Now()

	select {
	case <-cmdErr:
		return
	case <-t.ctx.Done():
		return
	case <-time.After(t.dur):
		t.stats.Err = errors.New("command timeout exceeded")
		t.stats.Duration = time.Since(start)
		out <- t.stats
	}
}
