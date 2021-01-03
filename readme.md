# nji

名称取自 inject 前三个字母

# 示例
每个示例是一个完整的 HTTP Handler 
## PathParam

返回 URL  `/api/***/:A` 中的 参数 `:A`

```go
// ./plugins/PathParam_test.go
type a struct {
	A plugins.PathParam
}

func (view *a) Handle(c *nji.Context) {
	c.Resp.String(200,view.A.Value)
}
```

## QueryParam
返回 URL  `/api/***/?A=Hello&B=World!` 中的 参数 `?A & B`
```go
// ./plugins/QueryParam_test.go
type c struct {
	A plugins.QueryParam
	B,C,D,E,F plugins.QueryParamOptional
}

func (v *c) Handle(c *nji.Context) {
	c.Resp.String(200,v.A.Value+v.B.Value)
}
```

## JSON

```go
// ./plugins/dyn.JSON_test.go
type j struct {
	Body struct{
		plugins.DynJSON
		A string
		B string
	}
}

func (v *j) Handle(c *nji.Context) {
	c.Resp.String(200,view.Body.A + view.Body.B)
}
```

# 特性说明

nji 通过使用依赖注入来节省业务代码的反复书写

它提供接口 `Plugin`  来达成这个目的， 一共有两种使用 Plugin 的方式

## 1. 简单插件

具名结构体作为 `Plugin`

```go
type view struct {
	A plugins.PathParam
}
```

## 2. 匿名结构体插件

整个匿名结构体作为 `Plugin` 

（外层结构体通过嵌套一个 Plugin 实体来实现 Plugin 接口）

```go
type json struct {
	Body struct{
		plugins.DynJSON
		A string
		B string
	}
}
```

## 使用情景

1. 从 request HEAD 获取信息， 比如 Auth
2. 从 BODY 解析信息并获取信息，比如 JSON
3. 验证信息， 是否存在、是否合格， 比如 验证验证码
4. 从 URL 获取信息， 比如 QueryParam


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
