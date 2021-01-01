package plugins_test

import (
	"net/http"
	"net/http/httptest"
	"nji"
	"nji/plugins"
	"strings"
	"testing"
)

var _ nji.ViewI = &a{}

type b struct {
	A plugins.PostParam
	B,C,D,E,F plugins.PostParamOptional
}

func (v *b) Handle(c *nji.Context) {
	c.ResponseWriter.WriteHeader(200)
	_, _ = c.ResponseWriter.Write([]byte(v.A.Value+v.B.Value))
}

func TestContextB(t *testing.T) {
	app := nji.Config{
		UnescapePathValues: true,
	}.New()
	app.POST("/api/", nji.Inject(&b{}))
	reader := strings.NewReader(`A=Hello &B=World!`)
	r, err := http.NewRequest(http.MethodPost, "/api/", reader)
	if err != nil {
		t.Error(err.Error())
		return
	}
	r.Header.Add("Content-Type","application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	app.ServeHTTP(w, r)

	t.Log(string(w.Body.Bytes()))
}
