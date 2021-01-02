package main

import (
	"nji"
	"nji/plugins"
)

type json struct {
	Body struct{
		plugins.DynJSON
		A string
		B string
	}
}

func (v *json) Handle(c *nji.Context) {
	c.Resp.String(200,v.Body.A+v.Body.B)
}

func main() {
	app := nji.NewLazyRouter()
	app.POST(&json{})
	app.Run(8003)
}
