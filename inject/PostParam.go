package inject

import (
	"nji"
	"reflect"
)

var _ nji.Plugin = new(PostParam)

var postParamFail = err{"PostParamFail"}

type PostParam struct {
	Val string
}

func (pl *PostParam) Inject(ctx *nji.Context, f reflect.StructField) error {
	if err := ctx.ParseForm(); err != nil {
		return postParamFail
	}
	vs := ctx.Request.PostForm[f.Name]
	if len(vs) == 0 {
		return postParamFail
	}
	pl.Val = ctx.Request.PostForm[f.Name][0]
	return nil
}