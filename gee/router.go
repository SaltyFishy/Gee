package gee

import (
	"net/http"
	"strings"
)

type router struct {
	roots map[string]*node
	handlers map[string]HandlerFunc
}

// 创建路由管理
func NewRouter() *router {
	return &router{
		roots: make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// 检验路由合法性(此处检验一个*)
func ParsePattern(pattern string) []string {
	parts := strings.Split(pattern, "/")
	text := make([]string, 0)
	for _, part := range parts {
		if part != "" {
			text = append(text, part)
			if part[0] == '*' {
				break
			}
		}
	}
	return text
}

// 添加路由-router
func (r *router) AddRouter(method string, pattern string, handlerFunc HandlerFunc) {
	parts := ParsePattern(pattern)
	key := method + "-" + pattern
	_, ok := r.roots[method]
	if ok == false {
		r.roots[method] = &node{}
	}
	r.handlers[key] = handlerFunc
	r.roots[method].Insert(pattern, parts, 0)
}

// 获取路由
func (r *router) GetRoute(method string, pattern string) (*node, map[string]string) {
	searchParts := ParsePattern(pattern)
	params := make(map[string]string)
	_, ok := r.roots[method]
	if ok == false {
		return nil, nil
	}
	n := r.roots[method].Search(searchParts, 0)
	if n != nil {
		parts := ParsePattern(n.pattern)
		for index, part := range parts {
			if strings.HasPrefix(part,":") {
				params[part[1:]] = searchParts[index]
			}
			if part == "*" && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

// 根据路由执行对应的操作
func (r *router) Handle(c *Context) {
	n, params := r.GetRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + c.Path
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 Page Not Found: %s\n", c.Path)
	}
}