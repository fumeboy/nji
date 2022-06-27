package nji

import (
	"reflect"
	"unsafe"
)

type iface struct {
	itab *struct{}
	data unsafe.Pointer
}

type field struct {
	typecast Plugin
	offset   uintptr
	f        reflect.StructField
}

func MakeHandle[V any, P interface {
	Handle(*Context)
	*V
}]() HandlerFunc {
	defaultv := *new(V)
	t := reflect.TypeOf(defaultv)
	length := t.NumField()
	var injector func(base unsafe.Pointer, c *Context) error
	if length > 0 {
		var fields []*field
		for i := 0; i < length; i++ {
			f := t.Field(i)
			if f.Type.Kind() == reflect.Interface {
				panic("invalid struct field")
			}
			if f.Type.Kind() == reflect.Ptr {
				panic("invalid struct field")
			}
			v := reflect.New(f.Type).Interface()
			if _, ok := v.(routeI); ok {
				continue
			}
			if fv, ok := v.(Plugin); ok {
				if fv2, ok := v.(PluginBuilder); ok {
					(*iface)(unsafe.Pointer(&fv2)).data = unsafe.Pointer(f.Offset + uintptr(unsafe.Pointer(&defaultv)))
					if err := fv2.Build(f); err != nil {
						panic(err)
					}
				}
				fields = append(fields, &field{
					typecast: fv,
					offset:   f.Offset,
					f:        f,
				})
			} else {
				panic("invalid struct field")
			}
		}
		if len(fields) > 0 {
			injector = func(base unsafe.Pointer, c *Context) error {
				for i := 0; i < len(fields); i++ {
					ifa := fields[i].typecast
					(*iface)(unsafe.Pointer(&ifa)).data = unsafe.Pointer(fields[i].offset + uintptr(base))
					err := ifa.Run(c, fields[i].f)
					if err != nil {
						return err
					}
				}
				return nil
			}
		}
	}
	return func(c *Context) {
		v := defaultv
		if injector != nil {
			if err := injector(unsafe.Pointer(&v), c); err != nil {
				c.Writer.WriteString(err.Error())
				return
			}
		}
		P(&v).Handle(c)
	}
}
