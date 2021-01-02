package nji_test

import (
	"nji"
	"nji/plugins"
	"testing"
)

type b struct {
	A plugins.PostParam
	B plugins.QueryParam
}

func (v *b) Handle(c *nji.Context) {
}

func TestMethodCheck(t *testing.T) {
	app := nji.NewLazyRouter()
	app.POST(&b{})
}