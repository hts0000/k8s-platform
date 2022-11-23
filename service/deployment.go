package service

import (
	"context"
	"errors"
	"github.com/wonderivan/logger"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/json"
	"strconv"
	"time"
)

var Deployment deployment

type deployment struct{}

// 定义deployments的返回内容 items是deployment列表吗total为deployment元素总数
type DeploymentsResp struct {
	Item  []appsv1.Deployment `json:"items"`
	Total int                 `json:"total"`
}

type DeploymentsNs struct {
	Namespace     string `json:"namespace"`
	DeploymentNum int    `json:"deployment_num"`
}

// 定义结构体用于创建deployment
type DeployCreate struct {
	Name          string            `json:"name"`
	Namespace     string            `json:"namespace"`
	Replicas      int32             `json:"replicas"`
	Image         string            `json:"image"`
	Label         map[string]string `json:"label"`
	Cpu           string            `json:"cpu"`
	Memory        string            `json:"memory"`
	ContainerPort int32             `json:"container_port"`
	HealthCheck   bool              `json:"health_check"`
	HealthPath    string            `json:"health_path"`
}

// 获取deployment列表
func (p *deployment) GetDeployments(filterName, namespace string, limit, page int) (deploymentsResp *DeploymentsResp, err error) {
	//通过clientset获取deployments完整列表
	deploymentList, err := K8s.ClientSet.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error("获取deployment列表失败", err)
		return nil, errors.New("获取deployment列表失败" + err.Error())
	}
	//实例化DataSelector对象
	selectableData := &DataSelector{
		GenericDataList: p.toCells(deploymentList.Items),
		DataSelectQuery: &DataSelect{
			FilterQuery: &Filter{filterName},
			PaginateQuery: &Paginate{
				Limit: limit,
				Page:  page,
			},
		},
	}
	//先过滤
	filtered := selectableData.Filter()
	//再拿total
	total := len(filtered.GenericDataList)
	//在排序和分页
	data := filtered.Sort().Paginate()
	//再将datacell切片数据转成原生deployment切片
	deployments := p.fromCells(data.GenericDataList)
	//返回
	return &DeploymentsResp{
		Item:  deployments,
		Total: total,
	}, nil
}

// 获取deployment详情
func (p *deployment) GetDeploymentDetail(deploymentName, namespace string) (deployment *appsv1.Deployment, err error) {
	deployment, err = K8s.ClientSet.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取Deployment详情失败" + err.Error())
		return nil, errors.New("获取Deployment详情失败" + err.Error())
	}
	return deployment, nil
}

// 删除deployment
func (p *deployment) DeleteDeployment(deploymentName, namespace string) (err error) {
	err = K8s.ClientSet.AppsV1().Deployments(namespace).Delete(context.TODO(), deploymentName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error("删除Deployment失败" + err.Error())
		return errors.New("获取Deployment详情失败" + err.Error())
	}
	return nil
}

