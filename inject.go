package nji

// #include <stdlib.h>
// #include <string.h>
import "C"
import (
	"reflect"
	"unsafe"
)

type View interface {
	Handle(c *Context)
}
type ViewAddress uintptr
func (c ViewAddress) Offset (o uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(c) + o)
}

type face struct {
	tab  *struct{}
	data unsafe.Pointer
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

func inject(view View, hook func(f reflect.StructField)) Handler {
	t := reflect.TypeOf(view).Elem()
	size := t.Size()
	length := t.NumField()
	var injector inj
	if length > 0 {
		injector = parse(view, hook)
	}

	return func(c *Context) {
		addr := C.malloc(C.ulong(size))
		defer C.free(unsafe.Pointer(addr))
		C.memcpy(unsafe.Pointer(addr), (*face)(unsafe.Pointer(&view)).data, C.size_t(size))
		new_view := view
		(*face)(unsafe.Pointer(&new_view)).data = unsafe.Pointer(addr)

		if injector != nil {
			injector(ViewAddress(addr), c)
			if c.Error != nil {
				_, _ = c.Resp.Writer.Write([]byte(c.Error.Error()))
				return
			}
		}
		new_view.Handle(c)
	}
}

func Inject(view View) Handler {
	return inject(view, nil)
}
