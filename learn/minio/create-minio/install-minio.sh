

# 在每个节点创建持久化目录
mkdir -p /data/gzf/minio

# 给MinIO 运行的节点都需要打标签，确定在哪些节点运行
kubectl label node dosec-9.201 minio=true
kubectl label node dosec-9.202 minio=true
kubectl label node dosec-9.203 minio=true
kubectl label node dosec-9.204 minio=true

# 创建命名空间
kubectl create ns minio-cluster

kubectl apply -f pv.yaml

kubectl apply -f minio.yaml

