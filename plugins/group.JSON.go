package plugins

import (
	"nji"
	"reflect"
)

var _ nji.PluginGroup = &GroupJSON{}

type GroupJSON struct {}

func (pg *GroupJSON) Proxy(f reflect.StructField) (fn func(base nji.ViewAddr, c *nji.Context), ok bool) {
	//name := f.Name
	return func(base nji.ViewAddr, c *nji.Context) {

	}, true
}

func (pg *GroupJSON) Control() func(base nji.ViewAddr, c *nji.Context){
	return nil
}