package plugins

import (
	jsoniter "github.com/json-iterator/go"
	"nji"
	"reflect"
	"unsafe"
)

var _ nji.Plugin = &DynJSON{}

var dynJSONFail = err{"DynJSONFail"}

type DynJSON struct {}

func (g *DynJSON) Support() nji.Method {
	return nji.MethodP
}

func (g *DynJSON) Inject(f reflect.StructField) func(base nji.ViewAddr, c *nji.Context) {
	fv := reflect.New(f.Type).Interface()
	offset := f.Offset
	return func(base nji.ViewAddr, c *nji.Context) {
		if ct := c.Request.Header.Get("Content-Type"); ct != "application/json" {
			c.Error = dynJSONFail
			return
		}
		vi := fv
		(*(*face)(unsafe.Pointer(&vi))).word = base.Offset(offset)
		c.Error = jsoniter.NewDecoder(c.Request.Body).Decode(&vi)
	}
}

