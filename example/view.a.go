package main

import (
	"nji"
	"nji/plugins"
)
// URL = /a/123
type a struct {
	A plugins.PathParam
}

func (view *a) Handle(c *nji.Context) {
	c.ResponseWriter.WriteHeader(200)
	_, _ = c.ResponseWriter.Write([]byte(view.A.Value))
}


