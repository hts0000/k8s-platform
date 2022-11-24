package service

import (
	"context"
	"errors"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Pv pv

type pv struct{}

// 定义pvs的返回内容 items是pv列表吗total为pv元素总数
type PvsResp struct {
	Item  []corev1.PersistentVolume `json:"items"`
	Total int                       `json:"total"`
}

// 获取pv列表
func (p *pv) GetPv(filterName string, limit, page int) (pvsResp *PvsResp, err error) {
	//通过clientset获取pvs完整列表
	pvList, err := K8s.ClientSet.CoreV1().PersistentVolumes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error("获取pv列表失败", err)
		return nil, errors.New("获取pv列表失败" + err.Error())
	}
	//实例化DataSelector对象
	selectableData := &DataSelector{
		GenericDataList: p.toCells(pvList.Items),
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
	//再将datacell切片数据转成原生pv切片
	pvs := p.fromCells(data.GenericDataList)
	//返回
	return &PvsResp{
		Item:  pvs,
		Total: total,
	}, nil
}

// 获取pv详情
func (p *pv) GetPvDetail(pvName string) (pv *corev1.PersistentVolume, err error) {
	pv, err = K8s.ClientSet.CoreV1().PersistentVolumes().Get(context.TODO(), pvName, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取Pv详情失败" + err.Error())
		return nil, errors.New("获取Pv详情失败" + err.Error())
	}
	return pv, nil
}

// 删除pv
func (p *pv) DeletePv(pvName string) (err error) {
	err = K8s.ClientSet.CoreV1().PersistentVolumes().Delete(context.TODO(), pvName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error("删除Pv失败" + err.Error())
		return errors.New("获取Pv详情失败" + err.Error())
	}
	return nil
}

// 把pvCell转成corev1 pv
func (p *pv) fromCells(cells []DataCell) []corev1.PersistentVolume {
	pv := make([]corev1.PersistentVolume, len(cells))
	for i := range cells {
		pv[i] = corev1.PersistentVolume(cells[i].(pvCell))
	}
	return pv
}

// 把corev1.pv转成datacell
func (p *pv) toCells(std []corev1.PersistentVolume) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = pvCell(std[i])

	}
	return cells
}
