package inject_test

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"nji"
	"nji/inject"
	"nji/schema"
	"testing"
)

var _ nji.View = &custom{}

type custom struct {
	A         inject.QueryParam
	B,C,D,E,F struct{
		inject.QueryParam
		schema.IsNull
	}
}

func (v *custom) Handle(c *nji.Context) {
	c.Resp.String(200,v.A.Val+v.B.Val)
}

func TestContextCustom(t *testing.T) {
	app := nji.NewServer()
	app.GET("/c", nji.Inj(new(c)))
	r, err := http.NewRequest("GET", "/c?A=Hello &B=World!", nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	w := httptest.NewRecorder()
	app.ServeHTTP(w, r)
	assert.Equal(t, "Hello World!", string(w.Body.Bytes()))
}


