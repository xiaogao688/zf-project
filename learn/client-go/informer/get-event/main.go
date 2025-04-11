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
	"time"
)

func main() {
	// 初始化命令行标志
	klog.InitFlags(nil)
	flag.Parse()

	// 读取 kubeconfig 文件路径（修改为你自己的 kubeconfig 路径）
	kubeconfig := flag.String("kubeconfig", "D:\\work\\zf-project\\learn\\client-go\\informer\\config", "Path to a kubeconfig file")
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

	selector := fields.OneTermEqualSelector("involvedObject.uid", "8f36c96d-2f2f-4fcc-b084-65199cf764f1")
	events, err := clientset.CoreV1().Events("default").List(ctx, metav1.ListOptions{
		FieldSelector: selector.String(),
	})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(events.Items)
}
