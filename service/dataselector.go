package service

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"sort"
	"strings"
	"time"
)

// DataSelector 结构体用于排序过滤和分页
type DataSelector struct {
	GenericDataList []DataCell
	DataSelectQuery *DataSelect
}

// DataCell 接口用于各种资源的类型转换， 排序 过滤 分页 统一对dataCell进行处理
type DataCell interface {
	GetCreation() time.Time
	GetName() string
}

// DataSelect 定义过滤和分页的属性
type DataSelect struct {
	FilterQuery   *Filter
	PaginateQuery *Paginate
}
type Filter struct {
	Name string
}

type Paginate struct {
	Limit int
	Page  int
}

//排序
//实现自定义排序 需要重写len swap less 方法
//len方法用于获取数组长度

func (d *DataSelector) Len() int {
	return len(d.GenericDataList)
}

// Swap 方法用于在len方法比较结果后。定义排序规则
func (d *DataSelector) Swap(i, j int) {
	d.GenericDataList[i], d.GenericDataList[j] = d.GenericDataList[j], d.GenericDataList[i]
}

// Less 方法用于定义数组中元素大小的比较方式
func (d *DataSelector) Less(i, j int) bool {
	a := d.GenericDataList[i].GetCreation()
	b := d.GenericDataList[j].GetCreation()
	return b.Before(a)
}

// 重写以上三个方法后，用sort.sort 方法进行排序
func (d *DataSelector) Sort() *DataSelector {
	sort.Sort(d)
	return d
}

// 过滤
// 比较元素中是否存在filterName相匹配的元素，若匹配则返回
func (d *DataSelector) Filter() *DataSelector {
	//若name为空则返回所有
	if d.DataSelectQuery.FilterQuery.Name == "" {
		return d
	}
	//若name 传参不为空则返回切片中包含name所有元素
	fileredList := make([]DataCell, 0)
	for _, value := range d.GenericDataList {
		matched := true
		objName := value.GetName()
		if !strings.Contains(objName, d.DataSelectQuery.FilterQuery.Name) {
			matched = false
			continue
		}
		if matched {
			fileredList = append(fileredList, value)
		}
	}
	d.GenericDataList = fileredList
	return d
}

// 分页
// 根据limit和page的穿插，返回数据
func (d *DataSelector) Paginate() *DataSelector {
	limit := d.DataSelectQuery.PaginateQuery.Limit
	page := d.DataSelectQuery.PaginateQuery.Page
	//验证参数是否合法，若参数不合法，则返回所有
	if limit <= 0 || page <= 0 {
		return d
	}
	startIndex := limit * (page - 1)
	endIndex := limit * page
	//处理最后一页
	if len(d.GenericDataList) < endIndex {
		endIndex = len(d.GenericDataList)
	}
	d.GenericDataList = d.GenericDataList[startIndex:endIndex]
	return d
}

// 定义podcell类型 实现datacell接口，用于类型转换
type podCell corev1.Pod

func (p podCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}

func (p podCell) GetName() string {
	return p.Name
}

// deployment
type deploymentCell appsv1.Deployment

func (d deploymentCell) GetCreation() time.Time {
	return d.CreationTimestamp.Time
}
func (d deploymentCell) GetName() string {
	return d.Name
}
