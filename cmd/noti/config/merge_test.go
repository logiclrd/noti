package config

import (
	"reflect"
	"testing"

	"github.com/variadico/noti/nsuser"
)

func TestMergeFields(t *testing.T) {
	ptrs := func(n *nsuser.Notification) []interface{} {
		return []interface{}{
			&n.Title,
			&n.Subtitle,
			&n.InformativeText,
			&n.ContentImage,
			&n.SoundName,
		}
	}

	blank := new(nsuser.Notification)
	preset := &nsuser.Notification{
		Title:           "testing",
		InformativeText: "hello",
	}

	if err := MergeFields(ptrs(blank), ptrs(preset)); err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if !reflect.DeepEqual(blank, preset) {
		t.Error("Failed equality")
		t.Errorf("got: %v; want: %v", blank, preset)
	}
}

func TestMergePointers(t *testing.T) {
	strp := func(v string) *string { return &v }
	intp := func(v int) *int { return &v }
	f64p := func(v float64) *float64 { return &v }
	boolp := func(v bool) *bool { return &v }

	a := []interface{}{
		strp("hello"),
	}
	b := []interface{}{
		strp("hola"),
		strp("mundo"),
	}

	if err := mergePointers(a, b); err == nil {
		t.Error("Length mismatch")
		t.Error("Unexpected success")
	}

	b = []interface{}{
		intp(42),
	}
	if err := mergePointers(a, b); err == nil {
		t.Error("Type mismatch: *string and *int")
		t.Error("Unexpected success")
	}

	b = []interface{}{
		"hello",
	}
	if err := mergePointers(a, b); err == nil {
		t.Error("Type mismatch: *string and string")
		t.Error("Unexpected success")
	}

	a = []interface{}{
		strp("hello"),
		intp(42),
		f64p(3.14),
		boolp(true),
	}
	b = []interface{}{
		strp(""),
		intp(0),
		f64p(0.0),
		boolp(false),
	}
	want := []interface{}{
		strp("hello"),
		intp(42),
		f64p(3.14),
		boolp(true),
	}

	t.Run("keep first struct values", func(t *testing.T) {
		if err := mergePointers(a, b); err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		for i, v := range want {
			if !reflect.DeepEqual(a[i], v) {
				t.Error("Failed equality")

				t.Errorf(
					"got: %v; want: %v",
					reflect.ValueOf(a[i]).Elem(),
					reflect.ValueOf(v).Elem(),
				)
			}
		}
	})

	a = []interface{}{
		strp(""),
		intp(0),
		f64p(0.0),
		boolp(false),
	}
	b = []interface{}{
		strp("hello"),
		intp(42),
		f64p(3.14),
		boolp(true),
	}

	t.Run("keep second struct values", func(t *testing.T) {
		if err := mergePointers(a, b); err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		for i, v := range b {
			if !reflect.DeepEqual(a[i], v) {
				t.Error("Failed equality")

				t.Errorf(
					"got: %v; want: %v",
					reflect.ValueOf(a[i]).Elem(),
					reflect.ValueOf(v).Elem(),
				)
			}
		}
	})

}
