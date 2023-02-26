package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	Writer http.ResponseWriter
	Req *http.Request
	Method string
	Path string
	StatusCode int
}

// 构造上下文
func NewContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req: req,
		Method: req.Method,
		Path: req.URL.Path,
	}
}

// 构造Post返回体
func (c *Context) PostForm (key string) string {
	return c.Req.FormValue(key)
}

// 指定键查询第一个value
func (c *Context) Query (key string) string {
	return c.Req.URL.Query().Get(key)
}

// 构造状态
func (c *Context) Status (code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// 设置表头
func (c *Context) SetHeader (key string, value string) {
	c.Writer.Header().Set(key, value)
}

// 构造String
func (c *Context) String (code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// 构造JSON
func (c *Context) JSON (code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
}

// 构造Data
func (c *Context) Data (code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

// 构造HTML
func (c *Context) HTML (code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}