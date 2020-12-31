package plugins_test

import (
	"net/http"
	"net/http/httptest"
	"nji"
	"nji/plugins"
	"testing"
)

var _ nji.ViewI = &a{}

type b struct {
	A plugins.PostParam
	B plugins.PostParamOptional
}

func (view *b) Handle(c *nji.Context) {
	c.ResponseWriter.WriteHeader(200)
	_, _ = c.ResponseWriter.Write([]byte("Hello !"))
}

func TestContextB(t *testing.T) {
	app := nji.Config{
		UnescapePathValues: true,
	}.New()
	app.POST("/api/", nji.Inject(&b{}))
	r, err := http.NewRequest("POST", "/api/", nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	app.ServeHTTP(httptest.NewRecorder(), r)
}
