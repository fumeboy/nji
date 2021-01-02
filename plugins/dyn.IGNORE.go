package plugins

import (
	"nji"
	"reflect"
)

var _ nji.PluginGroup = &DynIgnore{}

type DynIgnore struct {
}

func (pg *DynIgnore) Support() nji.Method {
	return nji.MethodAny
}

func (pg *DynIgnore) InjectAndControl(f reflect.StructField) func(base nji.ViewAddr, c *nji.Context) nji.PluginGroupCtrl{
	return func(base nji.ViewAddr, c *nji.Context) nji.PluginGroupCtrl{
		return nji.PluginGroupCtrlSkip
	}
}
