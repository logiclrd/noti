package triggers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"time"
)

type scanWriter struct {
	target []byte
	found  bool
}

func (s *scanWriter) Write(p []byte) (int, error) {
	s.found = bytes.Contains(p, s.target)
	return len(p), nil
}

type matchTrigger struct {
	stdin  io.Reader
	stdout *scanWriter
	stderr *scanWriter

	stats  Stats
	ctx    context.Context
	target string
}

func newMatchTrigger(ctx context.Context, s Stats, t string) *matchTrigger {
	scanStdout := &scanWriter{target: []byte(t)}
	scanStderr := &scanWriter{target: []byte(t)}

	return &matchTrigger{
		stdin:  os.Stdin,
		stdout: scanStdout,
		stderr: scanStderr,
		stats:  s,
		ctx:    ctx,
		target: t,
	}
}

func (t *matchTrigger) streams() (io.Reader, io.Writer, io.Writer) {
	return t.stdin, t.stdout, t.stderr
}

func (t *matchTrigger) run(cmdErr chan error, out chan Stats) {
	fmt.Println(">>> onContains")
	defer fmt.Println(">>> end onContains")
	start := time.Now()

	fmt.Println("starting scan loop")
	for {
		select {
		case <-t.ctx.Done():
			fmt.Println("contains cancelled!")
			return
		case <-cmdErr:
			fmt.Println("contains pulled error!")
			return
		default:
			t.stats.Duration = time.Since(start)
			t.stats.State = "running"

			if t.stdout.found || t.stderr.found {
				out <- t.stats
				t.stdout.found = false
				t.stderr.found = false
			}
		}
	}
}
