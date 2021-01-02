package plugins

import (
	"nji"
	"reflect"
)

var _ nji.Plugin = &QueryParam{}
var _ nji.Plugin = &QueryParamOptional{}

var queryParamFail = err{"queryParamFail"}

type QueryParam struct {
	Value string
}

func (pl QueryParam) Support() nji.Method {
	return nji.MethodGet | nji.MethodHead
}

type QueryParamOptional struct {
	QueryParam
	optional
}

func (pl QueryParam) Inject(f reflect.StructField) func(base nji.ViewAddr, c *nji.Context) {
	offset := f.Offset
	name := f.Name
	return func(base nji.ViewAddr, c *nji.Context) {
		var pl = (*QueryParam)(base.Offset(offset))
		var ok bool
		pl.Value,ok = c.QueryParam(name)
		if !ok {
			c.Error = queryParamFail
		}
	}
}

func (pl QueryParamOptional) Inject(f reflect.StructField) func(base nji.ViewAddr, c *nji.Context) {
	offset := f.Offset
	name := f.Name
	return func(base nji.ViewAddr, c *nji.Context) {
		var pl = (*QueryParamOptional)(base.Offset(offset))
		var ok bool
		pl.Value,ok = c.QueryParam(name)
		if ok {
			pl.optional.notEmpty = true
		}
	}
}