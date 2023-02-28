package gee

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func newTestRouter() *router {
	r := NewRouter()
	r.AddRouter("GET", "/", nil)
	r.AddRouter("GET", "/hello/:name", nil)
	r.AddRouter("GET", "/hello/18/19", nil)
	r.AddRouter("GET", "/hi/:name", nil)
	r.AddRouter("GET", "/assets/*filepath", nil)
	return r
}

func TestParsePattern(t *testing.T) {
	ok := reflect.DeepEqual(ParsePattern("/p/:name"), []string{"p", ":name"})
	ok = ok && reflect.DeepEqual(ParsePattern("/p/*"), []string{"p", "*"})
	ok = ok && reflect.DeepEqual(ParsePattern("/p/*name/*"), []string{"p", "*name"})
	if !ok {
		t.Fatal("test parsePattern failed")
	}
}

func TestGetRoute(t *testing.T) {
	r := newTestRouter()
	cases := []struct{
		method string
		pattern string
	}{
		{"GET", "/hello/18"},
		{"GET", "/hello/19"},
	}

	for _, c := range cases {
		n, _ := r.GetRoute(c.method, c.pattern)
		if n == nil {
			t.Fatalf("router may insert error\n")
		}

		parts := ParsePattern(n.pattern)

		if n.pattern != c.pattern && !strings.HasPrefix(parts[1], ":") {
			t.Fatalf("%s should match %s\n", n.pattern, c.pattern)
		}

		fmt.Printf("matched path: %s, params[%s]: %s\n", n.pattern, parts[1], ParsePattern(c.pattern)[1])
	}



}