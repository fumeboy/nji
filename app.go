package nji

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

// 引擎配置
type Config struct {
	UseRawPath         bool  // 使用url.RawPath查找参数
	UnescapePathValues bool  // 反转义路由参数
	MaxMultipartMemory int64 // 允许的请求Body大小(默认1 << 20 = 1MB)

	RootPath string // 根路径

	CORS             bool // 是否启用自动CORS处理
	AllowOrigins     string
	ExposeHeaders    string
	AllowMethods     string
	AllowHeaders     string
	AllowCredentials bool
}

// 引擎
type Engine struct {
	Router                 // 路由器
	Config      Config      // 配置
	contextPool sync.Pool   // context池
	trees       methodTrees // 路由树
}

// 创建一个新引擎
func (config Config) New() *Engine {
	if config.MaxMultipartMemory == 0 {
		config.MaxMultipartMemory = MaxMultipartMemory
	}
	if config.CORS {
		if config.AllowMethods == "" {
			config.AllowMethods = "GET,POST,PUT,DELETE,OPTIONS,PATCH"
		}
		if config.AllowHeaders == "" {
			config.AllowHeaders = "*"
		}
		if config.ExposeHeaders == "" {
			config.ExposeHeaders = "*"
		}
	}
	// 初始化一个引擎
	engine := &Engine{
		// 初始化根路由组
		Router: Router{
			Handlers: nil,
			basePath: "/",
			root:     true, // 标记为根路由组
		},
		Config: config,
		trees: make(methodTrees, 0, 7),
	}
	// 将引擎对象传入路由组中，便于访问引擎对象
	engine.engine = engine
	// 设置ctx池
	engine.contextPool.New = func() interface{} {
		return &Context{engine: engine}
	}
	return engine
}

// 实现http.Handler接口，并且是连接调度的入口
func (engine *Engine) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	defer func() {
		err := recover()
		if err != nil {
			// 触发panic事件
		}
	}()
	// 从池中取出一个ctx
	ctx := engine.contextPool.Get().(*Context)
	// 重置取出的ctx
	ctx.reset(req, resp)
	// 处理请求
	engine.handleRequest(ctx)
	// 将ctx放回池中
	engine.contextPool.Put(ctx)
}

func (engine *Engine) Run(port int) error {
	return http.ListenAndServe(fmt.Sprintf(":%d", port), engine)
}

// 处理连接请求
func (engine *Engine) handleRequest(ctx *Context) {
	if engine.Config.CORS {
		engine.setCORS(ctx.Request, ctx.ResponseWriter)
		if ctx.Request.Method == "OPTIONS" {
			ctx.ResponseWriter.WriteHeader(http.StatusNoContent)
			return
		}
	}

	rPath := ctx.Request.URL.Path
	unescape := false
	if engine.Config.UseRawPath && len(ctx.Request.URL.RawPath) > 0 {
		rPath = ctx.Request.URL.RawPath
		unescape = engine.Config.UnescapePathValues
	}

	for k := range engine.trees {
		if engine.trees[k].method != ctx.Request.Method {
			continue
		}
		root := engine.trees[k].root
		value := root.getValue(rPath, ctx.PathParams, unescape)
		if value.handlers != nil {
			ctx.handlers = value.handlers
			ctx.PathParams = value.params
			ctx.fullPath = value.fullPath
			ctx.Next()
			return
		}
		break
	}

	// 404
	ctx.ResponseWriter.WriteHeader(404)
	_, _ = ctx.ResponseWriter.Write([]byte("404 not found"))
}

// 添加路由
func (engine *Engine) addRoute(method, path string, handlers HandlersChain) {
	if path[0] != '/' {
		panic("The path must begin with '/'")
	}
	if method == "" {
		panic("HTTP method can not be empty")
	}
	if len(handlers) == 0 {
		panic("[" + method + "]" + path + " must be at least one handler")
	}

	// 查找方法是否存在
	root := engine.trees.get(method)
	// 如果方法不存在
	if root == nil {
		// 创建一个新的根节点
		root = new(node)
		root.fullPath = "/"
		engine.trees = append(engine.trees, methodTree{
			method: method,
			root:   root,
		})
	}
	root.addRoute(path, handlers)
}

// 在resp中设置CORS相关的头信息
func (engine *Engine) setCORS(req *http.Request, resp http.ResponseWriter) {
	if engine.Config.AllowOrigins == "" {
		resp.Header().Set("Access-Control-Allow-Origin", req.Header.Get("Origin"))
	} else {
		resp.Header().Set("Access-Control-Allow-Origin", engine.Config.AllowOrigins)
	}
	resp.Header().Set("Access-Control-Allow-Methods", engine.Config.AllowMethods)
	resp.Header().Set("Access-Control-Allow-Headers", engine.Config.AllowHeaders)
	resp.Header().Set("Access-Control-Expose-Headers", engine.Config.ExposeHeaders)
	resp.Header().Set("Access-Control-Allow-Credentials", strconv.FormatBool(engine.Config.AllowCredentials))
}
