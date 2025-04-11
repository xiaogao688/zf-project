package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	batchV1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	rbacV1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	"time"
)

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

	kc := KubeClient{ClientSet: clientset}
	yaml, err := kc.GetAssetsYaml("role-1", "assets", "role", 100)
	if err != nil {
		klog.Fatal(err)
	}
	klog.Info(yaml)
}

type KubeClient struct {
	ClientSet *kubernetes.Clientset
}
type K8sAssetType string

const (
	K8sAssetTypePod                         K8sAssetType = "pod"
	K8sAssetTypePVC                         K8sAssetType = "pvc"
	K8sAssetTypeSecret                      K8sAssetType = "secret"
	K8sAssetTypeConfigMap                   K8sAssetType = "cm" // configMap
	K8sAssetTypeEndpoint                    K8sAssetType = "endpoint"
	K8sAssetTypeNamespace                   K8sAssetType = "namespace"
	K8sAssetTypeIngress                     K8sAssetType = "ingress"
	K8sAssetTypeService                     K8sAssetType = "service"
	K8sAssetTypeRoutes                      K8sAssetType = "routes"
	K8sAssetTypePV                          K8sAssetType = "pv"
	K8sAssetTypeServiceAccount              K8sAssetType = "service_account"
	K8sAssetTypeRole                        K8sAssetType = "role"
	K8sAssetTypeRoleBinding                 K8sAssetType = "role_binding"
	K8sAssetTypeClusterRole                 K8sAssetType = "cluster_role"
	K8sAssetTypeClusterRoleBinding          K8sAssetType = "cluster_role_binding"
	WorkloadKindReplicationControllerString              = "replication_controller" // ReplicationController
	WorkloadKindReplicaSetString                         = "replica_set"            // ReplicaSet
	WorkloadKindDaemonSetString                          = "daemon_set"             // DaemonSet
	WorkloadKindDeploymentString                         = "deployment"             // Deployment
	WorkloadKindStatefulSetString                        = "stateful_set"           // StatefulSet
	WorkloadKindJobString                                = "job"                    // Job
	WorkloadKindCronJobString                            = "cron_job"               // CronJob
)

