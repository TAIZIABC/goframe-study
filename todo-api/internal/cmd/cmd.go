// internal/cmd/cmd.go
package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"

	"todo-api/internal/controller"
)

var Main = gcmd.Command{
	Name: "main",
	Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
		s := g.Server()
		s.Use(ghttp.MiddlewareHandlerResponse)
		s.Group("/api/v1", func(group *ghttp.RouterGroup) {
			group.Bind(controller.Todo)
		})
		s.Run()
		return nil
	},
}
