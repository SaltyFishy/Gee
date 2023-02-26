package gee

import (
	"net/http"
)

type router struct {
	handlers map[string]HandlerFunc
}

// 创建路由管理
func NewRouter() *router {
	return &router{handlers: make(map[string]HandlerFunc)}
}

// 添加路由-router
func (r *router) AddRouter (method string, pattern string, handlerFunc HandlerFunc) {
	key := method + "-" + pattern
	r.handlers[key] = handlerFunc
}

// 根据路由执行对应的操作
func (r *router) handle (c *Context) {
	key := c.Method + "-" + c.Path
	if handlerFunc, ok := r.handlers[key]; ok != false {
		handlerFunc(c)
	} else {
		c.String(http.StatusNotFound, "404 Page Not Found: %s\n", c.Path)
	}
}