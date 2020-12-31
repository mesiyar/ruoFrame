package ruo

import (
	"log"
	"net/http"
	"strings"
)

type HandlerFunc func(c *Context)

type RouterGroup struct {
	prefix      string // 路由组前缀
	middleWares []HandlerFunc
	parent      *RouterGroup
	engine      *Engine
}

// Engine 实现ServeHTTP接口
type Engine struct {
	*RouterGroup
	router *router // 路由
	groups []*RouterGroup
}

// ruo.Engine 的构造器
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// 路由分组
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		engine: engine,
		parent: group,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// 添加路由
func (group *RouterGroup) AddRoute(method string, pattern string, handler HandlerFunc) {
	pattern = group.prefix + pattern
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

// 处理get请求
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.AddRoute("GET", pattern, handler)
}

// 处理post请求
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.AddRoute("POST", pattern, handler)
}

// 启动指定端口的http服务
func (engine *Engine) Run(addr string) (err error) {
	log.Println("server start at ", addr)
	return http.ListenAndServe(addr, engine)
}

// 调用中间件
func (group *RouterGroup) Use(middleWares ...HandlerFunc) {
	group.middleWares = append(group.middleWares, middleWares...)
}

// 实现 serveHttp接口
// 解析请求的路径，查找路由映射表，如果查到，就执行注册的处理方法。如果查不到，就返回
// 当我们接收到一个具体请求时，要判断该请求适用于哪些中间件，在这里我们简单通过 URL 的前缀来判断。得到中间件列表后，赋值给 c.handlers。
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//key := req.Method + "-" + req.URL.Path
	//if handler, ok := engine.router[key]; ok {
	//	log.Printf("Handler router %s method %s \n", req.URL.Path, req.Method)
	//	handler(w, req)
	//} else {
	//	w.WriteHeader(http.StatusNotFound)
	//	fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	//}
	var middleWares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middleWares = append(middleWares, group.middleWares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middleWares
	engine.router.handle(c)
}
