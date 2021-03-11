package inject

import (
	"nji"
	"reflect"
)

var _ nji.Plugin = new(FormParam)

var formParamFail = err{"FormParamFail"}

type FormParam struct {
	Val string
}

func (pl *FormParam) Inject(ctx *nji.Context, f reflect.StructField) error {
	if err := ctx.ParseForm(); err != nil {
		return formParamFail
	}
	vs := ctx.Request.Form[f.Name]
	if len(vs) == 0 {
		return formParamFail
	}
	pl.Val = ctx.Request.Form[f.Name][0]
	return nil
}