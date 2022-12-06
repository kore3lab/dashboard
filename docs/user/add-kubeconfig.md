# Add a KUBECONFIG

## Prerequisites

* `kubectl` and runnable

## Using a ServiceAccount
> if `client-certificate` is not available

* Create a serviceaccount and secret
```
$ kubectl create serviceaccount token-sa -n kube-system

$ kubectl apply -f - <<EOF
apiVersion: v1
kind: Secret
metadata:
  name: token-sa-secret
  namespace: kube-system
  annotations:
    kubernetes.io/service-account.name: token-sa
type: kubernetes.io/service-account-token
EOF
```

* Grant `cluster-admin` permission to the serviceaccount

```
$ kubectl create clusterrolebinding token-binding --clusterrole=cluster-admin --serviceaccount=kube-system:token-sa
```

* Create and verify the KUBECONFIG file  (ca.crt, token)
```
$ echo -e "apiVersion: v1
kind: Config
clusters:
- cluster:
    certificate-authority-data: $(kubectl get secret token-sa-secret -n kube-system -o jsonpath='{.data.ca\.crt}')
    server: $(kubectl config view -o jsonpath="{.clusters[?(@.name==\"$(kubectl config view -o jsonpath="{.contexts[?(@.name==\"$(kubectl config current-context)\")].context.cluster}")\")].cluster.server}")
  name: token-cluster
contexts:
- context:
    cluster: token-cluster
    user: token-user
  name: token
current-context: token
users:
- name: token-user
  user:
    token: $(kubectl get secret token-sa-secret -n kube-system -o jsonpath='{.data.token}' | base64 --decode)
" > kubeconfig-token.yaml

$ kubectl get nodes --kubeconfig="$(pwd)/kubeconfig-token.yaml"
```

## Using ServiceAccount on GKE Autopilot

* `gcloud container` 명령으로 kubeconfig 정보 구성

```
$ export KUBECONFIG="kubeconfig-gke-autopilot.yaml"
$ gcloud container clusters get-credentials autopilot-cluster-1 --region asia-northeast3 --project kore3-etri-cloudbarista
$ kubectl get nodes
```

* 아래와 같은 에러 발생 시  gcloud config 설정 수정

```
ERROR: (gcloud.container.clusters.get-credentials) get-credentials requires edit permission on ....
```

```
$ gcloud config set container/use_client_certificate False
```

* serviceaccount, secret 생성하고 `cluster-admin` ClusterRole 권한 부여

```
$ kubectl create serviceaccount admin-sa -n kube-public

$ kubectl apply -f - <<EOF
apiVersion: v1
kind: Secret
metadata:
  name: admin-sa-secret
  namespace: kube-public
  annotations:
    kubernetes.io/service-account.name: admin-sa
type: kubernetes.io/service-account-token
EOF

$ kubectl create clusterrolebinding admin-token-binding --clusterrole=cluster-admin --serviceaccount=kube-public:admin-sa
```

* kubeconfig 파일 생성 및 검증

```
$ echo -e "apiVersion: v1
kind: Config
clusters:
- cluster:
    certificate-authority-data: $(kubectl get secret admin-sa-secret -n kube-public -o jsonpath='{.data.ca\.crt}')
    server: $(kubectl config view -o jsonpath="{.clusters[?(@.name==\"$(kubectl config view -o jsonpath="{.contexts[?(@.name==\"$(kubectl config current-context)\")].context.cluster}")\")].cluster.server}")
  name: token-cluster
contexts:
- context:
  name: token
  context:
    cluster: token-cluster
    user: token-user
current-context: token
users:
- name: token-user
  user:
    token: $(kubectl get secret admin-sa-secret -n kube-public -o jsonpath='{.data.token}' | base64 --decode)
" > kubeconfig-gke-autopilot-token.yaml


$ kubectl get sa admin-sa -n kube-public --kubeconfig="$(pwd)/kubeconfig-gke-autopilot-token.yaml"
$ kubectl get secret admin-sa-secret  -n kube-public --kubeconfig="$(pwd)/kubeconfig-gke-autopilot-token.yaml"
$ kubectl get nodes --kubeconfig="$(pwd)/kubeconfig-gke-autopilot-token.yaml"
```