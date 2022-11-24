package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"k8s-platform/service"
	"net/http"
)

var StatefulSet statefulSet

type statefulSet struct{}

// deloyment列表支持过滤。排序。分页
func (p *statefulSet) GetStatefulSet(ctx *gin.Context) {
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
	data, err := service.StatefulSet.GetStatefulSets(params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		logger.Error("获取数据错误" + err.Error())
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取statefulSet列表成功",
		"data": data,
	})

}

// Deloyment详情
func (p *statefulSet) GetStatefulSetDetail(ctx *gin.Context) {
	//匿名结构体用于定义入参,get请求为from格式其他请求为json格式
	params := new(struct {
		DeloymentName string `form:"deloyment_name"`
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
	data, err := service.StatefulSet.GetStatefulSetDetail(params.DeloymentName, params.Namespace)
	if err != nil {
		logger.Error("获取statefulSet详情错误" + err.Error())
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取statefulSet详情成功",
		"data": data,
	})

}

// 删除deloyment
func (p *statefulSet) DeleteStatefulSet(ctx *gin.Context) {
	//匿名结构体用于定义入参,get请求为from格式其他请求为json格式
	params := new(struct {
		StatefulSetName string `json:"statefulSetName"`
		Namespace       string `json:"namespace"`
	})
	//绑定参数给匿名结构体的属性赋值
	//from格式使用ctx.bind方法
	//json格式使用ctx.ShouldBind
	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error("删除statefulSet失败", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	//调用service方法获取数据
	err := service.StatefulSet.DeleteStatefulSet(params.StatefulSetName, params.Namespace)
	if err != nil {
		logger.Error("删除statefulSet失败" + err.Error())
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "删除statefulSet成功",
		"data": nil,
	})

}

// 更新deloyment
func (p *statefulSet) UpdateStatefulSet(ctx *gin.Context) {
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
	err := service.StatefulSet.UpdateStatefulSet(params.Content, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "更新statefulSet失败",
			"data": nil,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "更新statefulSet失败",
		"data": nil,
	})

}
