package config

import (
	"errors"
	"reflect"
)

func MergeFields(a ...interface{}) error {
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

func mergePointers(n1, n2 interface{}) error {
	a := reflect.ValueOf(n1)
	if err := validateType(a); err != nil {
		return err
	}

	b := reflect.ValueOf(n2)
	if err := validateType(b); err != nil {
		return err
	}

	// Grab struct at pointers.
	a = a.Elem()
	b = b.Elem()
	if !a.CanSet() || !b.CanSet() {
		return errors.New("notification is non-addressable")
	}

	ln := a.NumField()
	if ln != b.NumField() {
		return errors.New("fields number mismatch")
	}

	for i := 0; i < ln; i++ {
		switch a.Field(i).Kind() {
		case reflect.String:
			if v := b.Field(i).String(); v != "" {
				a.Field(i).SetString(v)
			}
		case reflect.Int:
			if v := b.Field(i).Int(); v != 0 {
				a.Field(i).SetInt(v)
			}
		case reflect.Float64:
			if v := b.Field(i).Float(); v != 0.0 {
				a.Field(i).SetFloat(v)
			}
		case reflect.Bool:
			if v := b.Field(i).Bool(); v != false {
				a.Field(i).SetBool(v)
			}
		}
	}

	return nil
}

func validateType(v reflect.Value) error {
	if v.Kind() != reflect.Ptr {
		return errors.New("notification must be pointer type")
	}

	if v.IsNil() {
		return errors.New("notification must be non-nil pointer type")
	}

	return nil
}
