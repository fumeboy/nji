package main

import (
	"nji"
	"nji/plugins"
)

type d struct {
	Body struct{
		plugins.DynIgnore
		A plugins.PostParamOptional
		B plugins.PostParam
	}
}

func (view *d) Handle(c *nji.Context) {
	c.ResponseWriter.WriteHeader(200)
	_, _ = c.ResponseWriter.Write([]byte(view.Body.A.Value + view.Body.B.Value))
}


