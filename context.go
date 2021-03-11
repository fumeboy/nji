package nji

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// 上下文
type Context struct {
	PathParams PathParams
	handlers   HandlersChain
	Resp       Response
	fullPath   string
	engine     *Engine
	Request    *http.Request
	index      int8

	parsed bool // 是否已解析body
	isJSON bool
	JSON []byte

	err error
}

var emptyValues url.Values

// 重置Context
func (ctx *Context) reset(req *http.Request, resp http.ResponseWriter) {
	ctx.Request = req
	ctx.Resp.Writer = resp
	ctx.PathParams = ctx.PathParams[0:0]
	ctx.handlers = nil
	ctx.index = -1
	ctx.fullPath = ""
	ctx.parsed = false
	ctx.isJSON = false
	ctx.JSON = nil
	ctx.err = nil
}

func (ctx *Context) IsJSON() bool {
	if ctx.parsed{
		return ctx.isJSON
	}
	if ct := ctx.Request.Header.Get("Content-Type"); ct != "application/json" {
		ctx.parsed = true
		return false
	}
	ctx.parsed = true
	data, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil{
		return false
	}
	ctx.JSON = data
	ctx.isJSON = true
	return true
}

// 解析form数据
func (ctx *Context) ParseForm() error {
	if ctx.parsed {
		return nil
	}
	ct := ctx.Request.Header.Get("Content-Type")
	if strings.HasPrefix(ct, "multipart/form-data") {
		if err := ctx.Request.ParseMultipartForm(ctx.engine.Config.MaxMultipartMemory); err != nil {
			return err
		}
	} else {
		if err := ctx.Request.ParseForm(); err != nil {
			return err
		}
	}
	ctx.parsed = true
	return nil
}

// 继续执行下一个处理器
func (ctx *Context) Next() {
	ctx.index++
	for ctx.index < int8(len(ctx.handlers)) {
		// 执行处理器
		ctx.handlers[ctx.index](ctx)
		ctx.index++
	}
}

// 中止执行
func (ctx *Context) Abort() {
	ctx.index = abortIndex
}

func (ctx *Context) IsAborted() bool {
	return ctx.index >= abortIndex
}
