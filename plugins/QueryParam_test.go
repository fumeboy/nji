package plugins_test

import (
	"net/http"
	"net/http/httptest"
	"nji"
	"nji/plugins"
	"testing"
)

var _ nji.ViewI = &a{}

type c struct {
	A plugins.QueryParam
	B,C,D,E,F plugins.QueryParamOptional

}

func (v *c) Handle(c *nji.Context) {
	c.ResponseWriter.WriteHeader(200)
	_, _ = c.ResponseWriter.Write([]byte(v.A.Value+v.B.Value))
}

func TestContextC(t *testing.T) {
	app := nji.Config{
		UnescapePathValues: true,
	}.New()
	app.GET("/api/", nji.Inject(&c{}))
	r, err := http.NewRequest("GET", "/api/?A=Hello &B=World!", nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	w := httptest.NewRecorder()
	app.ServeHTTP(w, r)
	t.Log(string(w.Body.Bytes()))
}
