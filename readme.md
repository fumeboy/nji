# nji

a web framework for write API quickly.

# feature

* simplify routing management
* parameter verification without go tag
* configurable dependency injection

# example

the user could write HTTP handler like this:

```go
// ./_example/view.get_query_params.go
type get_query_params struct {
	nji.Route[route.GET, route.ROOT] // define URL
    // and could be visit by [GET]http://127.0.0.1:8080/get_query_params?A=phonenumberis&B=12345678901


	A plugins.QueryParam[any] // define `plugin` to inject args automatically
	B plugins.QueryParam[struct {
        // use generic T as metadata for parameter verification
		schema.Must // must not null
		schema.IsPhoneNumber // check if valid as phonenumber
	}]
}

func (v *get_query_params) Handle(c *nji.Context) {
	c.Writer.WriteString(v.A.Value + v.B.Value)
}

func main() {
	app := nji.Default()
	nji.Register[get_query_params](app) 
    // `Register` get URL from `nji.Route`, so dont need write URL manualy
	app.Run(":8080")
}
```

instead of:

```go
func get_query_params(c *gin.Context) {
    a, ok := c.GetQuery("A")
    if !ok {
        // ...
    }
    b, ok := c.GetQuery("B")
    if !ok {
        // ...
    }
    if checkIsPhoneNumber(b) {
        // ...
    }

    c.Writer.WriteString(v.A.Value + v.B.Value)
}

func main() {
	router := gin.Default()
	router.POST("/get_query_params", get_query_params)
	router.Run(":8080")
}
```

please visit `./_example` and `./plugins/*_test.go` for more examples 









