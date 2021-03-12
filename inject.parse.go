package nji

import (
	"fmt"
	"nji/schema"
	"reflect"
	"unsafe"
)

type Plugin interface {
	Inject(c *Context, f reflect.StructField) error
}

var IsNullTyp = reflect.TypeOf((schema.IsNull)(true))

type field struct {
	typ            Plugin
	offset         uintptr
	len            uintptr
	f              reflect.StructField
	validators     []schema.RealV
	canBeNull      bool
	nullFlagOffset uintptr
}

func parse(stru interface{}) inj {
	t := reflect.TypeOf(stru).Elem()
	length := t.NumField()
	if length == 0 {
		return nil
	}
	var fields []*field
	for i := 0; i < length; i++ {
		f := t.Field(i)
		if f.Type.Kind() == reflect.Interface {
			panic("非法的 struct field")
		}
		if f.Type.Kind() == reflect.Ptr {
			panic("非法的 struct field")
		}

		if fv, ok := reflect.New(f.Type).Interface().(Plugin); ok {
			if f.Type.Kind() == reflect.Struct && f.Type.Name() == "" { // 匿名结构体，即使用了 schema 的 plugin
				ff := &field{
					offset: f.Offset,
					f:      f,
				}
				if fvv, ok := reflect.New(f.Type.Field(0).Type).Interface().(Plugin); ok {
					ff.typ = fvv
					ff.len = f.Type.Field(0).Type.Size()
				} else {
					panic("匿名结构体的首个 field 必须是 plugin 类型")
				}
				for i := 1; i < f.Type.NumField(); i++ { // 解析 schema field
					if v, ok := reflect.New(f.Type.Field(i).Type).Interface().(schema.V); ok {
						realv := v.INJ()
						rvt := reflect.ValueOf(realv).Type()
						if rvt.Kind() == reflect.Ptr {
							rvt = rvt.Elem()
						}
						if rvt.Size() != ff.len {
							panic("plugin 和 validator 不兼容")
						}
						ff.validators = append(ff.validators, realv)
					} else if f.Type.Field(i).Type == IsNullTyp {
						ff.canBeNull = true
						ff.nullFlagOffset = f.Type.Field(i).Offset
					} else {
						panic("非法的 struct field")
					}
				}
				fields = append(fields, ff)
			} else {
				fields = append(fields, &field{
					typ:    fv,
					offset: f.Offset,
					len:    f.Type.Size(),
					f:      f,
				})
			}
		} else {
			fmt.Println(f)
			panic("非法的 struct field")
		}
	}
	if len(fields) == 0 {
		return nil
	}
	return func(base unsafe.Pointer, c *Context) {
		for i := 0; i < len(fields); i++ {
			ifa := fields[i].typ
			(*face)(unsafe.Pointer(&ifa)).data = unsafe.Pointer(fields[i].offset +uintptr(base))
			if err := ifa.Inject(c, fields[i].f); err != nil { // 执行依赖注入
				if fields[i].canBeNull {
					*(*bool)(unsafe.Pointer(uintptr(base)+fields[i].offset + fields[i].nullFlagOffset)) = true
				} else {
					c.err = err
				}
				return
			}
			for j := 0; j < len(fields[i].validators); j++ { // 执行 schema 校验
				v := fields[i].validators[j]
				(*face)(unsafe.Pointer(&v)).data = unsafe.Pointer(uintptr(base)+fields[i].offset)
				if err := v.Check(); err != nil {
					c.err = err
					return
				}
			}
		}
	}
}
