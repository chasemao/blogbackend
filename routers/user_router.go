package routers

import (
	"github.com/chasemao/blogbackend/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.Engine) {
	userGroup := router.Group("/user")
	{
		userGroup.GET("/:id", handlers.GetUser)
		userGroup.POST("/", handlers.CreateUser)
	}
}
