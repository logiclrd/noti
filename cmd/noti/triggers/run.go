package triggers

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/variadico/noti/cmd/noti/run"
	"github.com/variadico/noti/cmd/noti/triggers/exit"
	"github.com/variadico/noti/cmd/noti/triggers/match"
	"github.com/variadico/noti/cmd/noti/triggers/timeout"
)

const (
	delim = "="
)

func Run(trigFlags []string, args []string, notify func(run.Stats) error) error {
	if len(trigFlags) == 0 {
		trigFlags = append(trigFlags, exit.FlagKey)
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Typically, the context will get cancelled by the first trigger to
	// finish. It'll also get called if we hit an error and return early. Or
	// return normally. It's okay if this gets called multiple times.
	defer cancel()

	sts := run.NewStats(args)

	trigs := make([]Trigger, 0, len(trigFlags))
	for _, t := range trigFlags {
		name, val := keyValue(t)

		switch name {
		case exit.FlagKey:
			trigs = append(trigs, exit.NewTrigger(ctx, sts))
		case match.FlagKey:
			trigs = append(trigs, match.NewTrigger(ctx, sts, val))
		case timeout.FlagKey:
			t, err := timeout.NewTrigger(ctx, sts, val)
			if err != nil {
				return err
			}

			trigs = append(trigs, t)
		default:
			return fmt.Errorf("unknown trigger: %s", name)
		}
	}

	cmd := newCmd(ctx, sts, trigs)
	errc := make(chan error)
	go func() { errc <- cmd.Run() }()

	var wg sync.WaitGroup
	statsc := make(chan run.Stats)

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
	inmap := make(map[io.Reader]struct{})
	outmap := make(map[io.Writer]struct{})
	errmap := make(map[io.Writer]struct{})

	// Make streams unique.
	for _, t := range ts {
		sin, sout, serr := t.Streams()
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

	fmt.Println("numstdins:", len(stdins))
	fmt.Println("numstdouts:", len(stdouts))
	fmt.Println("numstderrs:", len(stderrs))

	return io.MultiReader(stdins...), io.MultiWriter(stdouts...), io.MultiWriter(stderrs...)
}

func newCmd(ctx context.Context, sts run.Stats, ts []Trigger) *exec.Cmd {
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
