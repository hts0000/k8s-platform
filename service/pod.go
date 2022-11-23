package service

import (
	"bytes"
	"context"
	"errors"
	"github.com/wonderivan/logger"
	"io"
	"k8s-platform/config"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
)

var Pod pod

type pod struct{}

// 定义pods的返回内容 items是pod列表吗total为pod元素总数
type PodsResp struct {
	Item  []corev1.Pod `json:"items"`
	Total int          `json:"total"`
}

type PodsNs struct {
	Namespace string `json:"namespace"`
	PodNum    int    `json:"pod_num"`
}

// 获取pod列表
func (p *pod) GetPods(filterName, namespace string, limit, page int) (podsResp *PodsResp, err error) {
	//通过clientset获取pods完整列表
	podList, err := K8s.ClientSet.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error("获取pod列表失败", err)
		return nil, errors.New("获取pod列表失败" + err.Error())
	}
	//实例化DataSelector对象
	selectableData := &DataSelector{
		GenericDataList: p.toCells(podList.Items),
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
	//再将datacell切片数据转成原生pod切片
	pods := p.fromCells(data.GenericDataList)
	//返回
	return &PodsResp{
		Item:  pods,
		Total: total,
	}, nil
}

// 获取pod详情
func (p *pod) GetPodDetail(podName, namespace string) (pod *corev1.Pod, err error) {
	pod, err = K8s.ClientSet.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取Pod详情失败" + err.Error())
		return nil, errors.New("获取Pod详情失败" + err.Error())
	}
	return pod, nil
}

// 删除pod
func (p *pod) DeletePod(podName, namespace string) (err error) {
	err = K8s.ClientSet.CoreV1().Pods(namespace).Delete(context.TODO(), podName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error("删除Pod失败" + err.Error())
		return errors.New("获取Pod详情失败" + err.Error())
	}
	return nil
}

// 更新pod
func (p *pod) UpdatePod(namespace, content string) (err error) {
	//将content反序列化成为pod对象
	var pod = &corev1.Pod{}
	if err = json.Unmarshal([]byte(content), pod); err != nil {
		logger.Error("Content反序列化失败", err)
		return errors.New("Content反序列化失败" + err.Error())
	}
	//更新pod
	_, err = K8s.ClientSet.CoreV1().Pods(namespace).Update(context.TODO(), pod, metav1.UpdateOptions{})
	if err != nil {
		logger.Error("更新Pod失败" + err.Error())
		return errors.New("更新Pod失败" + err.Error())
	}
	return nil
}

// 获取pod日志
func (p *pod) GetPodLog(containerName, podName, namespace string) (log string, err error) {
	//设置日志的配置。容器名，tail行数
	lineLimit := int64(config.PodLogTailLine)
	option := &corev1.PodLogOptions{
		Container: containerName,
		TailLines: &lineLimit,
	}
	//获取request实例
	req := K8s.ClientSet.CoreV1().Pods(namespace).GetLogs(podName, option)
	//发起request请求。返回一个io.readcloser类型的，等同于response.body
	podLogs, err := req.Stream(context.TODO())
	if err != nil {
		logger.Error("获取podlog失败", err)
		return "", errors.New("获取podlog失败" + err.Error())
	}
	//将request body写入到缓冲区，目的是为了转成string返回
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		logger.Error("复制podlog失败", err)
		return "", errors.New("复制podlog失败" + err.Error())
	}
	return buf.String(), nil
}

// 获取pod中的容器，日志，终端功能使用
func (p *pod) GetPodConatiner(podName, namespace string) (containers []string, err error) {
	//获取pod详情
	pod, err := p.GetPodDetail(podName, namespace)
	if err != nil {
		return nil, err
	}
	//从pod对象中拿到容器名
	for _, container := range pod.Spec.Containers {
		containers = append(containers, container.Name)
	}
	return containers, nil
}

// 获取每个命名空间pod数量
func (p *pod) GetPodNumPerNs() (podsNss []*PodsNs, err error) {
	//获取namespace列表
	namespaceList, err := K8s.ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error("获取namespace列表失败", err)
		return nil, errors.New("获取namespace列表失败" + err.Error())
	}
	//for循环
	for _, namespace := range namespaceList.Items {
		//获取pod列表
		podList, err := K8s.ClientSet.CoreV1().Pods(namespace.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			logger.Error("获取pod列表失败", err)
			return nil, errors.New("获取pod列表失败" + err.Error())
		}
		//组装数据
		podsNs := &PodsNs{
			Namespace: namespace.Name,
			PodNum:    len(podList.Items),
		}
		podsNss = append(podsNss, podsNs)
	}
	return podsNss, nil
}

// 把podCell转成corev1 pod
func (p *pod) fromCells(cells []DataCell) []corev1.Pod {
	pods := make([]corev1.Pod, len(cells))
	for i := range cells {
		pods[i] = corev1.Pod(cells[i].(podCell))
	}
	return pods
}

// 把corev1 pods转成datacell
func (p *pod) toCells(std []corev1.Pod) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = podCell(std[i])

	}
	return cells
}
