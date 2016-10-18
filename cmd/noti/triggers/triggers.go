package triggers

import (
	"io"

	"github.com/variadico/noti/cmd/noti/run"
)

type Trigger interface {
	Streams() (stdin io.Reader, stdout io.Writer, stderr io.Writer)
	Run(chan error, chan run.Stats)
}
