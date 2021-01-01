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


type injectorManager struct {
	group     []PluginGroup
}

func isGroup(stru interface{}) bool{
	t := reflect.TypeOf(stru).Elem()
	length := t.NumField()
	if length == 0 {
		return false
	}
	for i := 0; i < length; i++ {
		if _, ok := reflect.New(t.Field(i).Type).Interface().(PluginGroup); ok {
			return true
		}
	}
	return false
}

func (mgr *injectorManager) ParseGroup(stru interface{}, offset uintptr) func(base ViewAddr, c *Context){
	t := reflect.TypeOf(stru).Elem()
	length := t.NumField()
	if length == 0 {
		return nil
	}
	var injectors []func(base ViewAddr, c *Context)
	pgNum := 0
	for i := 0; i < length; i++ {
		f := t.Field(i)
		if f.Type.Name() == "" {
			if f.Type.Kind().String() == "struct" {
				if stru := reflect.New(f.Type).Interface(); isGroup(stru) {
					fn := mgr.ParseGroup(stru, f.Offset)
					if fn != nil{
						injectors = append(injectors, fn)
					}
				}
			} else {
				panic("")
			}
		} else {
			if fv, ok := reflect.New(f.Type).Interface().(PluginGroup); ok {
				pgNum++
				mgr.group = append(mgr.group, fv)
				if fn := fv.Control();fn != nil{
					injectors = append(injectors, fn)
				}
				goto L
			}
		}
		// 不含 group plugin 的 无类型名 struct
		// 和非 PluginGroup 的其他类型
		for j := len(mgr.group) - 1; j >= 0; j-- {
			fn, ok := mgr.group[j].Proxy(f)
			if ok {
				if fn != nil{
					injectors = append(injectors, fn)
				}
				goto L
			}
		}
		panic("")
	L:
	}
	mgr.group = mgr.group[:len(mgr.group)-pgNum]
	return func(base ViewAddr, c *Context) {
		b := ViewAddr(uintptr(base) + offset)
		for i := 0;i<len(injectors);i++{
			injectors[i](b, c)
			if c.Error != nil{
				if _,ok := c.Error.(SkipE); ok {
					c.Error = nil
				}
				return
			}
		}
	}
}

func Inject(view ViewI) Handler {
	t := reflect.TypeOf(view).Elem()
	size := t.Size()
	length := t.NumField()
	if length == 0 {
		return nil
	}
	var injector = (&injectorManager{[]PluginGroup{&rootGroupPlugin{}}}).ParseGroup(view, 0)

	return func(c *Context) {
		addr := C.malloc(C.ulong(size))
		defer C.free(unsafe.Pointer(addr))
		C.memcpy(unsafe.Pointer(addr), (*face)(unsafe.Pointer(&view)).data, C.size_t(size))
		new_view := view
		(*face)(unsafe.Pointer(&new_view)).data = unsafe.Pointer(addr)

		injector(ViewAddr(addr), c)
		if c.Error != nil {
			_, _ = c.ResponseWriter.Write([]byte(c.Error.Error()))
			return
		}
		new_view.Handle(c)
	}
}
