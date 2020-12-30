package ruo

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

// 上下文结构体
type Context struct {
	// 原始对象
	Writer http.ResponseWriter
	Req    *http.Request

	// 请求信息
	Path   string
	Method string

	// 响应信息
	StatusCode int
}

// 初始化上下文
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
	}
}

// 获取 post 数据
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// 获取 querystring 的参数
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// 设置返回状态
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(c.StatusCode)
}

// 设置header
func (c *Context) SetHeader(key, value string) {
	c.Writer.Header().Set(key, value)
}

// 返回格式 string
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// 返回格式 json
func (c *Context) Json(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}


func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	_, _ = c.Writer.Write(data)
}

// 返回格式 html
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	_, _ = c.Writer.Write([]byte(html))
}