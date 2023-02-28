package gee

import (
	"fmt"
	"strings"
)

type node struct {
	pattern  string // 树上完整待匹配路由
	part     string // 路由的一部分，包括了模糊匹配，例:/*,/:lang,/article
	children []*node // 当前路由的子路由，切片形式，即[doc, travel]
	isWild   bool // 是否精确匹配，part 含有 : 或 * 时为true
}

// 辅助函数，用于插入失败时构造错误信息
func (n *node) String(pattern string, part string) string {
	return fmt.Sprintf("An error occur while insert, router: {%s}, part: {%s} is conflict with part: {%s}", pattern, part, n.part)
}

// 查询trie当前层匹配的结点——精确查找
func (n *node) MatchSpecificChild(part string) *node {
	for _, child := range n.children {
		if child.part == part {
			return child
		}
	}
	return nil
}

// 查询trie当前层匹配的结点——模糊查找
func (n *node) MatchWildChild() *node {
	for _, child := range n.children {
		if child.isWild == true {
			return child
		}
	}
	return nil
}

// 匹配当前节点所有孩子
func (n *node) MatchChildren(part string) []*node{
	children := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			children = append(children, child)
		}
	}
	return children
}

// 路由注册
// pattern:完整路由，parts:路由的拆分，height树高
// 只有在确实存在路由时n.pattern才有效，否则为""
func (n *node) Insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	part := parts[height]
	child := n.MatchWildChild()  // 进行模糊查找
	if child != nil && child.isWild == true && (strings.HasPrefix(part, "*") || strings.HasPrefix(part, ":")) { // 模糊查找存在模糊路由，并且当前路由也含有模糊信息
		errorText := child.String(pattern, part)
		panic(errorText)
	}
	child = n.MatchSpecificChild(part) // 精确查找
	if child == nil {
		child = &node{pattern:"", part: part, children: make([]*node, 0), isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.Insert(pattern, parts, height + 1)
}

// 查找路由表
func (n *node) Search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}
	part := parts[height]
	children := n.MatchChildren(part)
	for _, child := range children {
		res := child.Search(parts, height + 1)
		if res != nil {
			return res
		}
	}
	return nil
}

//
func (n *node) Travel(list *([]*node)) {
	if n.pattern != "" {
		*list = append(*list, n)
	}
	for _, child := range n.children {
		child.Travel(list)
	}

}