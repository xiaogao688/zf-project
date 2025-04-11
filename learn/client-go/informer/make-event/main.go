package main

import (
	"context"
	"fmt"
	"log"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// kubeconfig 路径，默认使用 ~/.kube/config
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		log.Fatalf("failed to load kubeconfig: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("failed to create clientset: %v", err)
	}

	now := time.Now()

	// 定义所有资源的 ObjectReference
	resources := []corev1.ObjectReference{
		{Kind: "Pod", Namespace: "default", Name: "example-pod"},
		{Kind: "PersistentVolume", Name: "example-pv"},
		{Kind: "PersistentVolumeClaim", Namespace: "default", Name: "example-pvc"},
		{Kind: "Secret", Namespace: "default", Name: "example-secret"},
		{Kind: "ConfigMap", Namespace: "default", Name: "example-configmap"},
		{Kind: "Endpoints", Namespace: "default", Name: "example-endpoints"},
		{Kind: "Namespace", Name: "default"},
		{Kind: "Service", Namespace: "default", Name: "example-service"},
		{Kind: "ServiceAccount", Namespace: "default", Name: "example-sa"},
		{Kind: "Role", Namespace: "default", Name: "example-role"},
		{Kind: "ClusterRole", Name: "example-clusterrole"},
		{Kind: "RoleBinding", Namespace: "default", Name: "example-rolebinding"},
		{Kind: "ClusterRoleBinding", Name: "example-crb"},
		{Kind: "ReplicationController", Namespace: "default", Name: "example-rc"},
		{Kind: "ReplicaSet", Namespace: "default", Name: "example-rs"},
		{Kind: "Deployment", Namespace: "default", Name: "example-deployment"},
		{Kind: "StatefulSet", Namespace: "default", Name: "example-sts"},
		{Kind: "DaemonSet", Namespace: "default", Name: "example-ds"},
		{Kind: "Job", Namespace: "default", Name: "example-job"},
		{Kind: "CronJob", Namespace: "default", Name: "example-cronjob"},
	}

	for _, ref := range resources {
		namespace := ref.Namespace
		if namespace == "" {
			namespace = "default" // fallback，如果资源是 cluster-scoped（如 PV、ClusterRole）
		}

		event := &corev1.Event{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%s-%s-%d", ref.Kind, ref.Name, now.UnixNano()),
				Namespace: namespace,
			},
			InvolvedObject: ref,
			Reason:         "ExampleEvent",
			Message:        fmt.Sprintf("This is a demo event for %s %s", ref.Kind, ref.Name),
			Type:           corev1.EventTypeNormal,
			Source:         corev1.EventSource{Component: "demo-event-generator"},
			FirstTimestamp: metav1.NewTime(now),
			LastTimestamp:  metav1.NewTime(now),
			Count:          1,
		}

		_, err := clientset.CoreV1().Events(namespace).Create(context.Background(), event, metav1.CreateOptions{})
		if err != nil {
			log.Printf("❌ Failed to create event for %s %s: %v", ref.Kind, ref.Name, err)
		} else {
			log.Printf("✅ Created event for %s %s", ref.Kind, ref.Name)
		}
	}

	log.Println("All events processed.")
}
