package triggers

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"
)

type exitTrigger struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer

	stats Stats
	ctx   context.Context
}

func newExitTrigger(ctx context.Context, s Stats) *exitTrigger {
	return &exitTrigger{
		stdin:  os.Stdin,
		stdout: os.Stdout,
		stderr: os.Stderr,
		stats:  s,
		ctx:    ctx,
	}
}

func (t *exitTrigger) streams() (io.Reader, io.Writer, io.Writer) {
	return t.stdin, t.stdout, t.stderr
}

func (t *exitTrigger) run(cmdErr chan error, out chan Stats) {
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
		fmt.Println("PULLED ERR!!!!", err)
		if err != nil {
			t.stats.Err = err
			t.stats.ExitStatus = exitStatus(err)
		}
		t.stats.Duration = time.Since(start)
		out <- t.stats
		fmt.Println("SENT STATS!!!!")
	case <-t.ctx.Done():
		fmt.Println("exit cancelled!!!!")
		return
	}
}
