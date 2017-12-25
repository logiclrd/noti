package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func runCommand(bin string, args ...string) error {
	out, err := exec.Command(bin, args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run %s: %s: %s", bin, err, string(out))
	}

	return nil
}

func commandOutput(bin string, args ...string) (string, error) {
	cmd := exec.Command(bin, args...)
	var sout bytes.Buffer
	var serr bytes.Buffer
	cmd.Stdout = &sout
	cmd.Stderr = &serr

	if err := cmd.Run(); err != nil {
		combined := sout.String() + serr.String()
		return "", fmt.Errorf("failed to run %s: %s: %s", bin, err, combined)
	}

	return sout.String(), nil
}

func gitCheckout(branch string) error {
	return runCommand("git", "checkout", branch)
}

func gitMergeFF(branch string) error {
	return runCommand("git", "merge", branch, "--ff-only")
}

func gitPushOrigin(branch string) error {
	return runCommand("git", "push", "origin", branch)
}

func gitTag(tag string) error {
	return runCommand("git", "tag", tag)
}

func currentBranch() (string, error) {
	out, err := commandOutput("git", "rev-parse", "--abbrev-ref", "HEAD")
	return strings.TrimSpace(out), err
}

func currentTag() (string, error) {
	out, err := commandOutput("git", "describe", "--abbrev=0", "--tags")
	return strings.TrimSpace(out), err
}
