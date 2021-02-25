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

func (pl *QueryParam) Inject(c *nji.Context, f reflect.StructField) error {
	var ok bool
	pl.Value,ok = c.QueryParam(f.Name)
	if ok{
		return nil
	}else{
		return queryParamFail
	}
}

func (pl *QueryParamOptional) Inject(c *nji.Context, f reflect.StructField) error {
	if err := pl.QueryParam.Inject(c,f); err != nil{
		pl.optional = false
	}else{
		pl.optional = true
	}
	return nil
}