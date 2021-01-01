package nji

import "unsafe"

type ViewAddr uintptr

type ViewI interface {
	Handle(c *Context)
}

func (c ViewAddr) Offset (o uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(c) + o)
}

