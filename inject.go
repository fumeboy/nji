package nji

import (
	"math"
	"reflect"
	"sync"
	"unsafe"
)

type inj = func(base unsafe.Pointer, c *Context)

type View interface {
	Handle(c *Context)
}

type inPool struct {
	address unsafe.Pointer
	typ     View // 接口变量仅作为类型信息使用， address 指向真实地址
}

var viewPools = struct {
	pools []sync.Pool
}{}

type face struct {
	itab *struct{}
	data unsafe.Pointer
}

type slice struct {
	array unsafe.Pointer
	len   int
	cap   int
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
		ret := &inPool{}
		b := make([]byte, size)
		ret.address = (*slice)(unsafe.Pointer(&b)).array
		ret.typ = default_view
		(*face)(unsafe.Pointer(&ret.typ)).data = ret.address
		return ret
	}
	return func(c *Context) {
		v := p.Get().(*inPool)
		copy((*[math.MaxInt32]byte)(v.address)[:size], (*[math.MaxInt32]byte)((*face)(unsafe.Pointer(&default_view)).data)[:size]) // reset
		if injector != nil {
			injector(v.address, c) // 执行依赖注入
			if c.err != nil {
				_, _ = c.Resp.Writer.Write([]byte(c.err.Error()))
				return
			}
		}
		v.typ.Handle(c)
		p.Put(v)
	}
}

func Inj(view View) Handler {
	return inject(view)
}
