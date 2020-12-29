package plugins

import (
	"nji"
	"reflect"
)

var _ nji.Plugin = &PathParam{}
var _ nji.Plugin = &PathParamOptional{}

type PathParam struct {
	Value string
}

type PathParamOptional struct {
	PathParam
	notEmpty
}

func (pl *PathParam) Inject(f reflect.StructField) func(c *nji.Context) {
	offset := f.Offset
	name := f.Name
	return func(c *nji.Context) {
		var ok bool
		var pl = (*PathParam)(c.OffsetV(offset))
		if pl.Value, ok = c.PathParams.Get(name); !ok{
			c.Abort()
		}
	}
}

func (pl *PathParamOptional) Inject(f reflect.StructField) func(c *nji.Context) {
	offset := f.Offset
	name := f.Name
	return func(c *nji.Context) {
		var ok bool
		var pl = (*PathParamOptional)(c.OffsetV(offset))
		if pl.Value, ok = c.PathParams.Get(name); ok{
			pl.notEmpty.value = true
		}
	}
}