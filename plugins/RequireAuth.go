package plugins

import (
	"nji"
	"reflect"
)

type RequireAuth struct {
	*Claims
}

func (pl *RequireAuth) Exec(c *nji.Context) (err error) {
	token := c.Request.Header.Get("Authorization")[4:] // `JWT YWxhZGRpbjpvcGVuc2VzYW1l`
	pl.Claims, err = parseToken(token)
	return
}

func (pl *RequireAuth) Support() nji.Method {
	return nji.MethodAny
}

func (pl RequireAuth) Inject(f reflect.StructField) func(base nji.ViewAddr, c *nji.Context) {
	offset := f.Offset
	return func(base nji.ViewAddr, c *nji.Context) {
		c.Error = (*RequireAuth)(base.Offset(offset)).Exec(c)
	}
}

