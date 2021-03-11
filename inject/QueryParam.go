package inject

import (
	"nji"
	"reflect"
)

var _ nji.Plugin = new(QueryParam)

var queryParamFail = err{"queryParamFail"}

type QueryParam struct {
	Val string
}

func (pl *QueryParam) Inject(c *nji.Context, f reflect.StructField) error {
	if len(c.Request.URL.Query()[f.Name]) == 0 {
		return queryParamFail
	}
	pl.Val = c.Request.URL.Query()[f.Name][0]
	return nil
}