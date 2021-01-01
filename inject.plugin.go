package nji

import "reflect"
type inj = func(base ViewAddr, c *Context)

type SkipE struct {}

func (e SkipE) Error() string {
	return ""
}

type Plugin interface {
	Inject(f reflect.StructField) func(base ViewAddr, c *Context)
}

type PluginGroup interface {
	Proxy(f reflect.StructField) (fn func(base ViewAddr, c *Context), ok bool)
	Control() func(base ViewAddr, c *Context)
}

var _ PluginGroup = &rootGroupPlugin{}

type rootGroupPlugin struct {
}

func (pg *rootGroupPlugin) Proxy(f reflect.StructField) (fn func(base ViewAddr, c *Context), ok bool) {
	if fv, ok := reflect.New(f.Type).Interface().(Plugin).(Plugin); ok {
		return fv.Inject(f),true
	}
	return nil,false
}

func (pg *rootGroupPlugin) Control() func(base ViewAddr, c *Context) {
	return nil
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


