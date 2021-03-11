package inject

import (
	"errors"
	"nji"
	"reflect"
)

var _ nji.Plugin = &PathParam{}

type PathParam struct {
	Value string
}

func (pl *PathParam) Inject(c *nji.Context, f reflect.StructField) error{
	var ok bool
	pl.Value,ok = c.PathParams.Get(f.Name)
	if ok {
		return nil
	}
	return errors.New("bad path: " + f.Name)
}