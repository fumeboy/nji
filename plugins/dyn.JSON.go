package plugins

import (
	jsoniter "github.com/json-iterator/go"
	"nji"
	"reflect"
	"unsafe"
)

var _ nji.Plugin = &DynJSON{}

type DynJSON struct {}

func (g *DynJSON) Inject(f reflect.StructField) func(base nji.ViewAddr, c *nji.Context) {
	fv := reflect.New(f.Type).Interface()
	offset := f.Offset
	return func(base nji.ViewAddr, c *nji.Context) {
		vi := fv
		(*(*face)(unsafe.Pointer(&vi))).word = base.Offset(offset)
		c.Error = jsoniter.NewDecoder(c.Request.Body).Decode(&vi)
	}
}

