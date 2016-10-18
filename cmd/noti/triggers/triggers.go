package triggers

import (
	"io"
	"os/exec"
	"syscall"
)

type trigger interface {
	streams() (stdin io.Reader, stdout io.Writer, stderr io.Writer)
	run(chan error, chan Stats)
}

func exitStatus(err error) int {
	eerr, is := err.(*exec.ExitError)
	if !is {
		return noExitStatus
	}

	if status, is := eerr.Sys().(syscall.WaitStatus); is {
		return status.ExitStatus()
	}

	return noExitStatus
}
