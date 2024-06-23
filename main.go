package main

import (
	"os"

	"github.com/chasemao/blogbackend/handlers"
	"github.com/chasemao/blogbackend/routers"
	"github.com/gin-gonic/gin"
)

func main() {
	args := os.Args[1:]
	if len(args) > 0 {
		handlers.SetArticleDic(args[0])
	}

	router := gin.Default()
	// 注册各个模块的路由
	routers.RegisterUserRoutes(router)
	// 启动Gin服务器
	router.Run("localhost:6666")

}
