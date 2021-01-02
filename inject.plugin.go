package nji

import "reflect"
type inj = func(base ViewAddr, c *Context)

type PluginGroupCtrl int8

const (
	PluginGroupCtrlSuccess PluginGroupCtrl = iota
	PluginGroupCtrlSkip
	PluginGroupCtrlFail
)

type Plugin interface {
	Inject(f reflect.StructField) func(base ViewAddr, c *Context)
}

type PluginGroup interface {
	InjectAndControl(f reflect.StructField) func(base ViewAddr, c *Context) PluginGroupCtrl
}


type InnerPluginPathParam struct {
	Value string
}

func (pl *InnerPluginPathParam) Exec(c *Context, name string) {
	pl.Value,_ = c.PathParams.Get(name)
}

func (pl InnerPluginPathParam) Inject(f reflect.StructField) func(base ViewAddr, c *Context) {
	offset := f.Offset
	name := f.Name
	return func(base ViewAddr, c *Context) {
		(*InnerPluginPathParam)(base.Offset(offset)).Exec(c, name)
	}
}


