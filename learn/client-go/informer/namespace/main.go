package main

import (
	"flag"
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	"os"
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

	// 创建 shared informer factory，设置 resync 时间（例如：30秒）
	factory := informers.NewSharedInformerFactory(clientset, 30*time.Second)

	podInformer := factory.Core().V1().Namespaces().Informer()

	// 为 Pod informer 添加事件处理器
	podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			_, ok := obj.(*v1.Namespace)
			if ok {

			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			_, ok := newObj.(*v1.Namespace)
			if ok {

			}
		},
		DeleteFunc: func(obj interface{}) {
			_, ok := obj.(*v1.Namespace)
			if !ok {
			}
		},
	})

	// 创建停止信号通道
	stopCh := make(chan struct{})
	defer close(stopCh)

	// 启动所有的 informer
	factory.Start(stopCh)

	// 等待所有 informer 的缓存同步完成
	if !cache.WaitForCacheSync(stopCh, podInformer.HasSynced) {
		klog.Fatalf("等待缓存同步超时")
	}

	obj, exists, err := podInformer.GetStore().GetByKey("dosec")
	if err != nil {
		fmt.Printf("获取 namespace 信息失败: %s\n", err.Error())
		os.Exit(1)
	}
	if !exists {
		fmt.Printf("未找到名称为 %s 的 namespace\n", "dosec")
		os.Exit(1)
	}
	ns := obj.(*v1.Namespace)
	fmt.Println(ns.UID, ns.Name)

	// 阻塞主线程，保持 informer 持续运行
	select {}
}
