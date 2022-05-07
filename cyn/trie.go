package cyn

import "strings"

type node struct {
	path     string
	part     string
	children []*node
	isBlurry bool //  是否模糊匹配, 当 path 中有 : 或 * 时为true
}

// matchChild 根据传入的part, 返回第一个匹配成功的节点,用于插入
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || n.isBlurry {
			return child
		}
	}
	return nil
}

//  所有匹配成功的节点，用于查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isBlurry {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

func (n *node) insert(path string, parts []string, height int) {
	// 遍历到底部了, 开始插入
	if len(parts) == height {
		n.path = path
		return
	}

	part := parts[height]
	child := n.matchChild(part)

	if child == nil {
		child = &node{part: part,
			isBlurry: part[0] == ':' || part[0] == '*',
		}
		n.children = append(n.children, child)
	}
	child.insert(path, parts, height+1)
}

// search 在 trie 树中查找符合条件的节点; 条件: parts 列表的最后一个, 或当前节点包含 `*`
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.path == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}
