package nji

// #include <stdlib.h>
// #include <string.h>
import "C"
import (
	"reflect"
	"unsafe"
)

type face struct {
	tab  *struct{}
	data unsafe.Pointer
}

func parseGroup(stru PluginGroup, offset uintptr, method *Method, hook func(f reflect.StructField)) inj {
	t := reflect.TypeOf(stru).Elem()
	length := t.NumField()
	if length == 0 {
		return nil
	}
	var injectors []func(base ViewAddr, c *Context)
	var ctrl = stru.InjectAndControl(t.Field(0))
	for i := 1; i < length; i++ {
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
				injectors = append(injectors, fv.Inject(f))
			}
		} else {
			panic("非法的 struct field")
		}
	}
	if ctrl == nil && len(injectors) == 0 {
		return nil
	}

	if ctrl != nil {
		return func(base ViewAddr, c *Context) {
			b := ViewAddr(uintptr(base) + offset)
			if ctrl(b, c) > PluginGroupCtrlSuccess {
				return
			}
			for i := 0; i < len(injectors); i++ {
				injectors[i](b, c)
				if c.Error != nil {
					return
				}
			}
		}
	}

	return func(base ViewAddr, c *Context) {
		b := ViewAddr(uintptr(base) + offset)
		for i := 0; i < len(injectors); i++ {
			injectors[i](b, c)
			if c.Error != nil {
				return
			}
		}
	}
}

func parse(stru interface{}, hook func(f reflect.StructField)) inj {
	t := reflect.TypeOf(stru).Elem()
	length := t.NumField()
	if length == 0 {
		return nil
	}
	var injectors []func(base ViewAddr, c *Context)
	var method Method = -1
	for i := 0; i < length; i++ {
		f := t.Field(i)
		if f.Type.Kind().String() == "interface" {
			panic("非法的 struct field")
		}
		if hook != nil {
			hook(f)
		}
		if fv, ok := reflect.New(f.Type).Interface().(PluginGroup); ok {
			method.Check(fv.Support())
			if f.Name == "" { // group plugin 在匿名结构体中，成组使用
				if fn := parseGroup(fv, f.Offset, &method, hook); fn != nil {
					injectors = append(injectors, fn)
				}
			} else {
				if fn := fv.InjectAndControl(f); fn != nil { // group plugin 不成组， 单独使用
					injectors = append(injectors, func(base ViewAddr, c *Context) {
						fn(base, c)
					})
				}
			}
		} else if fv, ok := reflect.New(f.Type).Interface().(Plugin); ok {
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
	return func(b ViewAddr, c *Context) {
		for i := 0; i < len(injectors); i++ {
			injectors[i](b, c)
			if c.Error != nil {
				return
			}
		}
	}
}

func inject(view ViewI, hook func(f reflect.StructField)) Handler {
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
			injector(ViewAddr(addr), c)
			if c.Error != nil {
				_, _ = c.Resp.Writer.Write([]byte(c.Error.Error()))
				return
			}
		}
		new_view.Handle(c)
	}
}

func Inject(view ViewI) Handler {
	return inject(view, nil)
}
