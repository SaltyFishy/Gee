package gee

import "strings"

type node struct {
	pattern  string // 待匹配路由，例如 /p/:lang
	part     string // 路由中的一部分，例如 :lang
	children []*node // 子节点，例如 [doc, tutorial, intro]
	isWild   bool // 是否精确匹配，part 含有 : 或 * 时为true
}

// 查询trie当前层匹配的结点
func (n *node) Match (part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 匹配当前节点所有孩子
func (n *node) MatchChildren (part string) []*node{
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
func (n *node) Insert (pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	part := parts[height]
	child := n.Match(part)
	if child == nil {
		child = &node{pattern:"", part: part, children: make([]*node, 0), isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.Insert(pattern, parts, height + 1)
}

// 查找路由表
func (n *node) Search (parts []string, height int) *node {
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