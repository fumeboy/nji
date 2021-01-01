package plugins

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