package routers

import (
	"github.com/chasemao/blogbackend/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.Engine) {
	articleLogic := handlers.NewArticleLogic()
	userGroup := router.Group("/api/v1/article")
	userGroup.GET("/list", articleLogic.ListArticles)
	userGroup.GET("/:id", articleLogic.GetArticle)
}
