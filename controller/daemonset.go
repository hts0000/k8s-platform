package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"k8s-platform/service"
	"net/http"
)

var Daemonset daemonSet

type daemonSet struct{}

// daemonset列表支持过滤。排序。分页
func (p *daemonSet) GetDaemonSet(ctx *gin.Context) {
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
	data, err := service.DaemonSet.GetDaemonSet(params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		logger.Error("获取数据错误" + err.Error())
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取daemonset列表成功",
		"data": data,
	})

}

// Daemonset详情
func (p *daemonSet) GetDaemonSetDetail(ctx *gin.Context) {
	//匿名结构体用于定义入参,get请求为from格式其他请求为json格式
	params := new(struct {
		DaemonSetName string `from:"daemonset_name"`
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
	data, err := service.DaemonSet.GetDaemonSetDetail(params.DaemonSetName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": data,
		})

	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取DaemonSet详情成功",
		"data": data,
	})

}

// 删除daemonset
func (p *daemonSet) DeleteDaemonSet(ctx *gin.Context) {
	//匿名结构体用于定义入参,get请求为from格式其他请求为json格式
	params := new(struct {
		DaemonsetName string `json:"daemonset_name"`
		Namespace     string `json:"namespace"`
	})
	//绑定参数给匿名结构体的属性赋值
	//from格式使用ctx.bind方法
	//json格式使用ctx.ShouldBind
	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error("删除daemonset失败", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	//调用service方法获取数据
	err := service.DaemonSet.DeleteDaemonSet(params.DaemonsetName, params.Namespace)
	if err != nil {
		logger.Error("删除daemonset失败" + err.Error())
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "删除daemonset成功",
		"data": nil,
	})

}

// 更新daemonset
func (p *daemonSet) UpdateDaemonSet(ctx *gin.Context) {
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
	err := service.DaemonSet.UpdateDaemonSet(params.Content, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "更新daemonSet失败",
			"data": nil,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "更新daemonset成功",
		"data": nil,
	})

}
