package main

import (
	"fmt"
	"nji"
)

func main() {
	fmt.Println(nji.MethodHead)
	app := nji.NewLazyRouter()
	app.GET(&a{})
	app.POST(&b{}, &JSON{})
	app.Run(8003)
}
