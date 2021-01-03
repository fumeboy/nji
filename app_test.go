package nji_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"nji"
	"nji/plugins"
	"testing"
	"time"
)

// 404
func TestNotFoundEvent(t *testing.T) {
	app := nji.NewServer()
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
	app := nji.NewServer()
	app.GET("/", func(ctx *nji.Context) {
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

type mock struct {
	A plugins.PathParam
}

func (v *mock) Handle(c *nji.Context) {
	c.Resp.String(200,v.A.Value)
}

func TestMock(t *testing.T) {
	app := nji.NewLazyRouter()
	app.GET(&mock{})
	go app.Run(8003)

	go func() {
		for i := 0;i<1000;i++{
			go func() {
				buf := make([]byte, 100)
				for j := 0;j<10;j++{
					r, err := http.Get(fmt.Sprintf("http://127.0.0.1:8003/mock/12"))
					if err == nil{
						n,err := r.Body.Read(buf)
						if err != io.EOF{
							fmt.Println(err)
							return
						}
						if "12" != string(buf[:n]){
							fmt.Println(123)
							return
						}
					}
				}
			}()
		}
	}()
	time.Sleep(19*time.Second)
}