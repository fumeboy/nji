package plugins_test

import (
	"github.com/fumeboy/nji"
	"github.com/fumeboy/nji/plugins"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var _ nji.View = &a{}

type a struct {
	A, B plugins.PathParam[any]
}

func (view *a) Handle(c *nji.Context) {
	c.Writer.WriteString(view.A.Value + view.B.Value)
}

func TestContext(t *testing.T) {
	app := nji.Default()
	app.GET("/a/:A/:B", nji.MakeHandle[a]())
	r, err := http.NewRequest("GET", "/a/Hello /World!", nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	w := httptest.NewRecorder()
	app.ServeHTTP(w, r)
	assert.Equal(t, "Hello World!", w.Body.String())
}
