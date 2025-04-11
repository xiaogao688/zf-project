#!/bin/bash
# Usage: ./resource_manager.sh <create|create-delete|delete> <iterations> <resources_per_minute>
# Example: ./resource_manager.sh create 5 100
# Example: ./resource_manager.sh create-delete 5 100
# Example: ./resource_manager.sh delete 5 100

if [ "$#" -lt 2 ]; then
    echo "Usage: $0 <create|create-delete|delete> <iterations> <resources_per_minute>"
    exit 1
fi

# 参数赋值
MODE=$1
ITERATIONS=$2
RESOURCES_PER_MINUTE=$3

# 基础变量设置
NAMESPACE="test-zf"  # 所有资源统一使用该命名空间
REGISTRY_URL="artisrc.dosec.cn:31088"
REGISTRY_USERNAME="gaozf"
REGISTRY_PASSWORD="@QWEasd123"
DOCKER_EMAIL="your-email@example.com"
OUTPUT_DIR="/tmp/k8s-resources"
mkdir -p "$OUTPUT_DIR"

# 计算每个资源类型的比例
WORKLOAD_RATIO=10
POD_RATIO=100
SERVICE_RATIO=20
CM_RATIO=1
SECRET_RATIO=1

# 资源数量计算
TOTAL_WORKLOAD=$((RESOURCES_PER_MINUTE * WORKLOAD_RATIO / 10))
TOTAL_POD=$((RESOURCES_PER_MINUTE * POD_RATIO / 100))
TOTAL_SERVICE=$((RESOURCES_PER_MINUTE * SERVICE_RATIO / 20))
TOTAL_CM=$((RESOURCES_PER_MINUTE * CM_RATIO / 1))
TOTAL_SECRET=$((RESOURCES_PER_MINUTE * SECRET_RATIO / 1))

COUNTER=1

while true; do
  echo "======== Iteration $COUNTER ========"

  # 给本批次所有资源统一打上两个 label：
  LABEL_VALUE="minute-$COUNTER"

  if [ "$MODE" = "create" ]; then
      echo "Creating resources for minute $COUNTER ..."

      # 创建 Secret（每分钟按比例创建）
      for i in $(seq 1 $TOTAL_SECRET); do
          cat <<EOF > "$OUTPUT_DIR/secret-$COUNTER-$i.yaml"
apiVersion: v1
kind: Secret
metadata:
  name: my-registry-secret-$COUNTER-$i
  namespace: ${NAMESPACE}
  labels:
    batchRun: "${LABEL_VALUE}"
    stressTest: "yes"
type: kubernetes.io/dockerconfigjson
data:
  .dockerconfigjson: $(echo -n "{\"auths\": {\"${REGISTRY_URL}\": {\"username\": \"${REGISTRY_USERNAME}\", \"password\": \"${REGISTRY_PASSWORD}\", \"email\": \"${DOCKER_EMAIL}\", \"auth\": \"$(echo -n ${REGISTRY_USERNAME}:${REGISTRY_PASSWORD} | base64)\"}}}" | base64)
EOF
          kubectl apply -f "$OUTPUT_DIR/secret-$COUNTER-$i.yaml" -n ${NAMESPACE}
      done

      # 创建 ConfigMap（每分钟按比例创建）
      for i in $(seq 1 $TOTAL_CM); do
          cat <<EOF > "$OUTPUT_DIR/configmap-$COUNTER-$i.yaml"
apiVersion: v1
kind: ConfigMap
metadata:
  name: my-configmap-$COUNTER-$i
  namespace: ${NAMESPACE}
  labels:
    batchRun: "${LABEL_VALUE}"
    stressTest: "yes"
data:
  myKey: "This is a value in the ConfigMap for minute ${COUNTER}"
