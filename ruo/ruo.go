package ruo

import (
	"fmt"
	"log"
	"net/http"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request)

// Engine 实现ServeHTTP接口
type Engine struct {
	router map[string]HandlerFunc // 路由
}

// ruo.Engine 的构造器
func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

// 添加路由
func (engine *Engine) AddRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	engine.router[key] = handler
}

// 处理get请求
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.AddRoute("GET", pattern, handler)
}

// 处理post请求
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.AddRoute("POST", pattern, handler)
}

// 启动指定端口的http服务
func (engine *Engine) Run(addr string) (err error) {
	log.Println("server start at ", addr)
	return http.ListenAndServe(addr, engine)
}

// 实现 serveHttp接口
// 解析请求的路径，查找路由映射表，如果查到，就执行注册的处理方法。如果查不到，就返回
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	if handler, ok := engine.router[key]; ok {
		log.Printf("Handler router %s method %s \n", req.URL.Path, req.Method)
		handler(w, req)
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}
