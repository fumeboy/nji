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

func (pl *PostParam) Exec(c *nji.Context, name string) error {
	var ok bool
	pl.Value,ok = c.PostParam(name)
	if ok{
		return nil
	}else{
		return postParamFail
	}
}

type PostParamOptional struct {
	PostParam
	optional
}

func (pl PostParam) Inject(f reflect.StructField) func(base nji.ViewAddr, c *nji.Context) {
	offset := f.Offset
	name := f.Name
	return func(base nji.ViewAddr, c *nji.Context) {
		c.Error = (*PostParam)(base.Offset(offset)).Exec(c,name)
	}
}

func (pl PostParamOptional) Inject(f reflect.StructField) func(base nji.ViewAddr, c *nji.Context) {
	offset := f.Offset
	name := f.Name
	return func(base nji.ViewAddr, c *nji.Context) {
		if (*PostParam)(base.Offset(offset)).Exec(c,name) != nil{
			pl.optional.notEmpty = true
		}
	}
}