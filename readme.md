# nji

名称取自 inject 前三个字母

# 示例

```go
package main

import (
	"net/http"
	"nji"
	"nji/plugins"
)

type a struct {
	A plugins.PathParam
}

func (view *a) Handle(c *nji.Context) {
	c.ResponseWriter.WriteHeader(200)
	_, _ = c.ResponseWriter.Write([]byte(view.A.Value))
}

func main() {
	app := nji.Config{
		UnescapePathValues: true,
	}.New()
	app.GET("/param/:A", nji.Inject(&a{}))
	app.Run(8003)
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