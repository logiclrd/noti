package triggers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

type timeoutTrigger struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer

	stats Stats
	ctx   context.Context
	dur   time.Duration
}

func newTimeoutTrigger(ctx context.Context, s Stats, wait string) (*timeoutTrigger, error) {
	d, err := time.ParseDuration(wait)
	if err != nil {
		return nil, err
	}

	fmt.Println(">>>>>> CALLED TIMEOUT TRIGGER")

	return &timeoutTrigger{
		stdin:  os.Stdin,
		stdout: os.Stdout,
		stderr: os.Stderr,
		stats:  s,
		ctx:    ctx,
		dur:    d,
	}, nil
}

func (t *timeoutTrigger) streams() (io.Reader, io.Writer, io.Writer) {
	fmt.Println(">>>>>> CALLED TIMEOUT STREAMS")
	return t.stdin, t.stdout, t.stderr
}

func (t *timeoutTrigger) run(cmdErr chan error, out chan Stats) {
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
