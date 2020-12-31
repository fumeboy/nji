package nji

// #include <stdlib.h>
// #include <string.h>
import "C"
import (
	"reflect"
	"unsafe"
)
type Plugin interface {
	Inject(f reflect.StructField) func(base ViewAddr, c *Context)
}

type ViewAddr uintptr

type ViewI interface {
	Handle(c *Context)
}

type face struct {
	tab  *struct{}
	data unsafe.Pointer
}

func (c ViewAddr) Offset (o uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(c) + o)
}

func Inject(view ViewI) Handler {
	t := reflect.TypeOf(view).Elem()
	size := t.Size()
	length := t.NumField()
	if length == 0{
		return nil
	}
	var plugins []func(base ViewAddr, c *Context)
	for i := 0;i<length;i++{
		if fv, ok := reflect.New(t.Field(i).Type).Interface().(Plugin); ok {
			plugins = append(plugins, fv.Inject(t.Field(i)))
		}
	}
	return func(c *Context) {
		addr := C.malloc(C.ulong(size))
		defer C.free(unsafe.Pointer(addr))
		C.memcpy(unsafe.Pointer(addr), (*face)(unsafe.Pointer(&view)).data, C.size_t(size))
		new_view := view
		(*face)(unsafe.Pointer(&new_view)).data = unsafe.Pointer(addr)
		for i := 0;i<len(plugins);i++{
			plugins[i](ViewAddr(addr), c)
			if c.Error != nil{
				_,_ = c.ResponseWriter.Write([]byte(c.Error.Error()))
				return
			}
		}
		new_view.Handle(c)
	}
}
