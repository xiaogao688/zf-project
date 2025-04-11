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

	podList := factory.Core().V1().Pods().Lister()
	pod, err := podList.Pods("dosec").Get("dosec-db-dao-tkjh7")

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

// GetPodState 获取pod状态,摘自k8s源码 pkg/printers/internalversion/printers.go -- printPod
func GetPodState(pod *v1.Pod) string {
	podPhase := pod.Status.Phase
	reason := string(podPhase)
	if pod.Status.Reason != "" {
		reason = pod.Status.Reason
	}
	// If the Pod carries {type:PodScheduled, reason:SchedulingGated}, set reason to 'SchedulingGated'.
	for _, condition := range pod.Status.Conditions {
		if condition.Type == v1.PodScheduled && condition.Reason == v1.PodReasonSchedulingGated {
			reason = v1.PodReasonSchedulingGated
		}
	}

	initializing := false
	for j := range pod.Status.InitContainerStatuses {
		container := pod.Status.InitContainerStatuses[j]
		switch {
		case container.State.Terminated != nil:
			// initialization is failed
			if len(container.State.Terminated.Reason) == 0 {
				if container.State.Terminated.Signal != 0 {
					reason = fmt.Sprintf("Init:Signal:%d", container.State.Terminated.Signal)
				} else {
					reason = fmt.Sprintf("Init:ExitCode:%d", container.State.Terminated.ExitCode)
				}
			} else {
				reason = "Init:" + container.State.Terminated.Reason
			}
			initializing = true
		case container.State.Waiting != nil && len(container.State.Waiting.Reason) > 0 && container.State.Waiting.Reason != "PodInitializing":
			reason = "Init:" + container.State.Waiting.Reason
			initializing = true
		default:
			reason = fmt.Sprintf("Init:%d/%d", j, len(pod.Spec.InitContainers))
			initializing = true
		}
		break
	}

	if !initializing || isPodInitializedConditionTrue(&pod.Status) {
		hasRunning := false
		for j := len(pod.Status.ContainerStatuses) - 1; j >= 0; j-- {
			container := pod.Status.ContainerStatuses[j]

			if container.State.Waiting != nil && container.State.Waiting.Reason != "" {
				reason = container.State.Waiting.Reason
			} else if container.State.Terminated != nil && container.State.Terminated.Reason != "" {
				reason = container.State.Terminated.Reason
			} else if container.State.Terminated != nil && container.State.Terminated.Reason == "" {
				if container.State.Terminated.Signal != 0 {
					reason = fmt.Sprintf("Signal:%d", container.State.Terminated.Signal)
				} else {
					reason = fmt.Sprintf("ExitCode:%d", container.State.Terminated.ExitCode)
				}
			} else if container.Ready && container.State.Running != nil {
				hasRunning = true
			}
		}

		// change pod status back to "Running" if there is at least one container still reporting as "Running" status
		if reason == "Completed" && hasRunning {
			if hasPodReadyCondition(pod.Status.Conditions) {
				reason = "Running"
			} else {
				reason = "NotReady"
			}
		}
	}

	if pod.DeletionTimestamp != nil && pod.Status.Reason == "NodeLost" {
		reason = "Unknown"
	} else if pod.DeletionTimestamp != nil && !IsPodPhaseTerminal(v1.PodPhase(podPhase)) {
		reason = "Terminating"
	}

	return reason
}
func IsPodPhaseTerminal(phase v1.PodPhase) bool {
	return phase == v1.PodFailed || phase == v1.PodSucceeded
}

func isPodInitializedConditionTrue(status *v1.PodStatus) bool {
	for _, condition := range status.Conditions {
		if condition.Type != v1.PodInitialized {
			continue
		}

		return condition.Status == v1.ConditionTrue
	}
	return false
}

func hasPodReadyCondition(conditions []v1.PodCondition) bool {
	for _, condition := range conditions {
		if condition.Type == v1.PodReady && condition.Status == v1.ConditionTrue {
			return true
		}
	}
	return false
}
