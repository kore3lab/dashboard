# Install on Kubernetes

## Prerequites

* Create a namespace
```
$ kubectl create ns acornsoft-dashboard
```

* Careate a configmap `acornsoft-dashboard-kubeconfig`  (in-cluster mode 제외)

```
$ kubectl create configmap acornsoft-dashboard-kubeconfig --from-file=config=${HOME}/.kube/config -n acornsoft-dashboard
```

* If exists a configmap then Careate a configmap `acornsoft-dashboard-kubeconfig`  (in-cluster mode 제외)
```
$ kubectl create configmap acornsoft-dashboard-kubeconfig --from-file=config=${HOME}/.kube/config --dry-run -o yaml | kubectl apply  -n acornsoft-dashboard -f -
```

* metrics-server 가 설치되어 있지 않다면 metrics-server 설치 (args --kubelet-insecure-tls 추가)
```
$ kubectl get po  -n kube-system | grep metrics-server

$ kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
```


## Install using Yaml

* Install
```
$ kubectl apply -f kuberntes/recommended.yaml
```

* Clean-up
```
$ kubectl delete -f kuberntes/recommended.yaml
```


## Install using Helm Chart

* dry-run

```
$ kubectl create ns acornsoft-dashboard
$ helm install --dry-run --debug -n acornsoft-dashboard acornsoft-dashboard ./kuberntes/helm-chart/ \
  --set backend.service.type=NodePort \
  --set backend.service.nodePort=30081 \
  --set frontend.service.type=NodePort \
  --set frontend.service.nodePort=30080
```

* Install

```
$ helm install -n acornsoft-dashboard acornsoft-dashboard ./kuberntes/helm-chart/ \
  --set backend.service.type=NodePort \
  --set backend.service.nodePort=30081 \
  --set frontend.service.type=NodePort \
  --set frontend.service.nodePort=30080

$ helm list
```

* Clean-up
```
$ helm uninstall acornsoft-dashboard
```


## Install "in-cluster" mode
> Install for Single Cluster

```
# metrics-scraper

$ kubectl run metrics-scraper -n ${NAMESPACE}\
  --image=ghcr.io/acornsoftlab/acornsoft-dashboard.metrics-scraper:latest --port=8000\
  -- --db-file=metrics.db
$ kubectl expose pod metrics-scraper -n ${NAMESPACE} --port=8000 --name=metrics-scraper


# backend

$ kubectl run backend -n ${NAMESPACE}\
  --image=ghcr.io/acornsoftlab/acornsoft-dashboard.backend:latest --port=3001
$ kubectl expose pod backend -n ${NAMESPACE} --name=backend --type='NodePort' --port=3001


# rbac
$ kubectl create role acornsoft-dashboard -n ${NAMESPACE} --resource=* --verb=*
$ kubectl create rolebinding acornsoft-dashboard -n ${NAMESPACE} --role=acornsoft-dashboard --serviceaccount=${NAMESPACE}:default
$ kubectl create clusterrolebinding acornsoft-dashboard --clusterrole=cluster-admin --serviceaccount=${NAMESPACE}:default


# frontend

$ BACKEND_PORT="$(kubectl get svc/backend -n ${NAMESPACE} -o jsonpath="{.spec.ports[0].nodePort}")"

$ kubectl run frontend -n ${NAMESPACE}\
  --image=ghcr.io/acornsoftlab/acornsoft-dashboard.frontend:latest\
  --port=3000\
  --env="BACKEND_PORT=${BACKEND_PORT}"

$ kubectl expose pod frontend -n ${NAMESPACE} --name=frontend --type='NodePort' --port=3000

```

* Clean-up

```
$ kubectl delete -n ${NAMESPACE} pod/backend pod/frontend pod/metrics-scraper
$ kubectl delete -n ${NAMESPACE} service/backend service/frontend service/metrics-scraper
$ kubectl delete -n ${NAMESPACE} role/acornsoft-dashboard rolebinding/acornsoft-dashboard
$ kubectl delete clusterrolebinding/acornsoft-dashboard
$ kubectl delete ns ${NAMESPACE}
```


### Verify

* Open in your browser

```
$ echo "http://<end-point ip>:$(kubectl get svc/frontend -n ${NAMESPACE} -o jsonpath="{.spec.ports[0].nodePort}")"
```


## Developments

* Packaging

```
$ helm package helm-chart/    # packaging  (tgz file 생성)
$ helm repo index .           # index.yaml 파일 생성/수정
```
