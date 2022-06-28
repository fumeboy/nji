package plugins_test

import (
	"github.com/fumeboy/nji"
	"github.com/fumeboy/nji/plugins"
	"github.com/fumeboy/nji/schema"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var _ nji.View = &b{}

type b struct {
	A             plugins.PostParam[any]
	B, C, D, E, F plugins.PostParam[struct {
		schema.NotNull
	}]
}

func (v *b) Handle(c *nji.Context) {
	if v.B.IsNull {
		c.Writer.WriteString("B is null")
	} else {
		c.Writer.WriteString(v.A.Value + v.B.Value)
	}
}

func TestContextB(t *testing.T) {
	app := nji.Default()
	app.POST("/b", nji.MakeHandle[b]())
	reader := strings.NewReader(`A=Hello &B=World!`)
	r, err := http.NewRequest(http.MethodPost, "/b", reader)
	if err != nil {
		t.Error(err.Error())
		return
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	app.ServeHTTP(w, r)

	assert.Equal(t, "Hello World!", w.Body.String())
}

func TestContextB2(t *testing.T) {
	app := nji.Default()
	app.POST("/b", nji.MakeHandle[b]())
	reader := strings.NewReader(`A=Hello`)
	r, err := http.NewRequest(http.MethodPost, "/b", reader)
	if err != nil {
		t.Error(err.Error())
		return
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	app.ServeHTTP(w, r)

	assert.Equal(t, "B is null", w.Body.String())
}
