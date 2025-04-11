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
	factory := informers.NewSharedInformerFactory(clientset, 0)

	podInformer := factory.Core().V1().Pods().Informer()

	// 为 Pod informer 添加事件处理器
	podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			//fmt.Printf("新增 Pod: %s/%s\n", pod.Namespace, pod.Name)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			//oldPod := oldObj.(*v1.Pod)
			newPod := newObj.(*v1.Pod)

			//if GetPodState(oldPod.Status, "") != GetPodState(newPod.Status, "") {
			//	fmt.Printf("pod 状态变更: %s -> %s \n", GetPodState(oldPod.Status, ""), GetPodState(newPod.Status, fmt.Sprintf("更新 Pod: %s/%s ", newPod.Namespace, newPod.Name)))
			//}
			fmt.Printf("更新pod %s/%s 当前状态: %s \n", newPod.Namespace, newPod.Name, GetPodState(newPod.Status, ""))
		},
		DeleteFunc: func(obj interface{}) {
			pod := obj.(*v1.Pod)
			//fmt.Println("删除pod 当前状态: ", GetPodState(pod.Status, fmt.Sprintf("删除 Pod: %s/%s ", pod.Namespace, pod.Name)))
			fmt.Printf("删除pod %s/%s 当前状态: %s \n", pod.Namespace, pod.Name, GetPodState(pod.Status, ""))
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

	// 阻塞主线程，保持 informer 持续运行
	select {}
}

func GetPodState(status v1.PodStatus, msg string) string {
	for _, initContainerStatus := range status.InitContainerStatuses {
		if initContainerStatus.State.Waiting != nil {
			if msg != "" {
				fmt.Printf("%s msg: %s podState: %s \n", msg, initContainerStatus.State.Waiting.Message, string(status.Phase))
			}
			return initContainerStatus.State.Waiting.Reason
		}
		if initContainerStatus.State.Terminated != nil {
			if msg != "" {
				fmt.Printf("%s msg: %s podState: %s \n", msg, initContainerStatus.State.Terminated.Message, string(status.Phase))
			}
			// 如果init正常结束,Terminated虽不为空,但状态不应该展示他
			if initContainerStatus.State.Terminated.Reason != "Completed" {
				return initContainerStatus.State.Terminated.Reason
			}
		}
	}

	for _, containerStatus := range status.ContainerStatuses {
		if containerStatus.State.Waiting != nil {
			if msg != "" {
				fmt.Printf("%s msg: %s podState: %s \n", msg, containerStatus.State.Waiting.Message, string(status.Phase))
			}
			return containerStatus.State.Waiting.Reason
		}
		if containerStatus.State.Terminated != nil {
			if containerStatus.State.Terminated.Reason == "Error" && containerStatus.State.Terminated.Message == "" {
				return "Terminating"
			}
			if containerStatus.State.Terminated.Reason != "" {
				return containerStatus.State.Terminated.Reason
			}
		}
	}
	return string(status.Phase)
}
