package main

import (
	"github.com/gin-gonic/gin"
	"k8s-platform/config"
	"k8s-platform/controller"
	"k8s-platform/service"
)

func main() {
	//初始化gin
	r := gin.Default()
	//初始化k8s client
	service.K8s.Init()
	//初始化路由规则
	controller.Router.InitAPiRouter(r)
	//启动gin
	r.Run(config.ListenAddr)
}
