package plugins_test

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"nji"
	"nji/plugins"
	"strings"
	"testing"
)

var _ nji.View = &a{}

type b struct {
	A plugins.PostParam
	B,C,D,E,F plugins.PostParamOptional
}

func (v *b) Handle(c *nji.Context) {
	c.Resp.String(200,v.A.Value+v.B.Value)
}

func TestContextB(t *testing.T) {
	app := nji.NewLazyRouter()
	app.POST(&b{})
	reader := strings.NewReader(`A=Hello &B=World!`)
	r, err := http.NewRequest(http.MethodPost, "/b", reader)
	if err != nil {
		t.Error(err.Error())
		return
	}
	r.Header.Add("Content-Type","application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	app.ServeHTTP(w, r)

	assert.Equal(t, "Hello World!", string(w.Body.Bytes()))
}
