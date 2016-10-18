package triggers

import (
	"io"

	"github.com/variadico/noti/cmd/noti/run"
)

type Trigger interface {
	Run(chan error, chan run.Stats)
}

type Streamer interface {
	Streams() (stdin io.Reader, stdout io.Writer, stderr io.Writer)
}
