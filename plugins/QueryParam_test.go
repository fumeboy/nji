package plugins_test

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"nji"
	"nji/plugins"
	"testing"
)

var _ nji.View = &a{}

type c struct {
	A plugins.QueryParam
	B,C,D,E,F plugins.QueryParamOptional

}

func (v *c) Handle(c *nji.Context) {
	c.Resp.String(200,v.A.Value+v.B.Value)
}

func TestContextC(t *testing.T) {
	app := nji.NewLazyRouter()
	app.GET(&c{})
	r, err := http.NewRequest("GET", "/c?A=Hello &B=World!", nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	w := httptest.NewRecorder()
	app.ServeHTTP(w, r)
	assert.Equal(t, "Hello World!", string(w.Body.Bytes()))
}
