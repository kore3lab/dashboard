# Installation Guide

## Prerequisites

### Metrics-Server

* Install metrics-server on target clusters (add startup option `--kubelet-insecure-tls` )
* https://github.com/kubernetes-sigs/metrics-server

```
# check for installation
$ kubectl get po  -n kube-system | grep metrics-server

# installation
$ kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
```

## Kubernetes

* first, create a kubeconfig configmap `kore-board-kubeconfig`

```
$ kubectl create ns kore
$ kubectl create configmap kore-board-kubeconfig --from-file=config=${HOME}/.kube/config -n kore
```

* if you modify it

```
$ kubectl create configmap kore-board-kubeconfig --from-file=config=${HOME}/.kube/config --dry-run -o yaml | kubectl apply  -n kore -f -
```

### Installation using Yaml

* Installation
```
$ kubectl apply -f ./scripts/install/kuberntes/recommended.yaml
```

* clean-up
```
$ kubectl delete -f ./scripts/install/kuberntes/recommended.yaml
```

### Installation using Helm Chart


* Installation

```
$ helm install -n kore kore-board ./scripts/install/kuberntes/helm-chart/ \
  --set backend.service.type=NodePort \
  --set backend.service.nodePort=30081 \
  --set frontend.service.type=NodePort \
  --set frontend.service.nodePort=30080

$ helm list
```

* clean-up
```
$ helm uninstall kore-board
```

## Docker

### Installation using `docker-compose`

* Installation
```
$ docker-compose -f ./scripts/install/docker-compose.yaml up -d
```

* clean-up
```
$ docker-compose -f ./scripts/install/docker-compose.yaml down
```

### Installation using `docker run`

* Installation

```
$ docker volume create data
$ docker volume create kubeconfig

$ docker run --rm -d --name metrics-scraper \
    -v "kubeconfig:/app/.kube"\
    -v "data:/app/data"\
    ghcr.io/acornsoftlab/kore-board.metrics-scraper:latest --kubeconfig=/app/.kube/config --db-file=/app/data/metrics.db

$ docker run --rm -d --name backend \
    -v "kubeconfig:/app/.kube" \
    --link metrics-scraper:metrics-scraper \
    --privileged \
    ghcr.io/acornsoftlab/kore-board.backend:latest --kubeconfig=/app/.kube/config --metrics-scraper-url=http://metrics-scraper:8000

$ docker run --rm -d --name frontend\
    -p 3000:80\
    --link backend:backend\
    ghcr.io/acornsoftlab/kore-board.frontend:latest
```

* clean-up
```
$ docker stop frontend backend metrics-scraper
$ docker volume rm data kubeconfig
```

