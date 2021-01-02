package plugins_test

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"nji"
	"nji/plugins"
	"testing"
)

var _ nji.ViewI = &a{}

type ig struct {
	Body struct{
		plugins.DynIgnore
		A plugins.QueryParam
		B plugins.QueryParamOptional
	}
}

func (v *ig) Handle(c *nji.Context) {
	c.ResponseWriter.WriteHeader(200)
	_, _ = c.ResponseWriter.Write([]byte(v.Body.A.Value+v.Body.B.Value))
}

func TestContextI(t *testing.T) {
	app := nji.NewLazyRouter()
	app.GET(&ig{})
	r, err := http.NewRequest("GET", "/ig?A=Hello &B=World!", nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	w := httptest.NewRecorder()
	app.ServeHTTP(w, r)
	assert.Equal(t, "", string(w.Body.Bytes()))
}
