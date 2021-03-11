package nji

import (
	"fmt"
	"math"
	"nji/schema"
	"reflect"
	"sync"
	"unsafe"
)

type inj = func(base ViewAddress, c *Context)

type Plugin interface {
	Inject(c *Context, f reflect.StructField) error
}

type View interface {
	Handle(c *Context)
}

type ViewAddress uintptr

func (c ViewAddress) Offset(o uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(c) + o)
}

type viewInPool struct {
	address unsafe.Pointer
	view    View
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
	iface          Plugin
	startAt        uintptr
	len            uintptr
	f              reflect.StructField
	validators     []schema.RealV
	canBeNull      bool
	nullFlagOffset uintptr
}

var IsNullTyp = reflect.TypeOf((schema.IsNull)(true))

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
			if f.Type.Kind() == reflect.Struct && f.Type.Name() == "" {
				ff := &field{
					startAt: f.Offset,
					f:       f,
				}
				if fvv, ok := reflect.New(f.Type.Field(0).Type).Interface().(Plugin); ok {
					ff.iface = fvv
					ff.len = f.Type.Field(0).Type.Size()
				} else {
					panic("匿名结构体的首个 field 必须是 plugin 类型")
				}
				for i := 1; i < f.Type.NumField(); i++ {
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
					iface:   fv,
					startAt: f.Offset,
					len:     f.Type.Size(),
					f:       f,
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
	return func(b ViewAddress, c *Context) {
		for i := 0; i < len(fields); i++ {
			ifa := fields[i].iface
			(*face)(unsafe.Pointer(&ifa)).data = b.Offset(fields[i].startAt)
			if err := ifa.Inject(c, fields[i].f); err != nil {
				if fields[i].canBeNull {
					*(*bool)(b.Offset(fields[i].startAt + fields[i].nullFlagOffset)) = true
				} else {
					c.err = err
				}
				return
			}
			for j := 0; j < len(fields[i].validators); j++ {
				v := fields[i].validators[j]
				(*face)(unsafe.Pointer(&v)).data = b.Offset(fields[i].startAt)
				if err := v.Check(); err != nil {
					c.err = err
					return
				}
			}
		}
	}
}

func inject(default_view View) Handler {
	t := reflect.TypeOf(default_view).Elem()
	size := t.Size()
	length := t.NumField()
	var injector inj
	if length > 0 {
		injector = parse(default_view)
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
			if c.err != nil {
				_, _ = c.Resp.Writer.Write([]byte(c.err.Error()))
				return
			}
		}
		v.view.Handle(c)
		p.Put(v)
	}
}

func Inj(view View) Handler {
	return inject(view)
}
