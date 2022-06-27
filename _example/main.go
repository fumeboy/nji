package main

import "github.com/fumeboy/nji"

func main() {
	app := nji.Default()
	nji.Register[get_path_params](app)
	nji.Register[get_post_params](app)
	nji.Register[get_query_params](app)
	app.Run(":8080")
}
