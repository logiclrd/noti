package config

import (
	"bytes"
	"reflect"
	"text/template"

	"github.com/variadico/noti/cmd/noti/triggers"
)

// EvalFields evaluates string fields as a text template. n should be
// a non-nil pointer type. It will be modified.
func EvalStringFields(n interface{}, st triggers.Stats) error {
	// Grab underlying value of n.
	v := reflect.ValueOf(n)
	if err := validateType(v); err != nil {
		return err
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

func eval(s string, st triggers.Stats) (string, error) {
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
