package plugins_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"nji"
	"nji/plugins"
	"testing"
	"unsafe"
)

var _ nji.ViewI = &a{}

type a struct {
	A plugins.PathParam
}

func (view *a) Handle(c *nji.Context) {
	fmt.Println(234, unsafe.Pointer(view))
	c.ResponseWriter.WriteHeader(200)
	_, _ = c.ResponseWriter.Write([]byte("Hello !"))
}

func TestContext(t *testing.T) {
	app := nji.Config{
		UnescapePathValues: true,
		MaxMultipartMemory: 20 << 20,
	}.New()
	app.GET("/context/:A", nji.Inject(&a{}))
	r, err := http.NewRequest("GET", "/context/2", nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	app.ServeHTTP(httptest.NewRecorder(), r)
}
