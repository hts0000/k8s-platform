package service

import (
	"context"
	"errors"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
)

var Pvc pvc

type pvc struct{}

// 定义pvcs的返回内容 items是pvc列表吗total为pvc元素总数
type PvcsResp struct {
	Item  []corev1.PersistentVolumeClaim `json:"items"`
	Total int                            `json:"total"`
}

// 获取pvc列表
func (p *pvc) GetPvc(filterName, namespace string, limit, page int) (pvcsResp *PvcsResp, err error) {
	//通过clientset获取pvcs完整列表
	pvcList, err := K8s.ClientSet.CoreV1().PersistentVolumeClaims(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error("获取pvc列表失败", err)
		return nil, errors.New("获取pvc列表失败" + err.Error())
	}
	//实例化DataSelector对象
	selectableData := &DataSelector{
		GenericDataList: p.toCells(pvcList.Items),
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
	//再将datacell切片数据转成原生pvc切片
	pvcs := p.fromCells(data.GenericDataList)
	//返回
	return &PvcsResp{
		Item:  pvcs,
		Total: total,
	}, nil
}

// 获取pvc详情
func (p *pvc) GetPvcDetail(pvcName, namespace string) (pvc *corev1.PersistentVolumeClaim, err error) {
	pvc, err = K8s.ClientSet.CoreV1().PersistentVolumeClaims(namespace).Get(context.TODO(), pvcName, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取Pvc详情失败" + err.Error())
		return nil, errors.New("获取Pvc详情失败" + err.Error())
	}
	return pvc, nil
}

// 删除pvc
func (p *pvc) DeletePvc(pvcName, namespace string) (err error) {
	err = K8s.ClientSet.CoreV1().PersistentVolumeClaims(namespace).Delete(context.TODO(), pvcName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error("删除Pvc失败" + err.Error())
		return errors.New("获取Pvc详情失败" + err.Error())
	}
	return nil
}

// 更新pvc
func (p *pvc) UpdatePvc(namespace, content string) (err error) {
	//将content反序列化成为pvc对象
	var pvc = &corev1.PersistentVolumeClaim{}
	if err = json.Unmarshal([]byte(content), pvc); err != nil {
		logger.Error("Content反序列化失败", err)
		return errors.New("Content反序列化失败" + err.Error())
	}
	//更新pvc
	_, err = K8s.ClientSet.CoreV1().PersistentVolumeClaims(namespace).Update(context.TODO(), pvc, metav1.UpdateOptions{})
	if err != nil {
		logger.Error("更新Pvc失败" + err.Error())
		return errors.New("更新Pvc失败" + err.Error())
	}
	return nil
}

// 把pvcCell转成appsv1 pvc
func (p *pvc) fromCells(cells []DataCell) []corev1.PersistentVolumeClaim {
	pvcs := make([]corev1.PersistentVolumeClaim, len(cells))
	for i := range cells {
		pvcs[i] = corev1.PersistentVolumeClaim(cells[i].(pvcCell))
	}
	return pvcs
}

// 把appsv1 pvcs转成datacell
func (p *pvc) toCells(std []corev1.PersistentVolumeClaim) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = pvcCell(std[i])

	}
	return cells
}
