package plugins

import (
	"nji"
	"reflect"
)

var _ nji.Plugin = &PostParam{}
var _ nji.Plugin = &PostParamOptional{}

var postParamFail = err{"PostParamFail"}

type PostParam struct {
	Value string
}

func (pl PostParam) Support() nji.Method {
	return nji.MethodP
}

type PostParamOptional struct {
	PostParam
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
			c.Error = postParamFail
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