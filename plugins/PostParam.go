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

func (pl *PostParam) Inject(c *nji.Context, f reflect.StructField) error {
	var ok bool
	pl.Value,ok = c.PostParam(f.Name)
	if ok{
		return nil
	}else{
		return postParamFail
	}
}

func (pl PostParam) Support() nji.Method {
	return nji.MethodP
}

type PostParamOptional struct {
	PostParam
	optional
}

func (pl *PostParamOptional)Inject(c *nji.Context, f reflect.StructField) error {
	if err := pl.PostParam.Inject(c,f); err != nil{
		pl.optional = false
	}else{
		pl.optional = true
	}
	return nil
}