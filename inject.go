package nji

// #include <stdlib.h>
// #include <string.h>
import "C"
import (
	"reflect"
	"unsafe"
)
type Plugin interface {
	Inject(f reflect.StructField) func(c *Context)
}

type ViewI interface {
	Handle(c *Context)
}

type face struct {
	tab  *struct{}
	data unsafe.Pointer
}

func (c *Context) OffsetV (o uintptr) unsafe.Pointer {
	return unsafe.Pointer(c.viewAddr + o)
}

func Inject(view ViewI) HandlersChain {
	t := reflect.TypeOf(view).Elem()
	size := t.Size()
	length := t.NumField()
	if length == 0{
		return nil
	}
	var handlers = HandlersChain{
		func(c *Context) {
			c.viewAddr = uintptr(C.malloc(C.ulong(size)))
			defer C.free(unsafe.Pointer(c.viewAddr))
			C.memcpy(unsafe.Pointer(c.viewAddr), (*face)(unsafe.Pointer(&view)).data, C.size_t(size))
			c.view = view
			(*face)(unsafe.Pointer(&c.view)).data = unsafe.Pointer(c.viewAddr)
			c.Next()
			if !c.IsAborted() {
				c.view.Handle(c)
			}
		},
	}
	for i := 0;i<length;i++{
		if fv, ok := reflect.New(t.Field(i).Type).Interface().(Plugin); ok {
			handlers = append(handlers, fv.Inject(t.Field(i)))
		}
	}
	return handlers
}
