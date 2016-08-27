package config

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/variadico/yaml"
)

const (
	Filename = ".noti.yaml"
)

func File() (Options, error) {
	ds, err := dirs()
	if err != nil {
		return Options{}, err
	}

	var data []byte
	for _, d := range ds {
		data, err = ioutil.ReadFile(filepath.Join(d, Filename))
		if err == nil {
			break
		}
	}
	if err != nil {
		return Options{}, fmt.Errorf("config not found in: %s", ds)
	}

	var opts Options
	if err := yaml.Unmarshal(data, &opts); err != nil {
		return Options{}, err
	}

	return opts, nil
}

func dirs() ([]string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	return decompose(dir), nil
}

func decompose(p string) []string {
	var out []string

	if !strings.HasSuffix(p, string(os.PathSeparator)) {
		p += string(os.PathSeparator)
	}

	for i := len(p) - 1; i >= 0; i-- {
		if p[i] == '/' {
			out = append(out, p[:i+1])
		}
	}

	return out
}

func WasSet(fs *flag.FlagSet, name string) bool {
	var wasSet bool

	fs.Visit(func(f *flag.Flag) {
		if f.Name == name {
			wasSet = true
		}
	})

	return wasSet
}
