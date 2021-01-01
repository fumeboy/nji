package main

import "nji"

func main() {
	app := nji.Config{
		UnescapePathValues: true,
	}.New()
	app.GET("/get/:A", nji.Inject(&a{}))
	app.POST("/post", nji.Inject(&b{}))
	app.Run(8003)
}
