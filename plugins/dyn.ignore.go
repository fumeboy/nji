package plugins

import (
	"nji"
	"reflect"
)

var _ nji.Plugin = &DynIgnore{}

type DynIgnore struct {
}

func (pg *DynIgnore) Support() nji.Method {
	return nji.MethodAny
}

func (pg *DynIgnore) Inject(f reflect.StructField) func(base nji.ViewAddress, c *nji.Context){
	return nil
}
