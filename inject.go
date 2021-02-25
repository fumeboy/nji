package nji

import (
	"fmt"
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

type field struct {
	iface Plugin
	startAt uintptr
	len uintptr
	f reflect.StructField
}
type fieldDyn struct {
	fn PluginDyn
	iface interface{}
	startAt uintptr
	len uintptr
	f reflect.StructField
}

func parse(stru interface{}, hook func(f reflect.StructField)) inj {
	t := reflect.TypeOf(stru).Elem()
	length := t.NumField()
	if length == 0 {
		return nil
	}
	var method Method = -1
	var fields []*field
	var fieldsDyn []*fieldDyn
	for i := 0; i < length; i++ {
		f := t.Field(i)
		if f.Type.Kind().String() == "interface" {
			panic("非法的 struct field")
		}
		if hook != nil {
			hook(f)
		}
		if f.Type.Kind().String() == "Ptr"{
			panic("非法的 struct field")
		}
		if fv, ok := reflect.New(f.Type).Interface().(Plugin); ok {
			fields = append(fields, &field{
				iface:   fv,
				startAt: f.Offset,
				len:     f.Type.Size(),
				f: f,
			})
			method.Check(fv.Support())
		} else if fv, ok := reflect.New(f.Type).Interface().(PluginDyn); ok {
			fieldsDyn = append(fieldsDyn, &fieldDyn{
				fn: fv,
				iface:   reflect.New(f.Type).Interface(),
				startAt: f.Offset,
				len:     f.Type.Size(),
				f:       f,
			})
			method.Check(fv.Support())
		} else {
			fmt.Println(f)
			panic("非法的 struct field")
		}
	}
	if len(fields) == 0 && len(fieldsDyn) == 0{
		return nil
	}
	return func(b ViewAddress, c *Context) {
		for i := 0; i < len(fields); i++ {
			ifa := fields[i].iface
			(*face)(unsafe.Pointer(&ifa)).data = b.Offset(fields[i].startAt)
			if err := ifa.Inject(c, fields[i].f); err != nil {
				c.Error = err
				return
			}
		}
		for i := 0; i < len(fieldsDyn); i++ {
			ifa := fieldsDyn[i].iface
			(*face)(unsafe.Pointer(&ifa)).data = b.Offset(fieldsDyn[i].startAt)
			if err := fieldsDyn[i].fn.Inject(c, fieldsDyn[i].f, ifa); err != nil {
				c.Error = err
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
