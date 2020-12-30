package main

import "nji"

func main() {
	app := nji.Config{
		UnescapePathValues: true,
	}.New()
	app.GET("/param/:A", nji.Inject(&a{}))
	app.Run(8003)
}
