package schema

import (
	"reflect"
)

type stringchecker struct {
	checker  stringCheckerIface
	metadata string
}

type StringChecker struct {
	checkers []stringchecker
}

func (sc *StringChecker) Check(v string) error {
	for _, c := range sc.checkers {
		if err := c.checker.check(v, c.metadata); err != nil {
			return err
		}
	}
	return nil
}

type stringCheckerIface interface {
	check(v string, metadata string) error
}

func BuildStringChecker[T any]() (checker StringChecker, mustNotNull bool) {
	t := reflect.TypeOf(new(T)).Elem()
	if t.Kind() == reflect.Interface {
		if t.NumMethod() == 0 {
			return
		}
		panic("invalid struct field")
	}
	if t.Kind() == reflect.Struct {
		if t.PkgPath() != "" {
			if t.Name() == reflect.TypeOf(Must{}).Name() {
				mustNotNull = true
				return
			} else if fv, ok := reflect.New(t).Interface().(stringCheckerIface); ok {
				checker.checkers = append(checker.checkers, stringchecker{fv, ""})
				return
			}
		} else {
			length := t.NumField()
			for i := 0; i < length; i++ {
				f := t.Field(i)
				if f.Type.Kind() == reflect.Interface {
					panic("invalid struct field")
				}
				if f.Type.Kind() == reflect.Ptr {
					panic("invalid struct field")
				}
				if f.Type.Name() == reflect.TypeOf(Must{}).Name() {
					mustNotNull = true
					continue
				}
				if fv, ok := reflect.New(f.Type).Interface().(stringCheckerIface); ok {
					checker.checkers = append(checker.checkers, stringchecker{fv, string(f.Tag)})
				} else {
					panic("StringChecker not support `" + f.Type.Name() + "`")
				}
			}
			return
		}
	}
	panic("StringChecker not support `" + t.Name() + "`")
}
