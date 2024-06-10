package main

import (
	"context"
	"log"

	"github.com/chasemao/blogbackend/routers"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
)

func main() {
	// initES()

	router := gin.Default()
	// 注册各个模块的路由
	routers.RegisterUserRoutes(router)
	// 启动Gin服务器
	router.Run("localhost:6666")

}

func initES() {
	// 设置 Elasticsearch 客户端
	client, err := elastic.NewClient(elastic.SetURL("http://localhost:9200"), elastic.SetSniff(false))
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	// 检查 Elasticsearch 节点信息
	info, code, err := client.Ping("http://localhost:9200").Do(context.Background())
	if err != nil {
		log.Fatalf("Error pinging Elasticsearch: %s", err)
	}
	log.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)
}
