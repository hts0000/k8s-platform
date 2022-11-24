package service

import (
	"context"
	"errors"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
)

var Secret secret

type secret struct{}

// 定义secrets的返回内容 items是secret列表吗total为secret元素总数
type SecretsResp struct {
	Item  []corev1.Secret `json:"items"`
	Total int             `json:"total"`
}

// 获取secret列表
func (p *secret) GetSecret(filterName, namespace string, limit, page int) (secretsResp *SecretsResp, err error) {
	//通过clientset获取secrets完整列表
	secretList, err := K8s.ClientSet.CoreV1().Secrets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error("获取secret列表失败", err)
		return nil, errors.New("获取secret列表失败" + err.Error())
	}
	//实例化DataSelector对象
	selectableData := &DataSelector{
		GenericDataList: p.toCells(secretList.Items),
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
	//再将datacell切片数据转成原生secret切片
	secrets := p.fromCells(data.GenericDataList)
	//返回
	return &SecretsResp{
		Item:  secrets,
		Total: total,
	}, nil
}

// 获取secret详情
func (p *secret) GetSecretDetail(secretName, namespace string) (secret *corev1.Secret, err error) {
	secret, err = K8s.ClientSet.CoreV1().Secrets(namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取Secret详情失败" + err.Error())
		return nil, errors.New("获取Secret详情失败" + err.Error())
	}
	return secret, nil
}

// 删除secret
func (p *secret) DeleteSecret(secretName, namespace string) (err error) {
	err = K8s.ClientSet.CoreV1().Secrets(namespace).Delete(context.TODO(), secretName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error("删除Secret失败" + err.Error())
		return errors.New("获取Secret详情失败" + err.Error())
	}
	return nil
}

// 更新secret
func (p *secret) UpdateSecret(namespace, content string) (err error) {
	//将content反序列化成为secret对象
	var secret = &corev1.Secret{}
	if err = json.Unmarshal([]byte(content), secret); err != nil {
		logger.Error("Content反序列化失败", err)
		return errors.New("Content反序列化失败" + err.Error())
	}
	//更新secret
	_, err = K8s.ClientSet.CoreV1().Secrets(namespace).Update(context.TODO(), secret, metav1.UpdateOptions{})
	if err != nil {
		logger.Error("更新Secret失败" + err.Error())
		return errors.New("更新Secret失败" + err.Error())
	}
	return nil
}

// 把secretCell转成appsv1 secret
func (p *secret) fromCells(cells []DataCell) []corev1.Secret {
	secrets := make([]corev1.Secret, len(cells))
	for i := range cells {
		secrets[i] = corev1.Secret(cells[i].(secretCell))
	}
	return secrets
}

// 把appsv1 secrets转成datacell
func (p *secret) toCells(std []corev1.Secret) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = secretCell(std[i])

	}
	return cells
}
