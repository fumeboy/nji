package nji

import (
	"errors"
	"reflect"
)

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
	Inject(c *Context, f reflect.StructField) error
	Support() Method
}

type PluginDyn interface {
	Inject(c *Context, f reflect.StructField, iface interface{}) error
	Support() Method
}

type InnerPluginPathParam struct {
	Value string
}

func (pl *InnerPluginPathParam) Support() Method {
	return MethodAny
}

func (pl *InnerPluginPathParam) Inject(c *Context, f reflect.StructField) error{
	var ok bool
	pl.Value,ok = c.PathParams.Get(f.Name)
	if ok {
		return nil
	}
	return errors.New("bad path")
}


