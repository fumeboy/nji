package nji

import (
	"fmt"
	"reflect"
	"unicode"
)

type lazyRouter struct { // lazyRouter 接收 View，并自动生成 URL
	*Engine
}

func NewLazyRouter() *lazyRouter{
	return &lazyRouter{Engine: NewServer()}
}

func (l *lazyRouter) GET(v ...View) {
	for i := 0;i<len(v);i++{
		l.Engine.GET(lazyRoutePath(v[i]))
	}
}
func (l *lazyRouter) POST(v ...View) {
	for i := 0;i<len(v);i++{
		l.Engine.POST(lazyRoutePath(v[i]))
	}
}

func (l *lazyRouter) PUT(v ...View) {
	for i := 0;i<len(v);i++{
		l.Engine.PUT(lazyRoutePath(v[i]))
	}
}

func (l *lazyRouter) PATCH(v ...View) {
	for i := 0;i<len(v);i++{
		l.Engine.PATCH(lazyRoutePath(v[i]))
	}
}

func (l *lazyRouter) DELETE(v ...View) {
	for i := 0;i<len(v);i++{
		l.Engine.DELETE(lazyRoutePath(v[i]))
	}
}

func (l *lazyRouter) HEAD(v ...View) {
	for i := 0;i<len(v);i++{
		l.Engine.HEAD(lazyRoutePath(v[i]))
	}
}

func (l *lazyRouter) OPTIONS(v ...View) {
	for i := 0;i<len(v);i++{
		l.Engine.OPTIONS(lazyRoutePath(v[i]))
	}
}

func lazyRoutePath(v View) (string, Handler){
	name := "/" + camel2Case(reflect.TypeOf(v).Elem().Name())
	var hook = func(f reflect.StructField){
		if _, ok := reflect.New(f.Type).Elem().Interface().(InnerPluginPathParam); ok {
			name += fmt.Sprintf("/:%s", f.Name)
		}
	}
	fn := inject(v, hook)
	return name, fn
}

func camel2Case(name string) string {
	str := []rune{}
	f := false
	for _, r := range name {
		if unicode.IsUpper(r) {
			if f {
				str = append(str, '/')
				f = false
			}
			str = append(str, (unicode.ToLower(r)))
		} else {
			f = true
			str = append(str, (r))
		}
	}
	return string(str)
}