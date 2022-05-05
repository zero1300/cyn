package cyn

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type H map[string]any

type Context struct {
	Writer     http.ResponseWriter
	Req        *http.Request
	Method     string
	Path       string
	StatusCode int
}

func NewContext(writer http.ResponseWriter, req *http.Request) *Context {
	return &Context{Writer: writer, Req: req, Method: req.Method, Path: req.URL.Path}
}

/* --- 获得参数的方法 --- */

// PostForm 获得Post请求的参数
func (c Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// PostBody 获得Post请求的Body
func (c Context) PostBody() string {
	body, err := ioutil.ReadAll(c.Req.Body)
	if err != nil {
		log.Fatal("Read Post body error")
	}
	fmt.Println(string(body))
	return ""
}

// Query 获得 URL 中的参数
func (c Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

/* 回写响应消息的方法 */

func (c Context) SetStatus(status int) {
	c.StatusCode = status
	c.Writer.WriteHeader(status)
}

func (c Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// String 调用我们的 SetHeader和Status 方法，构造string类型响应的状态码和消息头，然后将字符串转换成byte写入到响应头
func (c Context) String(code int, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.SetStatus(code)
	var str = ""
	for _, value := range values {
		str += value.(string)
	}
	_, err := c.Writer.Write([]byte(str))
	if err != nil {
		return
	}
}

// JSON 调用我们的 SetHeader和Status 方法，构造JSON类型响应的状态码和消息头，根据我们传入的对象来构造json数据写入
func (c Context) JSON(code int, obj any) {
	c.SetHeader("Content-Type", "application/json")
	c.SetStatus(code)
	encoder := json.NewEncoder(c.Writer)
	err := encoder.Encode(obj)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
}

// Data 同上 ，但是直接写入字节数组，不再构建消息头
func (c Context) Data(code int, data []byte) {
	c.SetStatus(code)
	_, err := c.Writer.Write(data)
	if err != nil {
		return
	}
}

// HTML 模版渲染 同上，消息体传入的是html文件
func (c Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.SetStatus(code)
	_, err := c.Writer.Write([]byte(html))
	if err != nil {
		return
	}
}
