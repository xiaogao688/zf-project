package main

import (
	"context"
	"flag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	"log"
)

type Asset struct {
	Namespace string
	Name      string
}

func main() {
	// 初始化命令行标志
	klog.InitFlags(nil)
	flag.Parse()

	// 读取 kubeconfig 文件路径（修改为你自己的 kubeconfig 路径）
	kubeconfig := flag.String("kubeconfig", "D:\\work\\zf-project\\learn\\client-go\\informer\\71config", "Path to a kubeconfig file")
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		klog.Fatalf("构建 kubeconfig 失败: %v", err)
	}

	// 创建 Kubernetes 客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Fatalf("创建 Kubernetes 客户端失败: %v", err)
	}

	allCJ, err := clientset.BatchV1().CronJobs("assets").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Fatalln(err)
		return
	}
	log.Println(allCJ)

	allCJ2, err := clientset.BatchV1beta1().CronJobs("assets").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Fatalln(err)
		return
	}
	log.Println(allCJ2)

	// 阻塞主线程，保持 informer 持续运行
	select {}
}
