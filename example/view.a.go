package main

import (
	"net/http"
	"nji"
	"nji/plugins"
)

var _ nji.ViewI = &a{}

type a struct {
	A plugins.PathParam
}

func (view *a) Handle(c *nji.Context) {
	c.ResponseWriter.WriteHeader(200)
	_, _ = c.ResponseWriter.Write([]byte(view.A.Value))
}

func main() {
	app := nji.Config{
		UnescapePathValues: true,
		MaxMultipartMemory: 20 << 20,
	}.New()
	app.GET("/param/:A", nji.Inject(&a{})...)
	http.ListenAndServe(":8003", app)
}
