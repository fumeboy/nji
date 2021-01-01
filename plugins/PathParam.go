package plugins

import (
	"nji"
)

var _ nji.Plugin = &PathParam{}

type PathParam = nji.InnerPluginPathParam