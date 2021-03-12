package inject_test

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"nji"
	"nji/inject"
	"strings"
	"testing"
)

var _ nji.View = &a{}

type json struct {
	A,B inject.JSONFieldStr
}

func (v *json) Handle(c *nji.Context) {
	c.Resp.String(200,v.A.Value+v.B.Value)
}

func TestContextJSON(t *testing.T) {
	app := nji.NewServer()
	app.POST("/json", nji.Inj(new(json)))
	reader := strings.NewReader(`{"A":"Hell ", "B": "World!"}`)
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
