package plugins

import (
	"reflect"

	"github.com/fumeboy/nji"
	"github.com/fumeboy/nji/schema"
)

type QueryParam[T any] struct {
	checker     schema.StringChecker
	mustNotNull bool

	IsNull bool
	Value  string
}

func (pl *QueryParam[T]) Build(_ reflect.StructField) error {
	pl.checker, pl.mustNotNull = schema.BuildStringChecker[T]()
	return nil
}

func (pl *QueryParam[T]) Run(c *nji.Context, f reflect.StructField) error {
	if q, ok := c.GetQuery(f.Name); !ok {
		if pl.mustNotNull {
			return err{"queryParamFail"}
		} else {
			pl.IsNull = true
			return nil
		}
	} else {
		pl.Value = q
		return pl.checker.Check(pl.Value)
	}
}