// 更新deployment
func (p *deployment) UpdateDeployment(namespace, content string) (err error) {
	//将content反序列化成为deployment对象
	var deployment = &appsv1.Deployment{}
	if err = json.Unmarshal([]byte(content), deployment); err != nil {
		logger.Error("Content反序列化失败", err)
		return errors.New("Content反序列化失败" + err.Error())
	}
	//更新deployment
	_, err = K8s.ClientSet.AppsV1().Deployments(namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
	if err != nil {
		logger.Error("更新Deployment失败" + err.Error())
		return errors.New("更新Deployment失败" + err.Error())
	}
	return nil
}

// 修改deployment副本数
func (p *deployment) ScaleDeployment(deploymentName, namespace string, scaleNum int) (replicas int32, err error) {
	//获取autoscaling.scale对象，能点出当前的副本数
	scale, err := K8s.ClientSet.AppsV1().Deployments(namespace).GetScale(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取deployment副本数失败", err.Error())
		return 0, errors.New("获取deployment副本数失败" + err.Error())
	}
	//修改副本数
	scale.Spec.Replicas = int32(scaleNum)
	//更新副本数
	newScale, err := K8s.ClientSet.AppsV1().Deployments(namespace).UpdateScale(context.TODO(), deploymentName, scale, metav1.UpdateOptions{})
	if err != nil {
		logger.Error("更新deployment副本数失败", err.Error())
		return 0, errors.New("更新deployment副本数失败" + err.Error())
	}
	return newScale.Spec.Replicas, nil
}

// 重启deployment
func (p *deployment) RestartDeployment(deploymentName, namespace string) (err error) {
	//使用patchData map 组装数据
	patchData := map[string]interface{}{
		"spec": map[string]interface{}{
			"template": map[string]interface{}{
				"spec": map[string]interface{}{
					"containers": []map[string]interface{}{
						{"name": deploymentName,
							"env": []map[string]string{{
								"name":  "RESTART_",
								"value": strconv.FormatInt(time.Now().Unix(), 10),
							}},
						},
					},
				},
			},
		},
	}
	//序列化为字节。因为patch方法值接收字节类型参数
	patchByte, err := json.Marshal(patchData)
	if err != nil {
		logger.Error("patchdata序列化失败", err)
		return errors.New("patchdata序列化失败" + err.Error())
	}
	//调用patch方法更新deployment副本数
	_, err = K8s.ClientSet.AppsV1().Deployments(namespace).Patch(context.TODO(), deploymentName, "application/strategic-merge-patch+json", patchByte, metav1.PatchOptions{})
	if err != nil {
		logger.Error("修改deployment副本数失败", err)
		return errors.New("修改deployment副本数失败" + err.Error())
	}
	return nil
}

// 创建deployment
func (p *deployment) CreateDeployment(data DeployCreate) (err error) {
	//将data中的数据组装成appsv1.deployment对象
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      data.Name,
			Namespace: data.Namespace,
			Labels:    data.Label,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &data.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: data.Label,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:   data.Name,
					Labels: data.Label,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  data.Name,
							Image: data.Image,
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: data.ContainerPort,
								},
							},
						},
					},
				},
			},
		},
		Status: appsv1.DeploymentStatus{},
	}
	//判断是否打开健康检测功能，若打开，则规定ReadinessProbe和LivenessProbe
	if data.HealthCheck {
		deployment.Spec.Template.Spec.Containers[0].ReadinessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: data.HealthPath,
					Port: intstr.IntOrString{
						Type:   0,
						IntVal: data.ContainerPort,
					},
				},
			},
			InitialDelaySeconds: 5,
			TimeoutSeconds:      5,
			PeriodSeconds:       5,
		}
		deployment.Spec.Template.Spec.Containers[0].LivenessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: data.HealthPath,
					Port: intstr.IntOrString{
						Type:   0,
						IntVal: data.ContainerPort,
					},
				},
			},
			InitialDelaySeconds: 15,
			TimeoutSeconds:      15,
			PeriodSeconds:       15,
		}
	}
	//定义容器的limit和request
	deployment.Spec.Template.Spec.Containers[0].Resources.Limits = map[corev1.ResourceName]resource.Quantity{
		corev1.ResourceCPU:    resource.MustParse(data.Cpu),
		corev1.ResourceMemory: resource.MustParse(data.Memory),
	}
	deployment.Spec.Template.Spec.Containers[0].Resources.Requests = map[corev1.ResourceName]resource.Quantity{
		corev1.ResourceCPU:    resource.MustParse(data.Cpu),
		corev1.ResourceMemory: resource.MustParse(data.Memory),
	}
	//创建deployment
	_, err = K8s.ClientSet.AppsV1().Deployments(data.Namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		logger.Error("创建deployment失败", err)
		return errors.New("创建deployment失败" + err.Error())
	}
	return nil
}

// 获取每个命名空间deployment数量
func (p *deployment) GetDeploymentNumPerNs() (deploymentsNss []*DeploymentsNs, err error) {
	//获取namespace列表
	namespaceList, err := K8s.ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error("获取namespace列表失败", err)
		return nil, errors.New("获取namespace列表失败" + err.Error())
	}
	//for循环
	for _, namespace := range namespaceList.Items {
		//获取deployment列表
		deploymentList, err := K8s.ClientSet.AppsV1().Deployments(namespace.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			logger.Error("获取deployment列表失败", err)
			return nil, errors.New("获取deployment列表失败" + err.Error())
		}
		//组装数据
		deploymentsNs := &DeploymentsNs{
			Namespace:     namespace.Name,
			DeploymentNum: len(deploymentList.Items),
		}
		deploymentsNss = append(deploymentsNss, deploymentsNs)
	}
	return deploymentsNss, nil
}

// 把deploymentCell转成appsv1 deployment
func (p *deployment) fromCells(cells []DataCell) []appsv1.Deployment {
	deployments := make([]appsv1.Deployment, len(cells))
	for i := range cells {
		deployments[i] = appsv1.Deployment(cells[i].(deploymentCell))
	}
	return deployments
}

// 把appsv1 deployments转成datacell
func (p *deployment) toCells(std []appsv1.Deployment) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = deploymentCell(std[i])

	}
	return cells
}
