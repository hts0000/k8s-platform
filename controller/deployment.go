package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"k8s-platform/service"
	"net/http"
)

var Deployment deployment

type deployment struct{}

// deloyment列表支持过滤。排序。分页
func (p *deployment) GetDeployment(ctx *gin.Context) {
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
	data, err := service.Deployment.GetDeployments(params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		logger.Error("获取数据错误" + err.Error())
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取deloyment列表成功",
		"data": data,
	})

}

// Deloyment详情
func (p *deployment) GetDeloymentDetail(ctx *gin.Context) {
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
	data, err := service.Deployment.GetDeploymentDetail(params.DeloymentName, params.Namespace)
	if err != nil {
		logger.Error("获取deloyment详情错误" + err.Error())
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取deloyment详情成功",
		"data": data,
	})

}

// 删除deloyment
func (p *deployment) DeleteDeloyment(ctx *gin.Context) {
	//匿名结构体用于定义入参,get请求为from格式其他请求为json格式
	params := new(struct {
		DeloymentName string `json:"deloyment_name"`
		Namespace     string `json:"namespace"`
	})
	//绑定参数给匿名结构体的属性赋值
	//from格式使用ctx.bind方法
	//json格式使用ctx.ShouldBind
	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error("删除deloyment失败", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	//调用service方法获取数据
	err := service.Deployment.DeleteDeployment(params.DeloymentName, params.Namespace)
	if err != nil {
		logger.Error("删除deloyment失败" + err.Error())
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "删除deloyment成功",
		"data": nil,
	})

}

// 更新deloyment
func (p *deployment) UpdateDeloyment(ctx *gin.Context) {
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
	err := service.Deployment.UpdateDeployment(params.Content, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "更新deployment失败",
			"data": nil,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "更新deloyment成功",
		"data": nil,
	})

}

// 修改deployment副本数
func (p *deployment) ScaleDeployment(ctx *gin.Context) {
	//匿名结构体用于定义入参,get请求为from格式其他请求为json格式
	params := new(struct {
		ScaleNum       int    `json:"scale_num"`
		DeploymentName string `json:"deployment_name"`
		Namespace      string `json:"namespace"`
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
	replicas, err := service.Deployment.ScaleDeployment(params.DeploymentName, params.Namespace, params.ScaleNum)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": replicas,
		})
		return
	}
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"msg":  "修改deployment成功",
		"data": replicas,
	})
	return

}

// 重启deployment
func (p *deployment) RestartDeployment(ctx *gin.Context) {
	//匿名结构体用于定义入参,get请求为from格式其他请求为json格式
	params := new(struct {
		DeploymentName string `json:"deployment_name"`
		Namespace      string `json:"namespace"`
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
	err := service.Deployment.RestartDeployment(params.DeploymentName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"msg":  "重启deployment成功",
		"data": nil,
	})
	return

}

// 创建deployment
func (p *deployment) CreateDeployment(ctx *gin.Context) {
	//匿名结构体用于定义入参,get请求为from格式其他请求为json格式
	deployCreate := service.DeployCreate{}
	//绑定参数给匿名结构体的属性赋值
	//from格式使用ctx.bind方法
	//json格式使用ctx.ShouldBind
	if err := ctx.ShouldBindJSON(deployCreate); err != nil {
		logger.Error("绑定参数失败", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	//调用service方法获取数据
	err := service.Deployment.CreateDeployment(deployCreate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"msg":  "创建deployment成功",
		"data": nil,
	})
	return

}

// 获取deloyment每个名称空间数量
func (p *deployment) GetDeloymentNumPerNs(ctx *gin.Context) {

	//调用service方法获取数据
	data, err := service.Deployment.GetDeploymentNumPerNs()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取deloyment数量成功",
		"data": data,
	})

}
