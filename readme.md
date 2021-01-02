# nji

名称取自 inject 前三个字母

## 示例 PathParam

```go
// ./plugins/PathParam_test.go

type a struct {
	A plugins.PathParam
}

func (view *a) Handle(c *nji.Context) {
	c.ResponseWriter.WriteHeader(200)
	_, _ = c.ResponseWriter.Write([]byte(view.A.Value))
}

func main() {
	app := nji.NewServer()
	app.GET("/param/:A", nji.Inject(&a{}))
	app.Run(8003)
}

```

## 示例 JSON

```go
// ./plugins/dyn.JSON_test.go
type json_t struct {
	Body struct{
		plugins.DynJSON
		A string
		B string
	}
}

func (v *json_t) Handle(c *nji.Context) {
	c.ResponseWriter.WriteHeader(200)
	_, _ = c.ResponseWriter.Write([]byte(v.Body.A+v.Body.B))
}

func TestContextJSON(t *testing.T) {
	app := nji.NewServer()
	app.POST("/api/", nji.Inject(&json_t{}))
	reader := strings.NewReader(`{"A":"Hello ", "B": "World!"}`)
	r, err := http.NewRequest(http.MethodPost, "/api/", reader)
	if err != nil {
		t.Error(err.Error())
		return
	}
	r.Header.Add("Content-Type","application/json")
	w := httptest.NewRecorder()
	app.ServeHTTP(w, r)
	assert.Equal(t, "Hello World!", string(w.Body.Bytes()))
}

```

ab 测试：

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