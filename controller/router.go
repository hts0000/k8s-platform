package controller

import (
	"github.com/gin-gonic/gin"
)

// 实例化router结构体,可使用改对象点出首字母大写的方式(挎包调用)
var Router router

// 创建router结构体
type router struct{}

// 初始化路由规则创建测试api接口
func (r *router) InitAPiRouter(router *gin.Engine) {
	router.
		//pod操作
		GET("/api/k8s/pods", Pod.GetPods).
		GET("/api/k8s/pod/detail", Pod.GetPodDetail).
		DELETE("/api/k8s/pod/del", Pod.DeletePod).
		PUT("/api/k8s/pod/update", Pod.UpdatePod).
		GET("/api/k8s/pod/container", Pod.GetPodContainer).
		GET("/api/k8s/pod/log", Pod.GetPodLog).
		GET("/api/k8s/pod/numns", Pod.GetPodNumPerNs).
		//deployment操作
		GET("/api/k8s/deployments", Deployment.GetDeployment).
		GET("/api/k8s/deployment/detail", Deployment.GetDeloymentDetail).
		DELETE("/api/k8s/deployment/del", Deployment.DeleteDeloyment).
		PUT("/api/k8s/deployment/update", Deployment.UpdateDeloyment).
		PUT("/api/k8s/deployment/restart", Deployment.RestartDeployment).
		PUT("/api/k8s/deployment/scale", Deployment.ScaleDeployment).
		POST("/api/k8s/deployment/create", Deployment.CreateDeployment).
		GET("/api/k8s/deployment/numns", Deployment.GetDeloymentNumPerNs)
}
