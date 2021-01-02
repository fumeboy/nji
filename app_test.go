package nji

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// 404
func TestNotFoundEvent(t *testing.T) {
	app := NewServer()
	r, err := http.NewRequest("GET", "/404", nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	w := httptest.NewRecorder()
	app.ServeHTTP(w, r)
	assert.Equal(t, w.Body.String(), "404 not found")
}

// 500
func TestPanicEvent(t *testing.T) {
	app := NewServer()
	app.GET("/", func(ctx *Context) {
		panic("")
	})
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	w := httptest.NewRecorder()
	app.ServeHTTP(w, r)
	assert.Equal(t, w.Body.String(), "Internal Server Error")
}
