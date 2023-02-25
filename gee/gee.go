package gee

import (
	"fmt"
	"log"
	"net/http"
)

type HandlerFunc func (http.ResponseWriter, *http.Request)

type Engine struct {
	router map[string]HandlerFunc
}

func New() *Engine{
	return &Engine{router: make(map[string]HandlerFunc)}
}

func (e *Engine) AddRouter (method string, pattern string, handlerFunc HandlerFunc) {
	key := method + "-" + pattern
	e.router[key] = handlerFunc
}

func (e *Engine) GET (pattern string, handlerFunc HandlerFunc) {
	e.AddRouter("GET", pattern, handlerFunc)
}

func (e *Engine) POST (pattern string, handlerFunc HandlerFunc) {
	e.AddRouter("POST", pattern, handlerFunc)
}

func (e *Engine) Run (addr string) error {
	return http.ListenAndServe(addr, e)
}

func (e *Engine) ServeHTTP (w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	if handlerFunc, ok := e.router[key]; ok != false {
		handlerFunc(w, req)
	} else {
		w.WriteHeader(http.StatusNotFound)
		_, err := fmt.Fprintf(w, "404 not found: %s\n", req.URL)
		if err != nil {
			log.Printf("写入失败\n")
			return
		}
	}
}
