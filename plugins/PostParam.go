package plugins

import (
	"nji"
	"reflect"
)

var _ nji.Plugin = &PostParam{}
var _ nji.Plugin = &PostParamOptional{}

var PostParamFail = err{"PostParamFail"}

type PostParam struct {
	Value string
}

type PostParamOptional struct {
	Value string
	optional
}

func (pl PostParam) Inject(f reflect.StructField) func(base nji.ViewAddr, c *nji.Context) {
	offset := f.Offset
	name := f.Name
	return func(base nji.ViewAddr, c *nji.Context) {
		var pl = (*PostParam)(base.Offset(offset))
		var ok bool
		pl.Value,ok = c.PostParam(name)
		if !ok {
			c.Error = PostParamFail
		}
	}
}

func (pl PostParamOptional) Inject(f reflect.StructField) func(base nji.ViewAddr, c *nji.Context) {
	offset := f.Offset
	name := f.Name
	return func(base nji.ViewAddr, c *nji.Context) {
		var pl = (*PostParamOptional)(base.Offset(offset))
		var ok bool
		pl.Value,ok = c.PostParam(name)
		if ok {
			pl.optional.notEmpty = true
		}
	}
}