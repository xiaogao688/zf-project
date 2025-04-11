package main

import (
	"flag"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/shirou/gopsutil/process"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	"log"
	"os"
	"sync/atomic"
	"time"
)

type Asset struct {
	Namespace string
	Name      string
}

func main() {
	// 初始化命令行标志
	klog.InitFlags(nil)
	flag.Parse()

	// 资源监控
	go logResourceUsage("D:\\work\\zf-project\\learn\\client-go\\informer\\watch-pod-resource\\log")

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

	// 假设你希望针对 Pod 资源使用 30 秒的重同步周期，而其他资源使用默认的 10 分钟
	//customResync := map[metav1.Object]time.Duration{
	//	&v1.Pod{}: 30 * time.Second,
	//}
	//factory := informers.NewSharedInformerFactoryWithOptions(clientset, time.Minute,
	//	informers.WithTransform(func(obj interface{}) (interface{}, error) {
	//		// 例如：将对象转换为只包含 Name 与 Namespace 的结构
	//		pod := obj.(*v1.Pod)
	//		transformed := Asset{
	//			Namespace: pod.Namespace,
	//			Name:      pod.Name,
	//		}
	//		return &transformed, nil
	//	}),
	// 用来设置监听范围，可以在其中添加 labelSelector、fieldSelector 等过滤条件
	//informers.WithTweakListOptions(func(options *metav1.ListOptions) {
	//	// 只关注 label 为 "app=nginx" 的 Pod
	//	options.LabelSelector = "app=nginx"
	//}),
	//informers.WithCustomResyncConfig(customResync),
	//)

	// 创建 shared informer factory，设置 resync 时间（例如：30秒）
	factory := informers.NewSharedInformerFactory(clientset, 30*time.Second)

	podInformer := factory.Core().V1().Pods().Informer()

	// 定义一个原子变量来标识是否已经完成初始缓存同步，0 表示未完成，1 表示已完成
	var initialSyncComplete int32 = 0

	// 为 Pod informer 添加事件处理器
	podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			//pod := obj.(*v1.Pod)
			if atomic.LoadInt32(&initialSyncComplete) == 0 {
				// 此时为初始加载阶段
				//fmt.Printf("初始加载 Pod: %s/%s\n", pod.Namespace, pod.Name)
			} else {
				// 后续新增的 Pod
				//fmt.Printf("新增 Pod: %s/%s\n", pod.Namespace, pod.Name)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldPod := oldObj.(*v1.Pod)
			newPod := newObj.(*v1.Pod)
			//fmt.Printf("更新 Pod: %s/%s\n", newPod.Namespace, newPod.Name)
			// 比较 Pod Spec 的差异
			specDiff := cmp.Diff(oldPod.Spec, newPod.Spec)
			if specDiff != "" {
				//fmt.Printf("Pod 配置（Spec）更新: %s/%s\nDiff: %s\n", newPod.Namespace, newPod.Name, specDiff)
			}
			// 比较 Pod Status 的差异
			statusDiff := cmp.Diff(oldPod.Status.ContainerStatuses, newPod.Status.ContainerStatuses)
			if statusDiff != "" {
				//fmt.Printf("Pod 状态（Status）更新: %s/%s\nDiff: %s\n", newPod.Namespace, newPod.Name, statusDiff)
			}

			// 如果两个部分都没变化，也可以记录一下更新事件（通常不会发生）
			if specDiff == "" && statusDiff == "" {
				//fmt.Printf("Pod 更新，但无明显差异: %s/%s\n", newPod.Namespace, newPod.Name)
			}
		},
		DeleteFunc: func(obj interface{}) {
			pod := obj.(*v1.Pod)
			fmt.Printf("删除 Pod: %s/%s\n", pod.Namespace, pod.Name)
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

	// 同步完成后设置标志
	atomic.StoreInt32(&initialSyncComplete, 1)

	// 阻塞主线程，保持 informer 持续运行
	select {}
}

func logResourceUsage(filename string) {
	// 以追加模式打开文件（不存在则创建）
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Println("打开文件错误:", err)
		return
	}
	defer f.Close()

	// 获取当前进程对象
	p, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		log.Println("获取进程信息错误:", err)
		return
	}

	// 每隔5秒触发一次
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		timestamp := time.Now().Format(time.RFC3339)

		// 获取内存信息（RSS: 常驻内存, VMS: 虚拟内存）
		memInfo, err := p.MemoryInfo()
		if err != nil {
			log.Println("获取内存信息错误:", err)
			continue
		}

		// 获取CPU使用时间（user: 用户态, system: 系统态）
		cpuTimes, err := p.Times()
		if err != nil {
			log.Println("获取CPU信息错误:", err)
			continue
		}

		// 组织日志记录内容
		logStr := fmt.Sprintf("%s - Memory: RSS=%vKB, VMS=%vKB; CPU Times: user=%.2fs, system=%.2fs\n",
			timestamp, memInfo.RSS/1024, memInfo.VMS/1024, cpuTimes.User, cpuTimes.System)

		// 写入文件
		if _, err := f.WriteString(logStr); err != nil {
			log.Println("写入文件错误:", err)
		}
	}
}
