package config

import (
	"reflect"
	"testing"

	"github.com/variadico/noti/nsuser"
)

func TestMergeFields(t *testing.T) {
	blank := new(nsuser.Notification)
	preset := &nsuser.Notification{
		Title:           "testing",
		InformativeText: "hello",
	}

	if err := MergeFields(blank, preset); err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if !reflect.DeepEqual(blank, preset) {
		t.Error("Failed equality")
		t.Errorf("got: %v; want: %v", blank, preset)
	}
}
