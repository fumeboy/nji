# Plugin 实现建议


不是所有场景都适合用 inject 调用 plugin， 还需要考虑手动执行 plugin 功能的情况

建议为 plugin 加上额外的 Exec 函数用于手动执行 plugin 功能

比如

```go
func (pl *DynJSON) Exec(i io.Reader, obj interface{}) error {
	return jsoniter.NewDecoder(i).Decode(&obj)
}

func (pl *PostParam) Exec(c *nji.Context, name string) error {
	var ok bool
	pl.Value,ok = c.PostParam(name)
	if ok{
		return nil
	}else{
		return postParamFail
	}
}

func (pl *QueryParam) Exec(c *nji.Context, name string) error {
	var ok bool
	pl.Value,ok = c.QueryParam(name)
	if ok{
		return nil
	}else{
		return queryParamFail
	}
}

func (pl *PathParam) Exec(c *Context, name string) {
	pl.Value,_ = c.PathParams.Get(name)
}
```