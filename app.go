package nji

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

// body 大小限制
const MaxMultipartMemory = 1<<10<<2 // 4k

type CORS struct {
	AllowOrigins     string
	ExposeHeaders    string
	AllowMethods     string
	AllowHeaders     string
	AllowCredentials bool
}

// 配置
type Config struct {
	UseRawPath         bool  // 使用url.RawPath查找参数
	UnescapePathValues bool  // 反转义路由参数
	MaxMultipartMemory int64 // 允许的请求Body大小

	RootPath string // 根路径
	CORS *CORS
}

type Engine struct {
	Router
	Config      Config
	contextPool sync.Pool   // context池
	trees       methodTrees // 路由树
}

func NewServer() *Engine {
	return Config{
		UnescapePathValues: true,
		MaxMultipartMemory: MaxMultipartMemory,
	}.New()
}

// 创建一个新引擎
func (config Config) New() *Engine {
	if config.MaxMultipartMemory == 0 {
		config.MaxMultipartMemory = MaxMultipartMemory
	}
	if config.CORS != nil {
		if config.CORS.AllowMethods == "" {
			config.CORS.AllowMethods = "GET,POST,PUT,DELETE,OPTIONS,PATCH"
		}
		if config.CORS.AllowHeaders == "" {
			config.CORS.AllowHeaders = "*"
		}
		if config.CORS.ExposeHeaders == "" {
			config.CORS.ExposeHeaders = "*"
		}
	}
	if config.RootPath == ""{
		config.RootPath = "/"
	}
	// 初始化一个引擎
	engine := &Engine{
		// 初始化根路由组
		Router: Router{
			Handlers: nil,
			basePath: config.RootPath,
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

// 实现http.Handler接口
func (engine *Engine) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := engine.contextPool.Get().(*Context)
	defer func() {
		if err := recover(); err != nil {
			ctx.Resp.String(500,"Internal Server err")
		}
	}()
	ctx.reset(req, resp)

	if engine.Config.CORS != nil {
		engine.setCORS(ctx.Request, ctx.Resp.Writer)
		if ctx.Request.Method == "OPTIONS" {
			ctx.Resp.Writer.WriteHeader(http.StatusNoContent)
			engine.contextPool.Put(ctx)
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
			goto END
		}
		break
	}

	// 404
	ctx.Resp.String(404,"404 not found")
END:
	engine.contextPool.Put(ctx)
}

func (engine *Engine) Run(port int) error {
	return http.ListenAndServe(fmt.Sprintf(":%d", port), engine)
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
	if engine.Config.CORS.AllowOrigins == "" {
		resp.Header().Set("Access-Control-Allow-Origin", req.Header.Get("Origin"))
	} else {
		resp.Header().Set("Access-Control-Allow-Origin", engine.Config.CORS.AllowOrigins)
	}
	resp.Header().Set("Access-Control-Allow-Methods", engine.Config.CORS.AllowMethods)
	resp.Header().Set("Access-Control-Allow-Headers", engine.Config.CORS.AllowHeaders)
	resp.Header().Set("Access-Control-Expose-Headers", engine.Config.CORS.ExposeHeaders)
	resp.Header().Set("Access-Control-Allow-Credentials", strconv.FormatBool(engine.Config.CORS.AllowCredentials))
}
