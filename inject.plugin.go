package nji

import "reflect"
type inj = func(base ViewAddr, c *Context)

type Method int16

const (
	_             Method = iota
	MethodGet            = 1
	MethodHead           = 1 << 1
	MethodPost           = 1 << 2
	MethodPut            = 1 << 3
	MethodPatch          = 1 << 4
	MethodDelete         = 1 << 5
	MethodConnect        = 1 << 6
	MethodOptions        = 1 << 7
	MethodAny            = MethodGet | MethodHead | MethodPost | MethodPut | MethodPatch | MethodDelete | MethodConnect | MethodOptions
	MethodP              = MethodPost | MethodPut | MethodPatch
)

type PluginGroupCtrl int8

const (
	PluginGroupCtrlSuccess PluginGroupCtrl = iota
	PluginGroupCtrlSkip
	PluginGroupCtrlFail
)

type Plugin interface {
	Inject(f reflect.StructField) func(base ViewAddr, c *Context)
	Support() Method
}

type PluginGroup interface {
	InjectAndControl(f reflect.StructField) func(base ViewAddr, c *Context) PluginGroupCtrl
	Support() Method
}


type InnerPluginPathParam struct {
	Value string
}

func (pl *InnerPluginPathParam) Exec(c *Context, name string) {
	pl.Value,_ = c.PathParams.Get(name)
}

func (pl *InnerPluginPathParam) Support() Method {
	return MethodAny
}

func (pl InnerPluginPathParam) Inject(f reflect.StructField) func(base ViewAddr, c *Context) {
	offset := f.Offset
	name := f.Name
	return func(base ViewAddr, c *Context) {
		(*InnerPluginPathParam)(base.Offset(offset)).Exec(c, name)
	}
}


