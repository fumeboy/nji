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

func (pl *QueryParam) Exec(c *nji.Context, name string) error {
	var ok bool
	pl.Value,ok = c.QueryParam(name)
	if ok{
		return nil
	}else{
		return queryParamFail
	}
}

func (pl QueryParam) Inject(f reflect.StructField) func(base nji.ViewAddr, c *nji.Context) {
	offset := f.Offset
	name := f.Name
	return func(base nji.ViewAddr, c *nji.Context) {
		c.Error = (*QueryParam)(base.Offset(offset)).Exec(c, name)
	}
}

func (pl QueryParamOptional) Inject(f reflect.StructField) func(base nji.ViewAddr, c *nji.Context) {
	offset := f.Offset
	name := f.Name
	return func(base nji.ViewAddr, c *nji.Context) {
		if (*QueryParam)(base.Offset(offset)).Exec(c, name) != nil {
			pl.optional.notEmpty = true
		}
	}
}