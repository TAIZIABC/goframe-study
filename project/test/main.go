package main

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type Response struct {
	Message string      `json:"message" dc:"消息提示"`
	Data    interface{} `json:"data" dc:"执行结果"`
}

type HelloReq struct {
	g.Meta `path:"/say" method:"get"`
	Name   string `v:"required" json:"name" dc:"姓名"`
	Age    int    `v:"required" json:"age" dc:"年龄"`
}

type HelloRes struct {
	Content string `json:"content" dc:"内容"`
}

type Hello struct{}

func (Hello) Say(ctx context.Context, req *HelloReq) (res *HelloRes, err error) {
	// r := g.RequestFromCtx(ctx)
	// r.Response.Writef("hello %s, You age is %d", req.Name, req.Age)
	res = &HelloRes{
		Content: fmt.Sprintf("hello %s, You age is %d", req.Name, req.Age),
	}
	return
}

func ErrorHandler(r *ghttp.Request) {
	r.Middleware.Next()
	if err := r.GetError(); err != nil {
		r.Response.Write("error at: ", err.Error())
		return
	}
}

func ResponseHandler(r *ghttp.Request) {
	r.Middleware.Next()
	var (
		msg string
		res = r.GetHandlerResponse()
		err = r.GetError()
	)

	if err != nil {
		msg = err.Error()
	} else {
		msg = "success"
	}
	r.Response.WriteJson(&Response{
		Message: msg,
		Data:    res,
	})
}

func main() {
	fmt.Println("Hello, World!", gf.VERSION)

	s := g.Server()
	// s.BindHandler("/", func(r *ghttp.Request) {
	// 	var req HelloReq
	// 	if err := r.Parse(&req); err != nil {
	// 		r.Response.Write(err.Error())
	// 		return
	// 	}
	// 	if req.Name == "" {
	// 		r.Response.Write("name is empty")
	// 		return
	// 	}

	// 	if req.Age <= 0 {
	// 		r.Response.Write("age is empty")
	// 		return
	// 	}
	// 	r.Response.Write("Hello, World!    ", req.Name, req.Age)
	// })

	s.Group("/", func(group *ghttp.RouterGroup) {
		// group.Middleware(ErrorHandler)
		group.Middleware(ResponseHandler)
		group.Bind(
			new(Hello),
		)
	})

	s.SetOpenApiPath("/api.json")
	s.SetSwaggerPath("/swagger")
	s.SetPort(8199)
	s.Run()
}
