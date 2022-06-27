package main

import (
	"github.com/fumeboy/nji"
	"github.com/fumeboy/nji/plugins"
	"github.com/fumeboy/nji/route"
)

// please visit http://127.0.0.1:8080/a_prefix/get_path_params/123/456
type get_path_params struct {
	nji.Route[route.GET, BaseRoute]

	A, B plugins.PathParam[any]
}

func (v *get_path_params) Handle(c *nji.Context) {
	c.Writer.WriteString(v.A.Value + v.B.Value)
}
