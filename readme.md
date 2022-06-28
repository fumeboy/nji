# nji

a web framework for write API quickly.

# feature

* simplify routing management
* parameter verification without go tag
* configurable dependency injection
# basic usage

## get parameters automatically with `plugins`

```go
// get params (`A`, `B`) from URL query part, and return `A + B`
type get_query_params struct {
	A plugins.QueryParam[any]
	B plugins.QueryParam[any]
}
func (v *get_query_params) Handle(c *nji.Context) {
	c.Writer.WriteString(v.A.Value + v.B.Value)
}
```

## parameter verification with `schema`

```go
type get_query_params struct {
	A plugins.QueryParam[schema.NotNull] // check if `A` is null

	B plugins.QueryParam[struct { // multi-verificator
		schema.NotNull 
		schema.IsPhoneNumber `could pass gotag to schema` // check if `B` is valid as phonenumber
	}]
	
	// or you could write multi-verificator in one-line by this way:
	C plugins.QueryParam[func(schema.NotNull, schema.IsPhoneNumber)]
}
```

## compute route automatically with `route`

```go
type BaseRoute struct {
	nji.Route[route.ANY, route.ROOT] `a_prefix` // if no gotag, path will use `/BaseRoute/` as component
}

type get_query_params struct {
	nji.Route[route.GET, BaseRoute] // output URL = `/a_prefix/get_query_params`

	A plugins.QueryParam[any]
	// ...
}
```









