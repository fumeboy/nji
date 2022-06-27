package main

import (
	"github.com/fumeboy/nji"
	"github.com/fumeboy/nji/route"
)

type BaseRoute struct {
	nji.Route[route.ANY, route.ROOT] `a_prefix`
}
