package config

import (
	"testing"

	"github.com/variadico/noti/runstat"
)

func TestEvalStringFields(t *testing.T) {
	st := runstat.Result{
		Cmd: "testing",
	}

	s := struct {
		Title string
		Num   int
	}{
		Title: "{{.Cmd}}",
		Num:   42,
	}

	if err := EvalStringFields(&s, st); err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if s.Title != st.Cmd {
		t.Error("Unexpected eval")
		t.Errorf("got: %v; want: %v", s.Title, st.Cmd)
	}
	if s.Num != 42 {
		t.Error("Unexpected eval")
		t.Errorf("got: %v; want: %v", s.Num, 42)
	}
}
