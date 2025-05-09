package main

import (
	"context"
	"flag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	"log"
	"strings"
	"time"
)

func main() {
	// 初始化命令行标志
	klog.InitFlags(nil)
	flag.Parse()

	// 读取 kubeconfig 文件路径（修改为你自己的 kubeconfig 路径）
	kubeconfig := flag.String("kubeconfig", "D:\\work\\zf-project\\learn\\client-go\\informer\\config4-137", "Path to a kubeconfig file")
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		klog.Fatalf("构建 kubeconfig 失败: %v", err)
	}

	// 创建 Kubernetes 客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Fatalf("创建 Kubernetes 客户端失败: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()

	selector := fields.OneTermEqualSelector("involvedObject.uid", "ca9d980e-0267-4e95-95d4-020e8fe1d4d2")
	//selector := fields.OneTermEqualSelector("regarding.uid", "7dd5a644-1a3e-4d9b-b2aa-1feeab4a38ff")
	events, err := clientset.EventsV1beta1().Events("assets").List(ctx, metav1.ListOptions{
		FieldSelector: selector.String(),
	})
	if err != nil && strings.Contains(err.Error(), "the server could not find the requested resource") {
		log.Fatalln(err)
	}
	log.Println(events.Items)
}
