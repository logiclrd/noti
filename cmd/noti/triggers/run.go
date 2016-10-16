package triggers

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

const (
	delim = "="
)

type NotifyFn func(Stats) error

func Run(trigs []string, args []string, notify NotifyFn) error {
	fmt.Println(">>>>", trigs)

	if len(trigs) == 0 {
		trigs = append(trigs, "exit")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	out := make(chan Stats)
	var wg sync.WaitGroup

	for _, t := range trigs {
		name, val := keyValue(t)
		wg.Add(1)

		switch name {
		case "exit":
			go func() { onExit(ctx, args, out); wg.Done() }()
		case "contains":
			go func() { onContains(ctx, val, args, out); wg.Done() }()
		case "timeout":
			go func() { onTimeout(ctx, val, args, out); wg.Done() }()
		default:
			return fmt.Errorf("unknown trigger: %s", name)
		}
	}

	go func() { wg.Wait(); close(out) }()

	for stat := range out {
		if err := notify(stat); err != nil {
			return err
		}
	}

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