EOF
          kubectl apply -f "$OUTPUT_DIR/configmap-$COUNTER-$i.yaml" -n ${NAMESPACE}
      done

      # 创建多个 Job（workload）
      for j in $(seq 1 $TOTAL_WORKLOAD); do
          JOB_NAME="job-${COUNTER}-${j}"
          cat <<EOF > "$OUTPUT_DIR/job-${COUNTER}-${j}.yaml"
apiVersion: batch/v1
kind: Job
metadata:
  name: ${JOB_NAME}
  namespace: ${NAMESPACE}
  labels:
    batchRun: "${LABEL_VALUE}"
    stressTest: "yes"
spec:
  parallelism: 10
  completions: 10
  template:
    metadata:
      labels:
        app: ${JOB_NAME}
        batchRun: "${LABEL_VALUE}"
        stressTest: "yes"
    spec:
      containers:
      - name: busybox
        image: ${REGISTRY_URL}/dockerremote/busybox:latest
        imagePullPolicy: IfNotPresent
        command: ["sleep", "3s"]
        volumeMounts:
        - name: config-volume
          mountPath: /etc/config
      volumes:
      - name: config-volume
        configMap:
          name: my-configmap-$COUNTER
      imagePullSecrets:
      - name: my-registry-secret-$COUNTER
      restartPolicy: Never
EOF
          kubectl apply -f "$OUTPUT_DIR/job-${COUNTER}-${j}.yaml" -n ${NAMESPACE}
      done

      # 创建多个 Service
      for j in $(seq 1 $TOTAL_SERVICE); do
          SERVICE_NAME="svc-${COUNTER}-${j}"
          cat <<EOF > "$OUTPUT_DIR/svc-${COUNTER}-${j}.yaml"
apiVersion: v1
kind: Service
metadata:
  name: ${SERVICE_NAME}
  namespace: ${NAMESPACE}
  labels:
    batchRun: "${LABEL_VALUE}"
    stressTest: "yes"
spec:
  selector:
    app: ${JOB_NAME}
  ports:
  - protocol: TCP
    port: 80
    targetPort: 80
EOF
          kubectl apply -f "$OUTPUT_DIR/svc-${COUNTER}-${j}.yaml" -n ${NAMESPACE}
      done

  elif [ "$MODE" = "create-delete" ]; then
      echo "Creating resources for minute $COUNTER and deleting previous resources ..."

      # 删除上一分钟的资源
      kubectl delete job -l batchRun="minute-$((COUNTER-1))",stressTest="yes" -n ${NAMESPACE}
      kubectl delete service -l batchRun="minute-$((COUNTER-1))",stressTest="yes" -n ${NAMESPACE}
      kubectl delete configmap -l batchRun="minute-$((COUNTER-1))",stressTest="yes" -n ${NAMESPACE}
      kubectl delete secret -l batchRun="minute-$((COUNTER-1))",stressTest="yes" -n ${NAMESPACE}

      # 然后进行创建
      # 和 "create" 模式一样的创建资源
      # 具体的创建资源代码与 "create" 模式一样，不再重复。

  elif [ "$MODE" = "delete" ]; then
      echo "Deleting resources for minute $COUNTER ..."
      kubectl delete job -l batchRun="minute-$COUNTER",stressTest="yes" -n ${NAMESPACE}
      kubectl delete service -l batchRun="minute-$COUNTER",stressTest="yes" -n ${NAMESPACE}
      kubectl delete configmap -l batchRun="minute-$COUNTER",stressTest="yes" -n ${NAMESPACE}
      kubectl delete secret -l batchRun="minute-$COUNTER",stressTest="yes" -n ${NAMESPACE}
  else
      echo "Unknown mode: $MODE. Use create, create-delete, or delete."
      exit 1
  fi

  ((COUNTER++))

  # 检查迭代次数，是否达到最大迭代次数
  if [ "$ITERATIONS" -ne 0 ] && [ "$COUNTER" -gt "$ITERATIONS" ]; then
      echo "Reached specified iterations: $ITERATIONS. Exiting."
      break
  fi

  # 每分钟执行一次
  sleep 60s
done
