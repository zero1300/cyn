package cyn

import (
	"net/http"
	"strings"
)

type router struct {
	handlers map[string]HandleFunc
	// 路由前缀树的根节点
	roots map[string]*node
}

func newRouter() *router {
	return &router{make(map[string]HandleFunc), make(map[string]*node)}
}

// parsePath 根据分隔符 `/` 拆分path; /hello/:lang -> ["hello", ":lang"]
func parsePath(path string) (parts []string) {
	vs := strings.Split(path, "/")
	parts = make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

// addRoute 注册路由
func (r *router) addRoute(method string, path string, handler HandleFunc) {
	parts := parsePath(path)

	key := method + "-" + path

	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	// 建树
	r.roots[method].insert(path, parts, 0)
	r.handlers[key] = handler
}

//  getRoute 根据传入的方法和路径, 返回符合元素的节点(通常是path的最后一段), 和path参数
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePath(path)
	params := make(map[string]string)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	node := root.search(searchParts, 0)
	if node != nil {
		parts := parsePath(node.path)
		for i, part := range parts {
			// 解析参数
			if part[0] == ':' {
				params[part[1:]] = searchParts[i]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[i:], "/")
				break
			}
		}
		return node, params
	}
	return nil, nil
}

func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.path
		// invoke HandleFun
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND %s\n", c.Path)
	}
}
