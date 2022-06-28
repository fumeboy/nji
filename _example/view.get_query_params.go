package main

import (
	"github.com/fumeboy/nji"
	"github.com/fumeboy/nji/plugins"
	"github.com/fumeboy/nji/route"
	"github.com/fumeboy/nji/schema"
)

var _ nji.View = &get_query_params{}

// please visit http://127.0.0.1:8080/get_query_params?A=phonenumberis&B=12345678901
type get_query_params struct {
	nji.Route[route.GET, route.ROOT]

	A plugins.QueryParam[schema.NotNull]
	B plugins.QueryParam[struct {
		schema.NotNull
		schema.IsPhoneNumber
	}]
}

func (v *get_query_params) Handle(c *nji.Context) {
	c.Writer.WriteString(v.A.Value + v.B.Value)
}
