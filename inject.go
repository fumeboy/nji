package nji

import "C"
import (
	"math"
	"reflect"
	"sync"
	"unsafe"
)

type View interface {
	Handle(c *Context)
}

type ViewAddress uintptr
func (c ViewAddress) Offset (o uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(c) + o)
}

type viewInPool struct {
	address unsafe.Pointer
	view View
}
var viewPools = struct {
	pools []sync.Pool
}{}


type face struct {
	tab  *struct{}
	data unsafe.Pointer
}

type slice struct {
	array unsafe.Pointer
	len   int
	cap   int
}

func parse(stru interface{}, hook func(f reflect.StructField)) inj {
	t := reflect.TypeOf(stru).Elem()
	length := t.NumField()
	if length == 0 {
		return nil
	}
	var injectors []func(base ViewAddress, c *Context)
	var method Method = -1
	for i := 0; i < length; i++ {
		f := t.Field(i)
		if f.Type.Kind().String() == "interface" {
			panic("非法的 struct field")
		}
		if hook != nil {
			hook(f)
		}
		if fv, ok := reflect.New(f.Type).Interface().(Plugin); ok {
			if fn := fv.Inject(f); fn != nil {
				method.Check(fv.Support())
				injectors = append(injectors, fn)
			}
		} else {
			panic("非法的 struct field")
		}
	}
	if len(injectors) == 0 {
		return nil
	}
	return func(b ViewAddress, c *Context) {
		for i := 0; i < len(injectors); i++ {
			injectors[i](b, c)
			if c.Error != nil {
				return
			}
		}
	}
}

func inject(default_view View, hook func(f reflect.StructField)) Handler {
	t := reflect.TypeOf(default_view).Elem()
	size := t.Size()
	length := t.NumField()
	var injector inj
	if length > 0 {
		injector = parse(default_view, hook)
	}
	viewPools.pools = append(viewPools.pools, sync.Pool{})
	p := &viewPools.pools[len(viewPools.pools)-1]
	p.New = func() interface{} {
		ret := &viewInPool{}
		b := make([]byte, size)
		ret.address = (*slice)(unsafe.Pointer(&b)).array
		ret.view = default_view
		(*face)(unsafe.Pointer(&ret.view)).data = ret.address
		return ret
	}
	return func(c *Context) {
		v := p.Get().(*viewInPool)
		copy((*[math.MaxInt32]byte)(v.address)[:size], (*[math.MaxInt32]byte)((*face)(unsafe.Pointer(&default_view)).data)[:size]) // reset
		if injector != nil {
			injector(ViewAddress(v.address), c)
			if c.Error != nil {
				_, _ = c.Resp.Writer.Write([]byte(c.Error.Error()))
				return
			}
		}
		v.view.Handle(c)
		p.Put(v)
	}
}

func Inject(view View) Handler {
	return inject(view, nil)
}