func (k *KubeClient) GetAssetsYaml(assetsName, Namespace string, kind K8sAssetType, timeout int32) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	switch kind {
	case K8sAssetTypePod:
		assetInfo, err := k.ClientSet.CoreV1().Pods(Namespace).Get(ctx, assetsName, metav1.GetOptions{})
		if err != nil {
			return "", err
		}
		assetInfo.SetManagedFields(nil)
		s := runtime.NewScheme()
		if err := v1.AddToScheme(s); err != nil {
			return "", fmt.Errorf("failed to add v1 scheme: %w", err)
		}
		// **设置完整的 GVK（Group, Version, Kind）信息**
		assetInfo.GetObjectKind().SetGroupVersionKind(v1.SchemeGroupVersion.WithKind("Pod"))
		// **使用 YAML 序列化**
		serializer := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
		var buf bytes.Buffer
		// **写入 YAML 文件**
		err = serializer.Encode(assetInfo, &buf)
		if err != nil {
			return "", fmt.Errorf("failed to encode %s  to YAML: %w", kind, err)
		}
		return buf.String(), err
	case K8sAssetTypePVC:
		assetInfo, err := k.ClientSet.CoreV1().PersistentVolumeClaims(Namespace).Get(ctx, assetsName, metav1.GetOptions{})
		if err != nil {
			return "", err
		}
		assetInfo.SetManagedFields(nil)
		s := runtime.NewScheme()
		if err := v1.AddToScheme(s); err != nil {
			return "", fmt.Errorf("failed to add v1 scheme: %w", err)
		}
		// **设置完整的 GVK（Group, Version, Kind）信息**
		assetInfo.GetObjectKind().SetGroupVersionKind(v1.SchemeGroupVersion.WithKind("PersistentVolumeClaim"))
		// **使用 YAML 序列化**
		serializer := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
		var buf bytes.Buffer
		// **写入 YAML 文件**
		err = serializer.Encode(assetInfo, &buf)
		if err != nil {
			return "", fmt.Errorf("failed to encode %s to YAML: %w", kind, err)
		}
		return buf.String(), err
	case K8sAssetTypeSecret:
		assetInfo, err := k.ClientSet.CoreV1().Secrets(Namespace).Get(ctx, assetsName, metav1.GetOptions{})
		if err != nil {
			return "", err
		}
		assetInfo.SetManagedFields(nil)
		s := runtime.NewScheme()
		if err := v1.AddToScheme(s); err != nil {
			return "", fmt.Errorf("failed to add v1 scheme: %w", err)
		}
		// **设置完整的 GVK（Group, Version, Kind）信息**
		assetInfo.GetObjectKind().SetGroupVersionKind(v1.SchemeGroupVersion.WithKind("Secret"))
		// **使用 YAML 序列化**
		serializer := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
		var buf bytes.Buffer
		// **写入 YAML 文件**
		err = serializer.Encode(assetInfo, &buf)
		if err != nil {
			return "", fmt.Errorf("failed to encode %s to YAML: %w", kind, err)
		}
		return buf.String(), err
	case K8sAssetTypeConfigMap:
		assetInfo, err := k.ClientSet.CoreV1().ConfigMaps(Namespace).Get(ctx, assetsName, metav1.GetOptions{})
		if err != nil {
			return "", err
		}
		assetInfo.SetManagedFields(nil)
		s := runtime.NewScheme()
		if err := v1.AddToScheme(s); err != nil {
			return "", fmt.Errorf("failed to add v1 scheme: %w", err)
		}
		// **设置完整的 GVK（Group, Version, Kind）信息**
		assetInfo.GetObjectKind().SetGroupVersionKind(v1.SchemeGroupVersion.WithKind("ConfigMap"))
		// **使用 YAML 序列化**
		serializer := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
		var buf bytes.Buffer
		// **写入 YAML 文件**
		err = serializer.Encode(assetInfo, &buf)
		if err != nil {
			return "", fmt.Errorf("failed to encode %s to YAML: %w", kind, err)
		}
		return buf.String(), err
	case WorkloadKindReplicationControllerString:
		assetInfo, err := k.ClientSet.CoreV1().ReplicationControllers(Namespace).Get(ctx, assetsName, metav1.GetOptions{})
		if err != nil {
			return "", err
		}
		assetInfo.SetManagedFields(nil)
		s := runtime.NewScheme()
		if err := v1.AddToScheme(s); err != nil {
			return "", fmt.Errorf("failed to add v1 scheme: %w", err)
		}
		// **设置完整的 GVK（Group, Version, Kind）信息**
		assetInfo.GetObjectKind().SetGroupVersionKind(v1.SchemeGroupVersion.WithKind("ReplicationControllers"))
		// **使用 YAML 序列化**
		serializer := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
		var buf bytes.Buffer
		// **写入 YAML 文件**
		err = serializer.Encode(assetInfo, &buf)
		if err != nil {
			return "", fmt.Errorf("failed to encode %s to YAML: %w", kind, err)
		}
		return buf.String(), err
	case WorkloadKindReplicaSetString:
		assetInfo, err := k.ClientSet.AppsV1().ReplicaSets(Namespace).Get(ctx, assetsName, metav1.GetOptions{})
		if err != nil {
			return "", err
		}
		assetInfo.SetManagedFields(nil)
		s := runtime.NewScheme()
		if err := appsv1.AddToScheme(s); err != nil {
			return "", fmt.Errorf("failed to add v1 scheme: %w", err)
		}
		// **设置完整的 GVK（Group, Version, Kind）信息**
		assetInfo.GetObjectKind().SetGroupVersionKind(appsv1.SchemeGroupVersion.WithKind("ReplicaSets"))
		// **使用 YAML 序列化**
		serializer := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
		var buf bytes.Buffer
		// **写入 YAML 文件**
		err = serializer.Encode(assetInfo, &buf)
		if err != nil {
			return "", fmt.Errorf("failed to encode %s to YAML: %w", kind, err)
		}
		return buf.String(), err
	case WorkloadKindDaemonSetString:
		assetInfo, err := k.ClientSet.AppsV1().DaemonSets(Namespace).Get(ctx, assetsName, metav1.GetOptions{})
		if err != nil {
			return "", err
		}
		assetInfo.SetManagedFields(nil)
		s := runtime.NewScheme()
		if err := appsv1.AddToScheme(s); err != nil {
			return "", fmt.Errorf("failed to add v1 scheme: %w", err)
		}
		// **设置完整的 GVK（Group, Version, Kind）信息**
		assetInfo.GetObjectKind().SetGroupVersionKind(appsv1.SchemeGroupVersion.WithKind("DaemonSets"))
		// **使用 YAML 序列化**
		serializer := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
		var buf bytes.Buffer
		// **写入 YAML 文件**
		err = serializer.Encode(assetInfo, &buf)
		if err != nil {
			return "", fmt.Errorf("failed to encode %s to YAML: %w", kind, err)
		}
		return buf.String(), err
	case WorkloadKindDeploymentString:
		assetInfo, err := k.ClientSet.AppsV1().Deployments(Namespace).Get(ctx, assetsName, metav1.GetOptions{})
		if err != nil {
			return "", err
		}
		assetInfo.SetManagedFields(nil)
		s := runtime.NewScheme()
		if err := appsv1.AddToScheme(s); err != nil {
			return "", fmt.Errorf("failed to add v1 scheme: %w", err)
		}
		// **设置完整的 GVK（Group, Version, Kind）信息**
		assetInfo.GetObjectKind().SetGroupVersionKind(appsv1.SchemeGroupVersion.WithKind("Deployments"))
		// **使用 YAML 序列化**
		serializer := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
		var buf bytes.Buffer
		// **写入 YAML 文件**
		err = serializer.Encode(assetInfo, &buf)
		if err != nil {
			return "", fmt.Errorf("failed to encode %s to YAML: %w", kind, err)
		}
		return buf.String(), err
	case WorkloadKindStatefulSetString:
		assetInfo, err := k.ClientSet.AppsV1().StatefulSets(Namespace).Get(ctx, assetsName, metav1.GetOptions{})
		if err != nil {
			return "", err
		}
		assetInfo.SetManagedFields(nil)
		s := runtime.NewScheme()
		if err := appsv1.AddToScheme(s); err != nil {
			return "", fmt.Errorf("failed to add v1 scheme: %w", err)
		}
		// **设置完整的 GVK（Group, Version, Kind）信息**
		assetInfo.GetObjectKind().SetGroupVersionKind(appsv1.SchemeGroupVersion.WithKind("StatefulSets"))
		// **使用 YAML 序列化**
		serializer := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
		var buf bytes.Buffer
		// **写入 YAML 文件**
		err = serializer.Encode(assetInfo, &buf)
		if err != nil {
			return "", fmt.Errorf("failed to encode %s to YAML: %w", kind, err)
		}
		return buf.String(), err
	case WorkloadKindJobString:
		assetInfo, err := k.ClientSet.BatchV1().Jobs(Namespace).Get(ctx, assetsName, metav1.GetOptions{})
		if err != nil {
			return "", err
		}
		assetInfo.SetManagedFields(nil)
		s := runtime.NewScheme()
		if err := batchV1.AddToScheme(s); err != nil {
			return "", fmt.Errorf("failed to add v1 scheme: %w", err)
		}
		// **设置完整的 GVK（Group, Version, Kind）信息**
		assetInfo.GetObjectKind().SetGroupVersionKind(batchV1.SchemeGroupVersion.WithKind("Jobs"))
		// **使用 YAML 序列化**
		serializer := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
		var buf bytes.Buffer
		// **写入 YAML 文件**
		err = serializer.Encode(assetInfo, &buf)
		if err != nil {
			return "", fmt.Errorf("failed to encode %s to YAML: %w", kind, err)
		}
		return buf.String(), err
	case WorkloadKindCronJobString:
		assetInfo, err := k.ClientSet.BatchV1().CronJobs(Namespace).Get(ctx, assetsName, metav1.GetOptions{})
		if err != nil {
			return "", err
		}
		assetInfo.SetManagedFields(nil)
		s := runtime.NewScheme()
		if err := v1.AddToScheme(s); err != nil {
			return "", fmt.Errorf("failed to add v1 scheme: %w", err)
		}
		// **设置完整的 GVK（Group, Version, Kind）信息**
		assetInfo.GetObjectKind().SetGroupVersionKind(v1.SchemeGroupVersion.WithKind("CronJobs"))
		// **使用 YAML 序列化**
		serializer := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
		var buf bytes.Buffer
		// **写入 YAML 文件**
		err = serializer.Encode(assetInfo, &buf)
		if err != nil {
			return "", fmt.Errorf("failed to encode %s to YAML: %w", kind, err)
		}
		return buf.String(), err
	case K8sAssetTypeEndpoint:
		assetInfo, err := k.ClientSet.CoreV1().Endpoints(Namespace).Get(ctx, assetsName, metav1.GetOptions{})
		if err != nil {
			return "", err
		}
		assetInfo.SetManagedFields(nil)
		s := runtime.NewScheme()
		if err := v1.AddToScheme(s); err != nil {
			return "", fmt.Errorf("failed to add v1 scheme: %w", err)
		}
		// **设置完整的 GVK（Group, Version, Kind）信息**
		assetInfo.GetObjectKind().SetGroupVersionKind(v1.SchemeGroupVersion.WithKind("Endpoints"))
		// **使用 YAML 序列化**
		serializer := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
		var buf bytes.Buffer
		// **写入 YAML 文件**
		err = serializer.Encode(assetInfo, &buf)
		if err != nil {
			return "", fmt.Errorf("failed to encode %s to YAML: %w", kind, err)
		}
		return buf.String(), err
	case K8sAssetTypeNamespace:
		assetInfo, err := k.ClientSet.CoreV1().Namespaces().Get(ctx, assetsName, metav1.GetOptions{})
		if err != nil {
			return "", err
		}
		assetInfo.SetManagedFields(nil)
		s := runtime.NewScheme()
		if err := v1.AddToScheme(s); err != nil {
			return "", fmt.Errorf("failed to add v1 scheme: %w", err)
		}
		// **设置完整的 GVK（Group, Version, Kind）信息**
		assetInfo.GetObjectKind().SetGroupVersionKind(v1.SchemeGroupVersion.WithKind("Namespace"))
		// **使用 YAML 序列化**
		serializer := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
		var buf bytes.Buffer
		// **写入 YAML 文件**
		err = serializer.Encode(assetInfo, &buf)
		if err != nil {
			return "", fmt.Errorf("failed to encode %s to YAML: %w", kind, err)
		}
		return buf.String(), err
	case K8sAssetTypeIngress:
		//return k.IngressYaml(ctx, Namespace, assetsName)
	case K8sAssetTypeService:
		assetInfo, err := k.ClientSet.CoreV1().Services(Namespace).Get(ctx, assetsName, metav1.GetOptions{})
		if err != nil {
			return "", err
		}
		assetInfo.SetManagedFields(nil)
		s := runtime.NewScheme()
		if err := v1.AddToScheme(s); err != nil {
			return "", fmt.Errorf("failed to add v1 scheme: %w", err)
		}
		// **设置完整的 GVK（Group, Version, Kind）信息**
		assetInfo.GetObjectKind().SetGroupVersionKind(v1.SchemeGroupVersion.WithKind("Service"))
		// **使用 YAML 序列化**
		serializer := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
		var buf bytes.Buffer
		// **写入 YAML 文件**
		err = serializer.Encode(assetInfo, &buf)
		if err != nil {
			return "", fmt.Errorf("failed to encode %s to YAML: %w", kind, err)
		}
		return buf.String(), err
	case K8sAssetTypeRoutes:
		// 暂不支持
	case K8sAssetTypePV:
		assetInfo, err := k.ClientSet.CoreV1().PersistentVolumes().Get(ctx, assetsName, metav1.GetOptions{})
		if err != nil {
			return "", err
		}
		assetInfo.SetManagedFields(nil)
		s := runtime.NewScheme()
		if err := v1.AddToScheme(s); err != nil {
			return "", fmt.Errorf("failed to add v1 scheme: %w", err)
		}
		// **设置完整的 GVK（Group, Version, Kind）信息**
		assetInfo.GetObjectKind().SetGroupVersionKind(v1.SchemeGroupVersion.WithKind("PersistentVolume"))
		// **使用 YAML 序列化**
		serializer := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
		var buf bytes.Buffer
		// **写入 YAML 文件**
		err = serializer.Encode(assetInfo, &buf)
		if err != nil {
			return "", fmt.Errorf("failed to encode %s to YAML: %w", kind, err)
		}
		return buf.String(), err
	case K8sAssetTypeServiceAccount:
		assetInfo, err := k.ClientSet.CoreV1().ServiceAccounts(Namespace).Get(ctx, assetsName, metav1.GetOptions{})
		if err != nil {
			return "", err
		}
		assetInfo.SetManagedFields(nil)
		s := runtime.NewScheme()
		if err := v1.AddToScheme(s); err != nil {
			return "", fmt.Errorf("failed to add v1 scheme: %w", err)
		}
		// **设置完整的 GVK（Group, Version, Kind）信息**
		assetInfo.GetObjectKind().SetGroupVersionKind(v1.SchemeGroupVersion.WithKind("ServiceAccount"))
		// **使用 YAML 序列化**
		serializer := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
		var buf bytes.Buffer
		// **写入 YAML 文件**
		err = serializer.Encode(assetInfo, &buf)
		if err != nil {
			return "", fmt.Errorf("failed to encode %s to YAML: %w", kind, err)
		}
		return buf.String(), err
	case K8sAssetTypeRole:
		assetInfo, err := k.ClientSet.RbacV1().Roles(Namespace).Get(ctx, assetsName, metav1.GetOptions{})
		if err != nil {
			return "", err
		}
		assetInfo.SetManagedFields(nil)
		s := runtime.NewScheme()
		if err := rbacV1.AddToScheme(s); err != nil {
			return "", fmt.Errorf("failed to add v1 scheme: %w", err)
		}
		// **设置完整的 GVK（Group, Version, Kind）信息**
		assetInfo.GetObjectKind().SetGroupVersionKind(rbacV1.SchemeGroupVersion.WithKind("Role"))
		// **使用 YAML 序列化**
		serializer := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
		var buf bytes.Buffer
		// **写入 YAML 文件**
		err = serializer.Encode(assetInfo, &buf)
		if err != nil {
			return "", fmt.Errorf("failed to encode %s to YAML: %w", kind, err)
		}
		return buf.String(), err
	case K8sAssetTypeRoleBinding:
		assetInfo, err := k.ClientSet.RbacV1().RoleBindings(Namespace).Get(ctx, assetsName, metav1.GetOptions{})
		if err != nil {
			return "", err
		}
		assetInfo.SetManagedFields(nil)
		s := runtime.NewScheme()
		if err := rbacV1.AddToScheme(s); err != nil {
			return "", fmt.Errorf("failed to add v1 scheme: %w", err)
		}
		// **设置完整的 GVK（Group, Version, Kind）信息**
		assetInfo.GetObjectKind().SetGroupVersionKind(rbacV1.SchemeGroupVersion.WithKind("RoleBinding"))
		// **使用 YAML 序列化**
		serializer := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
		var buf bytes.Buffer
		// **写入 YAML 文件**
		err = serializer.Encode(assetInfo, &buf)
		if err != nil {
			return "", fmt.Errorf("failed to encode %s to YAML: %w", kind, err)
		}
		return buf.String(), err
	case K8sAssetTypeClusterRole:
		assetInfo, err := k.ClientSet.RbacV1().ClusterRoles().Get(ctx, assetsName, metav1.GetOptions{})
		if err != nil {
			return "", err
		}
		assetInfo.SetManagedFields(nil)
		s := runtime.NewScheme()
		if err := rbacV1.AddToScheme(s); err != nil {
			return "", fmt.Errorf("failed to add v1 scheme: %w", err)
		}
		// **设置完整的 GVK（Group, Version, Kind）信息**
		assetInfo.GetObjectKind().SetGroupVersionKind(rbacV1.SchemeGroupVersion.WithKind("ClusterRole"))
		// **使用 YAML 序列化**
		serializer := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
		var buf bytes.Buffer
		// **写入 YAML 文件**
		err = serializer.Encode(assetInfo, &buf)
		if err != nil {
			return "", fmt.Errorf("failed to encode %s to YAML: %w", kind, err)
		}
		return buf.String(), err
	case K8sAssetTypeClusterRoleBinding:
		assetInfo, err := k.ClientSet.RbacV1().ClusterRoleBindings().Get(ctx, assetsName, metav1.GetOptions{})
		if err != nil {
			return "", err
		}
		assetInfo.SetManagedFields(nil)
		s := runtime.NewScheme()
		if err := rbacV1.AddToScheme(s); err != nil {
			return "", fmt.Errorf("failed to add v1 scheme: %w", err)
		}
		// **设置完整的 GVK（Group, Version, Kind）信息**
		assetInfo.GetObjectKind().SetGroupVersionKind(rbacV1.SchemeGroupVersion.WithKind("ClusterRoleBinding"))
		// **使用 YAML 序列化**
		serializer := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
		var buf bytes.Buffer
		// **写入 YAML 文件**
		err = serializer.Encode(assetInfo, &buf)
		if err != nil {
			return "", fmt.Errorf("failed to encode %s to YAML: %w", kind, err)
		}
		return buf.String(), err
	default:
		return "", errors.New("type not supported ")

	}
	return "", nil
}
