package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	yesAns = iota
	noAns
	skipAns
)

func printTitle(s string) {
	fmt.Println("\n***")
	fmt.Println("*", s)
	fmt.Println("***")
}

func printWarning(s string) {
	fmt.Println("--- WARNING ---")
	fmt.Println(s)
}

func promptString(s string) (string, error) {
	fmt.Printf("%s: ", s)
	r := bufio.NewReader(os.Stdin)
	ans, err := r.ReadString('\n')
	return strings.TrimSpace(ans), err
}

func promptYesNo(s string) (bool, error) {
	ans, err := promptString(fmt.Sprintf("%s (y/n)", s))
	if err != nil {
		return false, err
	}

	ans = strings.ToLower(ans)
	if ans == "y" {
		return true, nil
	} else if ans == "n" {
		return false, nil
	}

	return false, fmt.Errorf("unknown answer: %s", ans)
}

func promptYesAbort(s string) (bool, error) {
	ok, err := promptYesNo(s)
	if err != nil {
		return false, err
	}

	if !ok {
		return false, fmt.Errorf("abort %q", s)
	}

	return true, nil
}

func promptYesNoSkip(s string) (int, error) {
	ans, err := promptString(fmt.Sprintf("%s (y/n/s)", s))
	if err != nil {
		return -1, err
	}

	ans = strings.ToLower(ans)
	if ans == "y" {
		return yesAns, nil
	} else if ans == "n" {
		return noAns, nil
	} else if ans == "s" {
		return skipAns, nil
	}

	return -1, fmt.Errorf("unknown answer: %s", ans)
}

func promptYesAbortSkip(s string) (int, error) {
	ans, err := promptYesNoSkip(s)
	if err != nil {
		return -1, err
	}

	if ans == noAns {
		return -1, fmt.Errorf("abort %q", s)
	}

	return ans, nil
}
