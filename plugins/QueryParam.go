package plugins

import (
	"nji"
	"reflect"
)

var _ nji.Plugin = &QueryParam{}
var _ nji.Plugin = &QueryParamOptional{}

var QueryParamFail = err{"QueryParamFail"}

type QueryParam struct {
	Value string
}

type QueryParamOptional struct {
	Value string
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
			c.Error = QueryParamFail
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