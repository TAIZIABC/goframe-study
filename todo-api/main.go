package main

import (
	_ "todo-api/internal/packed"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"

	"github.com/gogf/gf/v2/os/gctx"

	"todo-api/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
