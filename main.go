package main

import (
	"github.com/xiaohan1995/Gin-blog/models"
	"github.com/xiaohan1995/Gin-blog/routers"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库
	models.InitDB()
	//初始化gin
	r := gin.Default()
	//设置静态资源和模版路径
	r.Static("/statics", "./statics")
	r.LoadHTMLGlob("templates/*")

	//初始化路由
	routers.InitRouter(r)
	models.Log.Info("服务已经启动：8081")
	r.Run(":8081")
}
