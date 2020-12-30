package plugins

import (
	"nji"
	"reflect"
)

var _ nji.Plugin = &PathParam{}

type PathParam struct {
	Value string
}

func (pl PathParam) Inject(f reflect.StructField) func(base nji.ViewAddr, c *nji.Context) {
	offset := f.Offset
	name := f.Name
	return func(base nji.ViewAddr, c *nji.Context) {
		var pl = (*PathParam)(base.Offset(offset))
		pl.Value,_ = c.PathParams.Get(name)
	}
}