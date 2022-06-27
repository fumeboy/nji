package nji

import (
	"reflect"

	"github.com/gin-gonic/gin"
)

type Context = gin.Context
type Engine = gin.Engine

var Default = gin.Default

type HandlerFunc = gin.HandlerFunc

type View interface {
	Handle(c *Context)
}

type Plugin interface {
	Run(c *Context, f reflect.StructField) error
}

type PluginBuilder interface {
	Build(f reflect.StructField) error
}
