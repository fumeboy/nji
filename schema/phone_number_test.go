package schema_test

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"nji"
	"nji/inject"
	"nji/schema"
	"testing"
)

var _ nji.View = &a{}

type a struct {
	A struct {
		inject.PathParam
		schema.IsPhoneNumber
	}
}

func (view *a) Handle(c *nji.Context) {
	c.Resp.String(200, view.A.Value)
}

func TestContext(t *testing.T) {
	app := nji.NewServer()
	app.GET("/a/:A", nji.Inj(new(a)))
	r, err := http.NewRequest("GET", "/a/17979300086", nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	w := httptest.NewRecorder()
	app.ServeHTTP(w, r)
	assert.Equal(t, "17979300086", string(w.Body.Bytes()))
}
