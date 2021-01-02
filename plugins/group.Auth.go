package plugins

import (
	"nji"
	"reflect"
)

var _ nji.PluginGroup = &GroupAuth{}

type GroupAuth struct {
	Auth
}

func (g GroupAuth) InjectAndControl(f reflect.StructField) func(base nji.ViewAddr, c *nji.Context) nji.PluginGroupCtrl {
	panic("implement me")
}

func (g GroupAuth) Support() nji.Method {
	return nji.MethodAny
}