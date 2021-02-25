package plugins

import "unsafe"

type optional bool

var _ error = err{}
type err struct {
	msg string
}

func (e err) Error() string {
	return e.msg
}

type face struct {
	typ *struct{}
	word unsafe.Pointer
}