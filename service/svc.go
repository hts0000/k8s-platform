package service

import (
	"context"
	"errors"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
)

var Svc svc

type svc struct{}

// 定义svcs的返回内容 items是svc列表吗total为svc元素总数
type SvcsResp struct {
	Item  []corev1.Service `json:"items"`
	Total int              `json:"total"`
}

// 获取svc列表
func (p *svc) GetSvc(filterName, namespace string, limit, page int) (svcsResp *SvcsResp, err error) {
	//通过clientset获取svcs完整列表
	svcList, err := K8s.ClientSet.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error("获取svc列表失败", err)
		return nil, errors.New("获取svc列表失败" + err.Error())
	}
	//实例化DataSelector对象
	selectableData := &DataSelector{
		GenericDataList: p.toCells(svcList.Items),
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
	//再将datacell切片数据转成原生svc切片
	svcs := p.fromCells(data.GenericDataList)
	//返回
	return &SvcsResp{
		Item:  svcs,
		Total: total,
	}, nil
}

// 获取svc详情
func (p *svc) GetSvcDetail(svcName, namespace string) (svc *corev1.Service, err error) {
	svc, err = K8s.ClientSet.CoreV1().Services(namespace).Get(context.TODO(), svcName, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取Svc详情失败" + err.Error())
		return nil, errors.New("获取Svc详情失败" + err.Error())
	}
	return svc, nil
}

// 删除svc
func (p *svc) DeleteSvc(svcName, namespace string) (err error) {
	err = K8s.ClientSet.CoreV1().Services(namespace).Delete(context.TODO(), svcName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error("删除Svc失败" + err.Error())
		return errors.New("获取Svc详情失败" + err.Error())
	}
	return nil
}

// 更新svc
func (p *svc) UpdateSvc(namespace, content string) (err error) {
	//将content反序列化成为svc对象
	var svc = &corev1.Service{}
	if err = json.Unmarshal([]byte(content), svc); err != nil {
		logger.Error("Content反序列化失败", err)
		return errors.New("Content反序列化失败" + err.Error())
	}
	//更新svc
	_, err = K8s.ClientSet.CoreV1().Services(namespace).Update(context.TODO(), svc, metav1.UpdateOptions{})
	if err != nil {
		logger.Error("更新Svc失败" + err.Error())
		return errors.New("更新Svc失败" + err.Error())
	}
	return nil
}

// 把svcCell转成corev1 svc
func (p *svc) fromCells(cells []DataCell) []corev1.Service {
	svc := make([]corev1.Service, len(cells))
	for i := range cells {
		svc[i] = corev1.Service(cells[i].(serviceCell))
	}
	return svc
}

// 把corev1 svcs转成datacell
func (p *svc) toCells(std []corev1.Service) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = serviceCell(std[i])

	}
	return cells
}
