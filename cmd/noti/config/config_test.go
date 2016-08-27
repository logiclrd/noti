package config

import (
	"reflect"
	"testing"
)

func TestDecompose(t *testing.T) {
	got := decompose("/home/foo/bar/fizz/buzz/")
	want := []string{
		"/home/foo/bar/fizz/buzz/",
		"/home/foo/bar/fizz/",
		"/home/foo/bar/",
		"/home/foo/",
		"/home/",
		"/",
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got: '%v'; want: '%v'", got, want)
	}
}
