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
		panic("invalid type for schema")
	}
	var typecheck = func(t reflect.Type, gotag string) {
		if t.Kind() == reflect.Interface {
			panic("invalid type for schema")
		}
		if t.Kind() == reflect.Ptr {
			panic("invalid type for schema")
		}
		if t.String() == reflect.TypeOf(NotNull{}).String() {
			mustNotNull = true
			return
		}
		if fv, ok := reflect.New(t).Interface().(stringCheckerIface); ok {
			checker.checkers = append(checker.checkers, stringchecker{fv, gotag})
		} else {
			panic("StringChecker not support `" + t.Kind().String() + " " + t.String() + "`")
		}
	}
	if t.Kind() == reflect.Struct {
		if t.PkgPath() != "" {
			typecheck(t, "")
			return
		} else {
			for i := 0; i < t.NumField(); i++ {
				f := t.Field(i)
				typecheck(f.Type, string(f.Tag))
			}
			return
		}
	}
	if t.Kind() == reflect.Func && t.PkgPath() == "" {
		for i := 0; i < t.NumIn(); i++ {
			arg := t.In(i)
			typecheck(arg, "")
		}
		return
	}
	panic("StringChecker not support `" + t.Kind().String() + " " + t.String() + "`")
}
