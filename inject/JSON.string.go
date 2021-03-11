package inject

import (
	jsoniter "github.com/json-iterator/go"
	"nji"
	"reflect"
)

var _ nji.Plugin = &JSONFieldStr{}

var JSONFieldStrFail = err{"JSONFieldStrFail"}

type JSONFieldStr struct {
	Value string
}

func (pl *JSONFieldStr) Inject(c *nji.Context, f reflect.StructField) error{
	if !c.IsJSON() {
		return JSONFieldStrFail
	}
	pl.Value = jsoniter.Get(c.JSON, f.Name).ToString()
	return nil
}
