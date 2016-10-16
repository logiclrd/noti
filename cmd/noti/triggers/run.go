package triggers

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
)

const (
	delim = "="
)

type NotifyFn func(Stats) error

func Run(trigFlags []string, args []string, notify NotifyFn) error {
	fmt.Println(">>>>", trigFlags)

	if len(trigFlags) == 0 {
		trigFlags = append(trigFlags, "exit")
		fmt.Println(">>>> ADDED EXIT TRIGGER")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sts := statsFromArgs(args)

	trigs := make([]trigger, 0, len(trigFlags))
	for _, t := range trigFlags {
		name, val := keyValue(t)

		switch name {
		case "exit":
			trigs = append(trigs, newOnExitTrigger(ctx, sts))
		case "contains":
			trigs = append(trigs, newOnContainsTrigger(ctx, sts, val))
		case "timeout":
			t, err := newOnTimeoutTrigger(ctx, sts, val)
			if err != nil {
				return err
			}

			trigs = append(trigs, t)
		default:
			return fmt.Errorf("unknown trigger: %s", name)
		}
	}

	var wg sync.WaitGroup
	statsc := make(chan Stats)
	cmd := newCmd(ctx, sts, trigs)
	errc := make(chan error)

	go func() { errc <- cmd.Run() }()
	for _, t := range trigs {
		fmt.Printf("trigger: %T\n", t)

		wg.Add(1)
		go func(t trigger) {
			t.run(errc, statsc)
			wg.Done()
			cancel()
		}(t)
	}

	go func() { wg.Wait(); close(statsc) }()

	fmt.Println("waiting for stats")
	for stat := range statsc {
		if err := notify(stat); err != nil {
			return err
		}
	}
	fmt.Println("done waiting")

	return nil
}

func keyValue(trigger string) (string, string) {
	i := strings.Index(trigger, delim)
	if i == -1 {
		// Trigger is something like, "exit"
		return trigger, ""
	}

	// Trigger is something like, "contains=hello world"
	t := trigger[:i]
	return t, trigger[i+1:]
}

func uniqStreams(ts []trigger) (stdin io.Reader, stdout io.Writer, stderr io.Writer) {
	inmap := make(map[io.Reader]struct{})
	outmap := make(map[io.Writer]struct{})
	errmap := make(map[io.Writer]struct{})

	// Make streams unique.
	for _, t := range ts {
		sin, sout, serr := t.streams()
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

func newCmd(ctx context.Context, sts Stats, ts []trigger) *exec.Cmd {
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
