package triggers

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/variadico/noti/config"
	"github.com/variadico/noti/runstat"
	"github.com/variadico/noti/triggers/exit"
	"github.com/variadico/noti/triggers/match"
	"github.com/variadico/noti/triggers/pid"
	"github.com/variadico/noti/triggers/timeout"
)

const (
	delim = "="
)

func Run(trigFlags []string, args []string, notify func(runstat.Result) error) error {
	if conf, err := config.File(); err == nil {
		if len(trigFlags) == 0 {
			trigFlags = append(trigFlags, conf.DefaultTriggers...)
		}
	}

	if len(trigFlags) == 0 {
		trigFlags = append(trigFlags, exit.FlagKey)
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Typically, the context will get cancelled by the first trigger to
	// finish. It'll also get called if we hit an error and return early. Or
	// return normally. It's okay if this gets called multiple times.
	defer cancel()

	sts := runstat.NewResult(args)

	trigs := make([]Trigger, 0, len(trigFlags))
	for _, t := range trigFlags {
		name, val := keyValue(t)

		switch name {
		case exit.FlagKey:
			trigs = append(trigs, exit.NewTrigger(ctx, sts))
		case match.FlagKey:
			trigs = append(trigs, match.NewTrigger(ctx, sts, val))
		case timeout.FlagKey:
			d, err := time.ParseDuration(val)
			if err != nil {
				return fmt.Errorf("timeout trigger: %s", err)
			}

			trigs = append(trigs, timeout.NewTrigger(ctx, sts, d))
		case pid.FlagKey:
			id, err := strconv.Atoi(val)
			if err != nil {
				return fmt.Errorf("pid trigger: %s", err)
			}

			trigs = append(trigs, pid.NewTrigger(ctx, sts, id))
		default:
			return fmt.Errorf("unknown trigger: %s", name)
		}
	}

	cmd := newCmd(ctx, sts, trigs)
	errc := make(chan error)
	go func() { errc <- cmd.Run() }()

	var wg sync.WaitGroup
	statsc := make(chan runstat.Result)

	for _, t := range trigs {
		wg.Add(1)
		go func(t Trigger) {
			t.Run(errc, statsc)
			wg.Done()
			cancel()
		}(t)
	}

	go func() { wg.Wait(); close(statsc) }()

	for stat := range statsc {
		if err := notify(stat); err != nil {
			return err
		}
	}

	return nil
}

func keyValue(s string) (string, string) {
	i := strings.Index(s, delim)
	if i == -1 {
		// Trigger is something like, "exit".
		return s, ""
	}

	// Trigger is something like, "contains=hello world".
	return s[:i], s[i+1:]
}

func uniqStreams(ts []Trigger) (stdin io.Reader, stdout io.Writer, stderr io.Writer) {
	inmap := map[io.Reader]struct{}{os.Stdin: struct{}{}}
	outmap := map[io.Writer]struct{}{os.Stdout: struct{}{}}
	errmap := map[io.Writer]struct{}{os.Stderr: struct{}{}}

	// Make streams unique.
	for _, t := range ts {
		s, is := t.(Streamer)
		if !is {
			continue
		}

		sin, sout, serr := s.Streams()
		inmap[sin] = struct{}{}
		outmap[sout] = struct{}{}
		errmap[serr] = struct{}{}
	}

	stdins := make([]io.Reader, 0, 1)
	stdouts := make([]io.Writer, 0, 1)
	stderrs := make([]io.Writer, 0, 1)

	// Convert to slices.
	for s := range inmap {
		stdins = append(stdins, s)
	}
	for s := range outmap {
		stdouts = append(stdouts, s)
	}
	for s := range errmap {
		stderrs = append(stderrs, s)
	}

	return io.MultiReader(stdins...), io.MultiWriter(stdouts...), io.MultiWriter(stderrs...)
}

func newCmd(ctx context.Context, sts runstat.Result, ts []Trigger) *exec.Cmd {
	var cmd *exec.Cmd

	if len(sts.ExpandedAlias) == 0 {
		cmd = exec.CommandContext(ctx, sts.Cmd, sts.Args...)
	} else {
		cmd = exec.CommandContext(ctx, sts.ExpandedAlias[0], sts.ExpandedAlias[1:]...)
	}

	_, sout, serr := uniqStreams(ts)
	cmd.Stdout = sout
	cmd.Stderr = serr

	// Stdin doesn't work with a multireader for some reason.
	// cmd.Stdin = sin
	cmd.Stdin = os.Stdin

	return cmd
}
