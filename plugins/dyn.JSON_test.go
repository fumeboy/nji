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

var _ nji.ViewI = &a{}

type json struct {
	Body struct{
		plugins.DynJSON
		A string
		B string
	}
}

func (v *json) Handle(c *nji.Context) {
	c.ResponseWriter.WriteHeader(200)
	_, _ = c.ResponseWriter.Write([]byte(v.Body.A+v.Body.B))
}

func TestContextJSON(t *testing.T) {
	app := nji.NewLazyRouter()
	app.POST(&json{})
	reader := strings.NewReader(`{"A":"Hello ", "B": "World!"}`)
	r, err := http.NewRequest(http.MethodPost, "/json", reader)
	if err != nil {
		t.Error(err.Error())
		return
	}
	r.Header.Add("Content-Type","application/json")
	w := httptest.NewRecorder()
	app.ServeHTTP(w, r)
	assert.Equal(t, "Hello World!", string(w.Body.Bytes()))
}
