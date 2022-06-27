package plugins

var _ error = err{}

type err struct {
	msg string
}

func (e err) Error() string {
	return e.msg
}
