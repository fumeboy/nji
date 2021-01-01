package main

import (
	"nji"
	"nji/plugins"
)

type d struct {
	A plugins.PostParamOptional
	B plugins.PostParam
	C struct{
		plugins.GroupIgnore
		C string
	}
}

func (view *d) Handle(c *nji.Context) {
	c.ResponseWriter.WriteHeader(200)
	_, _ = c.ResponseWriter.Write([]byte(view.A.Value + view.B.Value))
}


