package nji

import "reflect"

type inj = func(base ViewAddress, c *Context)

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

func (method *Method) Check(m Method){
	if *method == -1 {
		*method = m
	} else if m&*method == 0 {
		panic("请检查插件是否可以放在一起使用")
	} else {
		*method = m & *method
	}
}

type Plugin interface {
	Inject(f reflect.StructField) func(base ViewAddress, c *Context)
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

func (pl InnerPluginPathParam) Inject(f reflect.StructField) func(base ViewAddress, c *Context) {
	offset := f.Offset
	name := f.Name
	return func(base ViewAddress, c *Context) {
		(*InnerPluginPathParam)(base.Offset(offset)).Exec(c, name)
	}
}


