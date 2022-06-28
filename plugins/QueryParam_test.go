package plugins_test

import (
	"github.com/fumeboy/nji"
	"github.com/fumeboy/nji/plugins"
	"github.com/fumeboy/nji/route"
	"github.com/fumeboy/nji/schema"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ nji.View = &c{}

type c struct {
	nji.Route[route.GET, route.ROOT]

	A plugins.QueryParam[schema.NotNull]
	B plugins.QueryParam[func(schema.NotNull, schema.IsPhoneNumber)]
}

func (v *c) Handle(c *nji.Context) {
	c.Writer.WriteString(v.A.Value + v.B.Value)
}

func TestContextC(t *testing.T) {
	app := nji.Default()
	nji.Register[c](app)
	r, err := http.NewRequest("GET", "/c?A=PhoneNumberIs &B=12345678912", nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	w := httptest.NewRecorder()
	app.ServeHTTP(w, r)
	assert.Equal(t, "PhoneNumberIs 12345678912", w.Body.String())
}
