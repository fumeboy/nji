package nji

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)


// 测试回应
func TestEcho(t *testing.T) {
	app := Config{
		UnescapePathValues: true,
		MaxMultipartMemory: 20 << 20,
	}.New()
	app.GET("/", func(ctx *Context) {
		ctx.ResponseWriter.WriteHeader(200)
		_, _ = ctx.ResponseWriter.Write([]byte("Hello !"))
	})
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	app.ServeHTTP(httptest.NewRecorder(), r)
}

// 测试 PathParams
func TestURLParams(t *testing.T) {
	app := Config{
		UnescapePathValues: true,
		MaxMultipartMemory: 20 << 20,
	}.New()
	app.GET("/:path/:file", func(ctx *Context) {
		t.Log(ctx.PathParams.Value("path"))
		t.Log(ctx.PathParams.Value("file"))
	})
	r, err := http.NewRequest("GET", "/haha/hehe||123", nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	app.ServeHTTP(httptest.NewRecorder(), r)
}

// 测试Context传值
func TestContext(t *testing.T) {
	app := Config{
		UnescapePathValues: true,
		MaxMultipartMemory: 20 << 20,
	}.New()
	app.GET("/context", func(ctx *Context) {
		// 在ctx中写入参数
		ctx.SetValue("test", "hehe")
		t.Log(1, ctx.Request.URL.Path, "写值")
	}, func(ctx *Context) {
		// 从ctx中读取参数
		t.Log(2, ctx.Request.URL.Path, "取值：", ctx.GetValue("test"))
	})
	r, err := http.NewRequest("GET", "/context", nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	app.ServeHTTP(httptest.NewRecorder(), r)
}

// 测试路由组
func TestGroup(t *testing.T) {
	app := Config{
		UnescapePathValues: true,
		MaxMultipartMemory: 20 << 20,
	}.New()
	group := app.Group("/group", func(ctx *Context) {
		ctx.SetValue("test", "haha")
		t.Log(1, ctx.Request.URL.Path, "写值")
	})
	group.GET("/object", func(ctx *Context) {
		ctx.ResponseWriter.WriteHeader(http.StatusNoContent)
		t.Log(2, ctx.Request.URL.Path, "取值：", ctx.GetValue("test"))
	}, func(ctx *Context) {
		t.Log(3, ctx.Request.URL.Path, "取值：", ctx.GetValue("test"))
	})
	r, err := http.NewRequest("GET", "/group/object", nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	app.ServeHTTP(httptest.NewRecorder(), r)
}

// 测试中止
func TestAbort(t *testing.T) {
	app := Config{
		UnescapePathValues: true,
		MaxMultipartMemory: 20 << 20,
	}.New()
	group := app.Group("/group")
	group.GET("/object", func(ctx *Context) {
		t.Log(1, ctx.Request.URL.Path)
		ctx.Abort()
		t.Log(2, ctx.IsAborted())
	}, func(ctx *Context) {
		t.Log(3, ctx.Request.URL.Path)
	})
	r, err := http.NewRequest("GET", "/group/object", nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	app.ServeHTTP(httptest.NewRecorder(), r)
}

// 测试Append
func TestAppend(t *testing.T) {
	app := Config{
		UnescapePathValues: true,
		MaxMultipartMemory: 20 << 20,
	}.New()
	app.Append(func(ctx *Context) {
		t.Log(1, "append handler 1")
	}, func(ctx *Context) {
		t.Log(2, "append handler 2")
	})
	app.GET("/test", func(ctx *Context) {
		t.Log(3, ctx.Request.URL.Path)
	})
	r, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	app.ServeHTTP(httptest.NewRecorder(), r)
}

// 测试QueryParams
func TestQueryParams(t *testing.T) {
	app := Config{
		UnescapePathValues: true,
		MaxMultipartMemory: 20 << 20,
	}.New()
	app.GET("/object", func(ctx *Context) {
		t.Log(ctx.QueryParams())
	})
	r, err := http.NewRequest("GET", "/object?a=1&b=2", nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	app.ServeHTTP(httptest.NewRecorder(), r)
}

// 测试PostParams
func TestPostParams(t *testing.T) {
	app := Config{
		UnescapePathValues: true,
		MaxMultipartMemory: 20 << 20,
	}.New()
	app.POST("/object", func(ctx *Context) {
		t.Log(ctx.PostParams())
	})

	v := url.Values{}
	v.Add("a", "1")
	v.Add("b", "2")
	r, err := http.NewRequest("POST", "/object", strings.NewReader(v.Encode()))
	if err != nil {
		t.Error(err.Error())
		return
	}
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	app.ServeHTTP(httptest.NewRecorder(), r)
}

// 测试FormParams
func TestFormParams(t *testing.T) {
	app := Config{
		UnescapePathValues: true,
		MaxMultipartMemory: 20 << 20,
	}.New()
	app.POST("/object", func(ctx *Context) {
		t.Log(ctx.FormParams())
	})

	v := url.Values{}
	v.Add("c", "3")
	v.Add("d", "4")
	r, err := http.NewRequest("POST", "/object?a=1&b=2", strings.NewReader(v.Encode()))
	if err != nil {
		t.Error(err.Error())
		return
	}
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	app.ServeHTTP(httptest.NewRecorder(), r)
}

// 测试UnmarshalJSON
func TestPostUnmarshalJSON(t *testing.T) {
	type Obj struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}
	app := Config{
		UnescapePathValues: true,
		MaxMultipartMemory: 20 << 20,
	}.New()
	app.POST("/", func(ctx *Context) {
		var obj Obj
		t.Log(ctx.UnmarshalJSON(&obj))
	})

	var o Obj
	o.ID = 123
	o.Name = "dxvgef"
	ob, err := json.Marshal(&o)
	if err != nil {
		t.Error(err.Error())
		return
	}
	r, err := http.NewRequest("POST", "/", bytes.NewBuffer(ob))
	if err != nil {
		t.Error(err.Error())
		return
	}
	r.Header.Set("Content-Type", "application/json")
	app.ServeHTTP(httptest.NewRecorder(), r)
}

// 测试404事件
func TestNotFoundEvent(t *testing.T) {
	app := Config{
		RootPath:           getRootPath(),
		UnescapePathValues: true,
	}.New()
	r, err := http.NewRequest("GET", "/404", nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	app.ServeHTTP(httptest.NewRecorder(), r)
}

// 测试405事件
func TestMethodNotAllowedEvent(t *testing.T) {
	app := Config{
		RootPath:           getRootPath(),
		UnescapePathValues: true,
	}.New()
	app.POST("/", func(ctx *Context)  {
	})
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	app.ServeHTTP(httptest.NewRecorder(), r)
}

// 测试panic事件
func TestPanicEvent(t *testing.T) {
	app := Config{
		RootPath:           getRootPath(),
		UnescapePathValues: true,
		MaxMultipartMemory: 2 << 20,
		Recover:            true,
	}.New()
	app.GET("/", func(ctx *Context) {
		panic("这是panic消息")
	})
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	app.ServeHTTP(httptest.NewRecorder(), r)
}
