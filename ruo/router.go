package ruo

import (
	"log"
	"net/http"
	"strings"
)

type router struct {
	handlers map[string]HandlerFunc // 存储请求方式的HandlerFunc
	roots    map[string]*node       // 存储请求方式的前缀树根节点
}

// roots key eg, roots['GET'] roots['POST']
// handlers key eg, handlers['GET-/p/:lang/doc'], handlers['POST-/p/book']

func newRouter() *router {
	return &router{
		handlers: make(map[string]HandlerFunc),
		roots:    make(map[string]*node),
	}
}

// 解析路由规则
// /hello/ruo => [hello ruo]
func parsePattern(pattern string) (parts []string) {
	vs := strings.Split(pattern, "/") // 切割字符串

	parts = make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return
}

// 添加路由
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	//key := method + "-" + pattern
	//r.handlers[key] = handler
	parts := parsePattern(pattern)

	key := method + "-" + pattern

	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

// 获取路由
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	log.Println("开始匹配路由", path)
	searchParts := parsePattern(path)
	log.Println("解析路由", searchParts)
	params := make(map[string]string)
	// 获取对应请求方式的前缀树
	root, ok := r.roots[method]

	if !ok {
		return nil, nil
	}
	// 查找对应规则
	n := root.search(searchParts, 0)
	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

// 实现 handle
// 将从路由匹配得到的 Handler 添加到 c.handlers列表中，执行c.Next()。
func (r *router) handle(c *Context) {
	//key := c.Method + "-" + c.Path
	//if handler, ok := r.handlers[key]; ok {
	//	handler(c)
	//}
	n, params := r.getRoute(c.Method, c.Path)
	log.Printf("[%s] 处理路由 %s ", c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	c.Next()
}
