package gee

import (
	"log"
	"net/http"
)

type HandlerFunc func (c *Context)

// 框架引擎
type (
	RouterGroup struct {
		prefix      string
		middlewares []HandlerFunc // 支持中间件
		engine      *Engine       // 所有分组共享的引擎实例
	}

	Engine struct {
		*RouterGroup 		  // 使Engine能够含有RouterGroup的属性
		router *router
		groups []*RouterGroup // 存储所有的路由分组
	}
)

// 构造路由引擎
func New() *Engine{
	engine := &Engine{router: NewRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// 添加路由-gee
func (e *Engine) AddRouter(method string, pattern string, handlerFunc HandlerFunc) {
	e.router.AddRouter(method, pattern, handlerFunc)
}

// 设置GET请求
func (e *Engine) GET(pattern string, handlerFunc HandlerFunc) {
	e.AddRouter("GET", pattern, handlerFunc)
}

// 设置POST请求
func (e *Engine) POST(pattern string, handlerFunc HandlerFunc) {
	e.AddRouter("POST", pattern, handlerFunc)
}

// 启动引擎
func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

// 实现ServeHTTP以实现接管所有的HTTP请求，并构造上下文
func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := NewContext(w, req)
	e.router.Handle(c)
}

// 创建新的路由组
// 所有的路由组共享同一个引擎
func (rg *RouterGroup) Group(prefix string) *RouterGroup {
	engine := rg.engine
	newGroup := &RouterGroup{
		prefix: rg.prefix + prefix,
		middlewares: nil,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// 针对路由组的AddRouter
func (rg *RouterGroup) AddRouter(method string, pattern string, handlerFunc HandlerFunc) {
	route := rg.prefix + pattern
	log.Printf("method %v %v", method, route)
	rg.engine.router.AddRouter(method, route, handlerFunc)
}

// 针对路由组的GET方法
func (rg *RouterGroup) GET(pattern string, handlerFunc HandlerFunc) {
	rg.AddRouter("GET", pattern, handlerFunc)
}

// 针对路由组的POST方法
func (rg *RouterGroup) POST(pattern string, handlerFunc HandlerFunc) {
	rg.AddRouter("POST", pattern, handlerFunc)
}