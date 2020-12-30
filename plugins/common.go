package plugins

type notEmpty struct {
	value bool
}

func (e notEmpty) NotEmpty() bool {
	return e.value
}