package match

import (
	"bytes"
	"context"
	"io"
	"os"
	"time"

	"github.com/variadico/noti/runstat"
)

const FlagKey = "match"

type scanWriter struct {
	target []byte
	found  bool
}

func (s *scanWriter) Write(p []byte) (int, error) {
	s.found = bytes.Contains(p, s.target)
	return len(p), nil
}

type Trigger struct {
	stdin  io.Reader
	stdout *scanWriter
	stderr *scanWriter

	stats  runstat.Result
	ctx    context.Context
	target string
}

func NewTrigger(ctx context.Context, s runstat.Result, target string) *Trigger {
	scanStdout := &scanWriter{target: []byte(target)}
	scanStderr := &scanWriter{target: []byte(target)}

	return &Trigger{
		stdin:  os.Stdin,
		stdout: scanStdout,
		stderr: scanStderr,
		stats:  s,
		ctx:    ctx,
		target: target,
	}
}

func (t *Trigger) Streams() (io.Reader, io.Writer, io.Writer) {
	return t.stdin, t.stdout, t.stderr
}

func (t *Trigger) Run(cmdErr chan error, stats chan runstat.Result) {
	start := time.Now()

	for {
		select {
		case <-t.ctx.Done():
			return
		case <-cmdErr:
			return
		default:
			t.stats.Duration = time.Since(start)
			t.stats.State = "running"

			if t.stdout.found || t.stderr.found {
				stats <- t.stats
				t.stdout.found = false
				t.stderr.found = false
			}
		}
	}
}
