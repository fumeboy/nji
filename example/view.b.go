package main

import (
	"nji"
	"nji/plugins"
)

type b struct {
	A plugins.PostParamOptional
	B plugins.PostParam
}

func (view *b) Handle(c *nji.Context) {
	c.ResponseWriter.WriteHeader(200)
	_, _ = c.ResponseWriter.Write([]byte(view.A.Value + view.B.Value))
}


