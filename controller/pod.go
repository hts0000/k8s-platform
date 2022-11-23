package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"k8s-platform/service"
	"net/http"
)

var Pod pod

type pod struct{}

// pod列表支持过滤。排序。分页
func (p *pod) GetPods(ctx *gin.Context) {
	//匿名结构体用于定义入参,get请求为from格式其他请求为json格式
	params := new(struct {
		FilterName string `form:"filter_name"`
		Namespace  string `form:"namespace"`
		Page       int    `form:"page"`
		Limit      int    `form:"limit"'`
	})
	//绑定参数给匿名结构体的属性赋值
	//from格式使用ctx.bind方法
	//json格式使用ctx.ShouldBind
	if err := ctx.Bind(params); err != nil {
		logger.Error("绑定参数失败", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	//调用service方法获取数据
	data, err := service.Pod.GetPods(params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		logger.Error("获取数据错误" + err.Error())
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取pod列表成功",
		"data": data,
	})

}

// Pod详情
func (p *pod) GetPodDetail(ctx *gin.Context) {
	//匿名结构体用于定义入参,get请求为from格式其他请求为json格式
	params := new(struct {
		PodName   string `form:"pod_name"`
		Namespace string `form:"namespace"`
	})
	//绑定参数给匿名结构体的属性赋值
	//from格式使用ctx.bind方法
	//json格式使用ctx.ShouldBind
	if err := ctx.Bind(params); err != nil {
		logger.Error("绑定参数失败", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	//调用service方法获取数据
	data, err := service.Pod.GetPodDetail(params.PodName, params.Namespace)
	if err != nil {
		logger.Error("获取pod详情错误" + err.Error())
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取pod详情成功",
		"data": data,
	})

}

// 删除pod
func (p *pod) DeletePod(ctx *gin.Context) {
	//匿名结构体用于定义入参,get请求为from格式其他请求为json格式
	params := new(struct {
		PodName   string `json:"pod_name"`
		Namespace string `json:"namespace"`
	})
	//绑定参数给匿名结构体的属性赋值
	//from格式使用ctx.bind方法
	//json格式使用ctx.ShouldBind
	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error("删除pod失败", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	//调用service方法获取数据
	err := service.Pod.DeletePod(params.PodName, params.Namespace)
	if err != nil {
		logger.Error("删除pod失败" + err.Error())
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "删除pod成功",
		"data": nil,
	})

}

// 更新pod
func (p *pod) UpdatePod(ctx *gin.Context) {
	//匿名结构体用于定义入参,get请求为from格式其他请求为json格式
	params := new(struct {
		Content   string `json:"content"`
		Namespace string `json:"namespace"`
	})
	//绑定参数给匿名结构体的属性赋值
	//from格式使用ctx.bind方法
	//json格式使用ctx.ShouldBind
	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error("绑定参数失败", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	//调用service方法获取数据
	err := service.Pod.UpdatePod(params.Content, params.Namespace)
	if err != nil {
		logger.Error("更新pod失败" + err.Error())
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "更新pod成功",
		"data": nil,
	})

}

// 获取pod的容器名
func (p *pod) GetPodContainer(ctx *gin.Context) {
	//匿名结构体用于定义入参,get请求为from格式其他请求为json格式
	params := new(struct {
		PodName   string `from:"pod_name"`
		Namespace string `from:"namespace"`
	})
	//绑定参数给匿名结构体的属性赋值
	//from格式使用ctx.bind方法
	//json格式使用ctx.ShouldBind
	if err := ctx.Bind(params); err != nil {
		logger.Error("绑定参数失败", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	//调用service方法获取数据
	data, err := service.Pod.GetPodConatiner(params.PodName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取pod容器名成功",
		"data": data,
	})

}

// 获取pod日志
func (p *pod) GetPodLog(ctx *gin.Context) {
	//匿名结构体用于定义入参,get请求为from格式其他请求为json格式
	params := new(struct {
		ContainerName string `from:"container"`
		PodName       string `form:"pod_name"`
		Namespace     string `form:"namespace"`
	})
	//绑定参数给匿名结构体的属性赋值
	//from格式使用ctx.bind方法
	//json格式使用ctx.ShouldBind
	if err := ctx.Bind(params); err != nil {
		logger.Error("绑定参数失败", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	//调用service方法获取数据
	data, err := service.Pod.GetPodLog(params.ContainerName, params.PodName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取pod日志成功",
		"data": data,
	})

}

// 获取pod每个名称空间数量
func (p *pod) GetPodNumPerNs(ctx *gin.Context) {

	//调用service方法获取数据
	data, err := service.Pod.GetPodNumPerNs()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取pod数量成功",
		"data": data,
	})

}
