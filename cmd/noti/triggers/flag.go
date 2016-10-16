package triggers

import "fmt"

type Flag []string

func (t *Flag) String() string {
	return fmt.Sprint(*t)
}

func (t *Flag) Set(value string) error {
	*t = append(*t, value)
	return nil
}
