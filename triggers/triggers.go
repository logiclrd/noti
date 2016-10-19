package triggers

import (
	"io"

	"github.com/variadico/noti/cmd/noti/runstat"
)

type Trigger interface {
	Run(chan error, chan runstat.Result)
}

type Streamer interface {
	Streams() (stdin io.Reader, stdout io.Writer, stderr io.Writer)
}
