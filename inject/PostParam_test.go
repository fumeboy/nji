package inject_test

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"nji"
	"nji/inject"
	"nji/schema"
	"strings"
	"testing"
)

var _ nji.View = &b{}

type b struct {
	A         inject.PostParam
	B,C,D,E,F struct{
		inject.PostParam
		schema.IsNull
	}
}

func (v *b) Handle(c *nji.Context) {
	if v.B.IsNull {
		c.Resp.String(200,"B is null")
	}else{
		c.Resp.String(200,v.A.Val+v.B.Val)
	}
}

func TestContextB(t *testing.T) {
	app := nji.NewServer()
	app.POST("/b", nji.Inj(new(b)))
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

func TestContextB2(t *testing.T) {
	app := nji.NewServer()
	app.POST("/b", nji.Inj(new(b)))
	reader := strings.NewReader(`A=Hello`)
	r, err := http.NewRequest(http.MethodPost, "/b", reader)
	if err != nil {
		t.Error(err.Error())
		return
	}
	r.Header.Add("Content-Type","application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	app.ServeHTTP(w, r)

	assert.Equal(t, "B is null", string(w.Body.Bytes()))
}
