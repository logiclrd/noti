package config

import (
	"errors"
	"fmt"
	"log"
	"runtime"
)

func MergeFields(a ...[]interface{}) error {
	ln := len(a)

	if ln < 2 {
		return errors.New("not enough arguments")
	}

	for i := 1; i < ln; i++ {
		if err := mergePointers(a[0], a[i]); err != nil {
			return err
		}
	}

	return nil
}

func mergePointers(a, b []interface{}) (err error) {
	if len(a) != len(b) {
		return errors.New("merge slice length mismatch")
	}
	ln := len(a)

	defer func() {
		switch v := recover().(type) {
		case *runtime.TypeAssertionError:
			err = errors.New("merge type mismatch")
		case error:
			err = v
		case nil:
			// Carry on. Nothing to see here.
		default:
			log.Println(`IDK. ¯\_(ツ)_/¯`)
			panic(v)
		}
	}()

	for i := 0; i < ln; i++ {
		switch v := a[i].(type) {
		case *string:
			if *b[i].(*string) != "" {
				*a[i].(*string) = *b[i].(*string)
			}
		case *int:
			if *b[i].(*int) != 0 {
				*a[i].(*int) = *b[i].(*int)
			}
		case *float64:
			if *b[i].(*float64) != 0.0 {
				*a[i].(*float64) = *b[i].(*float64)
			}
		case *bool:
			if *b[i].(*bool) != false {
				*a[i].(*bool) = *b[i].(*bool)
			}
		default:
			return fmt.Errorf("unsupported merge type: %T", v)
		}
	}

	return nil
}
