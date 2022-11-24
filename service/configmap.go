package service

import (
	"context"
	"errors"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
)

var Configmap configmap

type configmap struct{}

// 定义configmaps的返回内容 items是configmap列表吗total为configmap元素总数
type ConfigmapsResp struct {
	Item  []corev1.ConfigMap `json:"items"`
	Total int                `json:"total"`
}

// 获取configmap列表
func (p *configmap) GetConfigmap(filterName, namespace string, limit, page int) (configmapsResp *ConfigmapsResp, err error) {
	//通过clientset获取configmaps完整列表
	configmapList, err := K8s.ClientSet.CoreV1().ConfigMaps(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error("获取configmap列表失败", err)
		return nil, errors.New("获取configmap列表失败" + err.Error())
	}
	//实例化DataSelector对象
	selectableData := &DataSelector{
		GenericDataList: p.toCells(configmapList.Items),
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
	//再将datacell切片数据转成原生configmap切片
	configmaps := p.fromCells(data.GenericDataList)
	//返回
	return &ConfigmapsResp{
		Item:  configmaps,
		Total: total,
	}, nil
}

// 获取configmap详情
func (p *configmap) GetConfigmapDetail(configmapName, namespace string) (configmap *corev1.ConfigMap, err error) {
	configmap, err = K8s.ClientSet.CoreV1().ConfigMaps(namespace).Get(context.TODO(), configmapName, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取Configmap详情失败" + err.Error())
		return nil, errors.New("获取Configmap详情失败" + err.Error())
	}
	return configmap, nil
}

// 删除configmap
func (p *configmap) DeleteConfigmap(configmapName, namespace string) (err error) {
	err = K8s.ClientSet.CoreV1().ConfigMaps(namespace).Delete(context.TODO(), configmapName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error("删除Configmap失败" + err.Error())
		return errors.New("获取Configmap详情失败" + err.Error())
	}
	return nil
}

// 更新configmap
func (p *configmap) UpdateConfigmap(namespace, content string) (err error) {
	//将content反序列化成为configmap对象
	var configmap = &corev1.ConfigMap{}
	if err = json.Unmarshal([]byte(content), configmap); err != nil {
		logger.Error("Content反序列化失败", err)
		return errors.New("Content反序列化失败" + err.Error())
	}
	//更新configmap
	_, err = K8s.ClientSet.CoreV1().ConfigMaps(namespace).Update(context.TODO(), configmap, metav1.UpdateOptions{})
	if err != nil {
		logger.Error("更新Configmap失败" + err.Error())
		return errors.New("更新Configmap失败" + err.Error())
	}
	return nil
}

// 把configmapCell转成appsv1 configmap
func (p *configmap) fromCells(cells []DataCell) []corev1.ConfigMap {
	configmaps := make([]corev1.ConfigMap, len(cells))
	for i := range cells {
		configmaps[i] = corev1.ConfigMap(cells[i].(configMapCell))
	}
	return configmaps
}

// 把appsv1 configmaps转成datacell
func (p *configmap) toCells(std []corev1.ConfigMap) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = configMapCell(std[i])

	}
	return cells
}
