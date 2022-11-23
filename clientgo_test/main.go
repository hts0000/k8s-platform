package main

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	Kubeconfig := "C:\\Users\\init\\.kube\\config"
	//将kubeconfig转换成restclient,cconfig类型
	config, err := clientcmd.BuildConfigFromFlags("", Kubeconfig)
	if err != nil {
		panic(err)
	}
	//通过config创建clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	//通过cinetset获取pods，namespace 为默认
	//context.TODO() 控制上下文环境。比如请求超时
	//metav1.listOptions{} 可填写标签过滤配置。比如按标签
	pods, err := clientset.CoreV1().Pods("kube-system").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println(pods.Items)
}
