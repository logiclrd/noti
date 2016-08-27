package config

import (
	"bytes"
	"text/template"

	"github.com/variadico/noti/cmd/noti/run"
)

// EvalFields evaluates string fields as a text template.
func EvalFields(fs []interface{}, st run.Stats) error {
	var err error

	for _, field := range fs {
		strField, is := field.(*string)
		if !is {
			continue
		}

		*strField, err = eval(*strField, st)
		if err != nil {
			return err
		}
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
