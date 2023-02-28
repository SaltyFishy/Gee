package gee

import (
	"net/http"
)

type HandlerFunc func (c *Context)

type Engine struct {
	router *router
}

// 构造路由引擎
func New() *Engine{
	return &Engine{router: NewRouter()}
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
