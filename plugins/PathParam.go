package plugins

import (
	"github.com/fumeboy/nji"
	"github.com/fumeboy/nji/schema"
	"reflect"
)

type PathParam[T any] struct {
	checker schema.StringChecker

	Value string
}

func (pl *PathParam[T]) Build(_ reflect.StructField) error {
	pl.checker, _ = schema.BuildStringChecker[T]()
	return nil
}

func (pl *PathParam[T]) Run(c *nji.Context, f reflect.StructField) error {
	pl.Value = c.Param(f.Name)
	return pl.checker.Check(pl.Value)
}
