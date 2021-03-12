# nji

名称取自 inject 前三个字母

# 示例
每个示例是一个完整的 HTTP Handler 
## inject:PathParam

返回 URL  `/api/***/:A` 中的 参数 `:A`

```go
// ./inject/PathParam_test.go
type a struct {
	A inject.PathParam
}

func (view *a) Handle(c *nji.Context) {
	c.Resp.String(200,view.A.Value)
}
```

## inject:QueryParam
返回 URL  `/api/***/?A=Hello&B=World!` 中的 参数 `?A & B`
```go
// ./inject/QueryParam_test.go
type c struct {
    A inject.QueryParam 
    B,C,D,E,F struct{
        inject.QueryParam
        schema.IsNull
    }
}

func (v *c) Handle(c *nji.Context) {
	c.Resp.String(200,v.A.Value+v.B.Value)
}
```

## schema:IsPhoneNumber

校验是否是合法的手机号
```go
type a struct {
	A struct {
		inject.PathParam
		schema.IsPhoneNumber
	}
}

func (view *a) Handle(c *nji.Context) {
	c.Resp.String(200, view.A.Value)
}
```

# 特性说明

## inject

nji 通过使用依赖注入来节省业务代码的反复书写

它提供接口 `Plugin`  来达成这个目的

```go
type view struct {
	A inject.PathParam
}
func (view *a) Handle(c *nji.Context) {
	c.Resp.String(200,view.A.Value)
}
```
其中 PathParam 就是一个从URL获取信息的 Plugin

如果不使用依赖注入，那么大概是这样写业务代码：

```go
func Handle(c *nji.Context) {
    A,err := c.GetPathParam("A")
    if err != nil{
        ...
    }   
    c.Resp.String(200,A)
}
```

详见 `./inject` 文件夹

## schema

同时提供简单的 schema 机制，用于校验数据合法性 

详见 `./schema` 文件夹


# 性能测试：

`ab -n 10000 -c 100 http://127.0.0.1:8003/param/123`

```
Concurrency Level:      100
Time taken for tests:   0.856 seconds
Complete requests:      10000
Failed requests:        0
Total transferred:      1230000 bytes
HTML transferred:       70000 bytes
Requests per second:    11675.76 [#/sec] (mean)
Time per request:       8.565 [ms] (mean)
Time per request:       0.086 [ms] (mean, across all concurrent requests)
Transfer rate:          1402.46 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    4   1.1      4      10
Processing:     2    4   1.1      4      10
Waiting:        1    4   1.1      4      10
Total:          5    8   1.8      8      16
```
