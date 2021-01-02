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

type a struct {
	A plugins.PathParam
	B plugins.PathParam
}

func (view *a) Handle(c *nji.Context) {
	c.Resp.String(200,view.A.Value + view.B.Value)
}

func TestContext(t *testing.T) {
	app := nji.NewLazyRouter()
	app.GET(&a{})
	r, err := http.NewRequest("GET", "/a/Hello /World!", nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	w := httptest.NewRecorder()
	app.ServeHTTP(w, r)
	assert.Equal(t, "Hello World!", string(w.Body.Bytes()))
}
