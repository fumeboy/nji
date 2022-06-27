package nji

import (
	"fmt"
	"reflect"
	"strings"
)

type methodI interface {
	Method() string
}
type routeI interface {
	methodI
	Path() string
}

type Route[M methodI, R routeI] struct {
	_ R
	m M
}

func tname(s string) string {
	return strings.Split(s, "[")[0]
}

func (r Route[M, R]) Path() string {
	return ParseURL[R]()
}

func (r Route[M, R]) Method() string {
	return r.m.Method()
}

func ParseURL[R routeI]() string {
	rr := *new(R)
	t := reflect.TypeOf(rr)
	path := ""
	params := []string{}
	if t.Kind() == reflect.Struct {
		length := t.NumField()
		for i := 0; i < length; i++ {
			f := t.Field(i)
			if tname(f.Type.String()) == "plugins.PathParam" {
				params = append(params, f.Name)
			}
			v := reflect.New(f.Type).Interface()
			if fv, ok := v.(routeI); ok {
				if path != "" {
					panic("couldnt use `Route` more than one in `View` struct")
				}
				path = fv.Path()
				if string(f.Tag) == "" {
					path += "/" + t.Name()
				} else {
					path += "/" + string(f.Tag)
				}
			}
		}
	} else {
		panic("Route")
	}
	for _, p := range params {
		path += `/:` + p
	}
	return path
}

func Register[R routeI, P interface {
	Handle(*Context)
	*R
}](eg *Engine) {
	url := ParseURL[R]()
	fmt.Println(url)
	method := (*new(R)).Method()
	eg.Handle(method, url, MakeHandle[R, P]())
}
