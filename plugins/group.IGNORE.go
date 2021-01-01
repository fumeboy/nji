package plugins

import (
	"nji"
	"reflect"
)

var _ nji.PluginGroup = &GroupIgnore{}

type GroupIgnore struct {
}

func (pg *GroupIgnore) Proxy(f reflect.StructField) (fn func(base nji.ViewAddr, c *nji.Context), ok bool) {
	return nil, true
}

func (pg *GroupIgnore) Control() func(base nji.ViewAddr, c *nji.Context){
	return func(base nji.ViewAddr, c *nji.Context) {
		c.Error = nji.SkipE{}
	}
}
