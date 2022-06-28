package main

import (
	"github.com/fumeboy/nji"
	"github.com/fumeboy/nji/plugins"
	"github.com/fumeboy/nji/route"
	"github.com/fumeboy/nji/schema"
)

var _ nji.View = &get_post_params{}

// curl -d "A=1234&C=abcd" -X POST http://127.0.0.1:8080/a_prefix/get_post_params1234
// curl -d "A=1234&B=abcd" -X POST http://127.0.0.1:8080/a_prefix/get_post_params1234
// curl -d "C=1234&B=abcd" -X POST http://127.0.0.1:8080/a_prefix/get_post_params1234
type get_post_params struct {
	nji.Route[route.POST, BaseRoute] `get_post_params1234`

	A             plugins.PostParam[schema.NotNull]
	B, C, D, E, F plugins.PostParam[any]
}

func (v *get_post_params) Handle(c *nji.Context) {
	if v.B.IsNull {
		c.Writer.WriteString("B is null")
	} else {
		c.Writer.WriteString(v.A.Value + v.B.Value)
	}
}
