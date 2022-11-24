package service

import (
	"context"
	"errors"
	"github.com/wonderivan/logger"
	nwv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
)

var Ingress ingress

type ingress struct{}

// 定义ingresss的返回内容 items是ingress列表吗total为ingress元素总数
type IngresssResp struct {
	Item  []nwv1.Ingress `json:"items"`
	Total int            `json:"total"`
}

// 获取ingress列表
func (p *ingress) GetIngress(filterName, namespace string, limit, page int) (ingresssResp *IngresssResp, err error) {
	//通过clientset获取ingresss完整列表
	ingressList, err := K8s.ClientSet.NetworkingV1().Ingresses(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error("获取ingress列表失败", err)
		return nil, errors.New("获取ingress列表失败" + err.Error())
	}
	//实例化DataSelector对象
	selectableData := &DataSelector{
		GenericDataList: p.toCells(ingressList.Items),
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
	//再将datacell切片数据转成原生ingress切片
	ingresss := p.fromCells(data.GenericDataList)
	//返回
	return &IngresssResp{
		Item:  ingresss,
		Total: total,
	}, nil
}

// 获取ingress详情
func (p *ingress) GetIngressDetail(ingressName, namespace string) (ingress *nwv1.Ingress, err error) {
	ingress, err = K8s.ClientSet.NetworkingV1().Ingresses(namespace).Get(context.TODO(), ingressName, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取Ingress详情失败" + err.Error())
		return nil, errors.New("获取Ingress详情失败" + err.Error())
	}
	return ingress, nil
}

// 删除ingress
func (p *ingress) DeleteIngress(ingressName, namespace string) (err error) {
	err = K8s.ClientSet.NetworkingV1().Ingresses(namespace).Delete(context.TODO(), ingressName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error("删除Ingress失败" + err.Error())
		return errors.New("获取Ingress详情失败" + err.Error())
	}
	return nil
}

// 更新ingress
func (p *ingress) UpdateIngress(namespace, content string) (err error) {
	//将content反序列化成为ingress对象
	var ingress = &nwv1.Ingress{}
	if err = json.Unmarshal([]byte(content), ingress); err != nil {
		logger.Error("Content反序列化失败", err)
		return errors.New("Content反序列化失败" + err.Error())
	}
	//更新ingress
	_, err = K8s.ClientSet.NetworkingV1().Ingresses(namespace).Update(context.TODO(), ingress, metav1.UpdateOptions{})
	if err != nil {
		logger.Error("更新Ingress失败" + err.Error())
		return errors.New("更新Ingress失败" + err.Error())
	}
	return nil
}

// 把ingressCell转成appsv1 ingress
func (p *ingress) fromCells(cells []DataCell) []nwv1.Ingress {
	ingresss := make([]nwv1.Ingress, len(cells))
	for i := range cells {
		ingresss[i] = nwv1.Ingress(cells[i].(ingressCell))
	}
	return ingresss
}

// 把appsv1 ingresss转成datacell
func (p *ingress) toCells(std []nwv1.Ingress) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = ingressCell(std[i])

	}
	return cells
}
