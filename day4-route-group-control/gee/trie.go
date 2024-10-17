package gee

import "strings"

// 用trie树于实现动态路由
// 若采用day2的map映射，对于某一类url，无法简单通过map实现这样一类的映射
// 动态路由，与静态路由相对，其核心在于运行时生成URL。动态路由允许根据用户提供的或生成的参数自动生成页面URL

type node struct {
	pattern  string  // 待匹配的完整路由，例如 /p/:lang
	part     string  // 路由中的一部分，如 :lang
	children []*node // 当前节点的子节点
	isWild   bool    // 是否为模糊匹配, 即含有 : 或 * 时为true
	// :为参数匹配 例如 /p/:lang/doc 为参数匹配，可以匹配 /p/c/doc 和 /p/go/doc
	// *为通配 例如 /static/*filepath，可以匹配/static/fav.ico，也可以匹配/static/js/jQuery.js
	// 通配常用于静态服务器，能够递归地匹配子路径。
}

// 返回第一次匹配成功的结点，用于插入
func (n *node) FirstMatch(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 返回所有匹配成功的结点，用于查找
func (n *node) AllMatch(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			// 将 child 元素追加到 nodes 切片的末尾，返回新的切片并赋值给 nodes
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// trie树的插入
func (n *node) insert(pattern string, parts []string, height int) {
	//
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height]
	child := n.FirstMatch(part)
	// 子节点中没有能够匹配的，新建一个结点
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		// 将新建节点加入
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

// trie树的查找
func (n *node) search(parts []string, height int) *node {
	// * 代表通配，匹配成功时也需要退出
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.AllMatch(part)

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}
