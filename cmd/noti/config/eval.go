package config

import (
	"bytes"
	"errors"
	"reflect"
	"text/template"

	"github.com/variadico/noti/cmd/noti/run"
)

// EvalFields evaluates string fields as a text template. n should be
// a non-nil pointer type. It will be modified.
func EvalStringFields(n interface{}, st run.Stats) error {
	// Grab underlying value of n.
	v := reflect.ValueOf(n)

	if v.Kind() != reflect.Ptr {
		return errors.New("notification must be pointer type")
	}
	if v.IsNil() {
		return errors.New("notification must be non-nil pointer type")
	}

	// Grab the element at pointer address.
	v = v.Elem()

	var s string
	var err error
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Kind() != reflect.String {
			continue
		}

		s, err = eval(v.Field(i).String(), st)
		if err != nil {
			return err
		}
		v.Field(i).SetString(s)
	}

	return nil
}

func eval(s string, st run.Stats) (string, error) {
	tmpl, err := template.New("").Parse(s)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, st); err != nil {
		return "", err
	}

	return buf.String(), nil
}
