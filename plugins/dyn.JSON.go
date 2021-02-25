package plugins

import (
	jsoniter "github.com/json-iterator/go"
	"nji"
	"reflect"
)

var _ nji.PluginDyn = &DynJSON{}

var dynJSONFail = err{"DynJSONFail"}

type DynJSON struct {}

func (pl DynJSON) Inject(c *nji.Context, f reflect.StructField, iface interface{}) error {
	if ct := c.Request.Header.Get("Content-Type"); ct != "application/json" {
		return dynJSONFail
	}
	return  jsoniter.NewDecoder(c.Request.Body).Decode(&iface)
}

func (pl *DynJSON) Support() nji.Method {
	return nji.MethodP
}
