package inject_test

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"nji"
	"nji/inject"
	"testing"
)

var _ nji.View = &a{}

type a struct {
	A, B inject.PathParam
}

func (view *a) Handle(c *nji.Context) {
	c.Resp.String(200,view.A.Value + view.B.Value)
}

func TestContext(t *testing.T) {
	app := nji.NewServer()
	app.GET("/a/:A/:B", nji.Inj(new(a)))
	r, err := http.NewRequest("GET", "/a/Hello /World!", nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	w := httptest.NewRecorder()
	app.ServeHTTP(w, r)
	assert.Equal(t, "Hello World!", string(w.Body.Bytes()))
}
