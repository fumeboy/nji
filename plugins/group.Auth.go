package plugins

import (
	"nji"
	"reflect"
	"github.com/dgrijalva/jwt-go"
)

var _ nji.PluginGroup = &GroupAuth{}

type GroupAuth struct {
}

func (g GroupAuth) InjectAndControl(f reflect.StructField) func(base nji.ViewAddr, c *nji.Context) nji.PluginGroupCtrl {
	panic("implement me")
}

func (g GroupAuth) Support() nji.Method {
	return nji.MethodAny
}
