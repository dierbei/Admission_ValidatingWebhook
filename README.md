## Admission ValidatingWebhook
```shell
# 安装 cert maneger
kubectl apply -f deploy/cert-manager-1.5.3.yaml

# 创建证书 secret 
kubectl apply -f deploy/tls-secret.yaml

# 创建 configuration
kubectl apply -f deploy/validatingWebhookConfiguration.yaml

# 部署应用
kubectl apply -f deploy/deployment.yaml

# 卸载全部
make clear
```

## 构建镜像
```shell
docker build --platform=linux/amd64 -t dierbei/admission_validating:202304271401 .
```

## 命名空间添加 label
```shell
kubectl label namespace admission nginx=true
```

## 测试
```shell
kubectl run test1 --image=nginx:1.18-alpine --namespace=admission
kubectl run test2 --image=busybox --namespace=admission
```

## 生成证书
```shell
openssl genrsa -out validating-webhook.admission.svc.key 2048
openssl req -new -key validating-webhook.admission.svc.key -out validating-webhook.admission.svc.csr
openssl x509 -req -days 365 -in validating-webhook.admission.svc.csr -signkey validating-webhook.admission.svc.key -out validating-webhook.admission.svc.crt
```
