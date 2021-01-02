package main

import (
	"nji"
	"nji/plugins"
)

// URL = /json

type JSON struct {
	Body struct{
		plugins.DynJSON
		A string
		B string
	}
}

func (v *JSON) Handle(c *nji.Context) {
	c.ResponseWriter.WriteHeader(200)
	_, _ = c.ResponseWriter.Write([]byte(v.Body.A+v.Body.B))
}

