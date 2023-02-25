package main

import (
	"Gee/gee"
	"fmt"
	"net/http"
)

func main() {
	r := gee.New()
	r.GET("/", func (w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Request URL is: %s\n", req.URL.Path)
	})
	r.GET("/hello", func (w http.ResponseWriter, req *http.Request) {
		for k, v := range req.Header {
			fmt.Fprintf(w, "[key,value] is: %v %v\n", k, v)
		}
	})
	r.Run(":8888")
}
