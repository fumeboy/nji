package plugins

import (
	"github.com/fumeboy/nji"
	"github.com/fumeboy/nji/schema"
	"reflect"
)

type PostParam[T any] struct {
	checker     schema.StringChecker
	mustNotNull bool

	IsNull bool
	Value  string
}

func (pl *PostParam[T]) Build(_ reflect.StructField) error {
	pl.checker, pl.mustNotNull = schema.BuildStringChecker[T]()
	return nil
}

func (pl *PostParam[T]) Run(c *nji.Context, f reflect.StructField) error {
	if q, ok := c.GetPostForm(f.Name); !ok {
		if pl.mustNotNull {
			return err{"PostParamFail"}
		} else {
			pl.IsNull = true
			return nil
		}
	} else {
		pl.Value = q
		return pl.checker.Check(pl.Value)
	}
}
