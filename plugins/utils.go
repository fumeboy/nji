package plugins

import "unsafe"

type optional struct {
	notEmpty bool
}

func (e optional) NotEmpty() bool {
	return e.notEmpty
}

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