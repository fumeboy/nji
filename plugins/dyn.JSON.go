package plugins

import (
	jsoniter "github.com/json-iterator/go"
	"io"
	"nji"
	"reflect"
	"unsafe"
)

var _ nji.Plugin = &DynJSON{}

var dynJSONFail = err{"DynJSONFail"}

type DynJSON struct {}

func (pl *DynJSON) Exec(i io.Reader, obj interface{}) error {
	return jsoniter.NewDecoder(i).Decode(&obj)
}

func (pl *DynJSON) Support() nji.Method {
	return nji.MethodP
}

func (pl *DynJSON) Inject(f reflect.StructField) func(base nji.ViewAddress, c *nji.Context) {
	fv := reflect.New(f.Type).Interface()
	offset := f.Offset
	return func(base nji.ViewAddress, c *nji.Context) {
		if ct := c.Request.Header.Get("Content-Type"); ct != "application/json" {
			c.Error = dynJSONFail
			return
		}
		vi := fv
		(*(*face)(unsafe.Pointer(&vi))).word = base.Offset(offset)
		c.Error = jsoniter.NewDecoder(c.Request.Body).Decode(&vi)
	}
}

