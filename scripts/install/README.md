# Installation

## Prerequites

### Metrics-Server

* 조회 대상 클러스터에 metrics-server 설치 (args --kubelet-insecure-tls 추가)

```
$ kubectl get po  -n kube-system | grep metrics-server

$ kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
```

## Docker

### Installation using `docker-compose`

* Installation
```
$ export KUBECONFIG="${HOME}/.kube/config"

$ docker-compose -f docker-compose/docker-compose.yaml up -d
```
dck
* clean-up
```
$ docker-compose -f docker-compose/docker-compose.yaml down
```

### Installation using `docker run`

* Installation

```
$ export KUBECONFIG="${HOME}/.kube/config"

$ docker run --rm -d --name metrics-scraper \
    -v "${KUBECONFIG}:/app/.kube/config"\
    -v "$(pwd):/app/data"\
    ghcr.io/acornsoftlab/kore-board.metrics-scraper:latest --kubeconfig=/app/.kube/config --db-file=/app/data/metrics.db

$ docker run --rm -d --name backend \
    -p 3001:3001\
    -v "${KUBECONFIG}:/app/.kube/config"\
    --link metrics-scraper:metrics-scraper\
    ghcr.io/acornsoftlab/kore-board.backend:latest --kubeconfig=/app/.kube/config --metrics-scraper-url=http://metrics-scraper:8000

$ docker run --rm -d --name frontend\
    -p 3000:80\
    -v "$(pwd):/tmp"\
    -v "$(pwd)/docker-compose/default.conf:/etc/nginx/conf.d/default.conf"\
    -v "$(pwd)/docker-compose/nginx.conf:/etc/nginx/nginx.conf"\
    --link backend:backend\
    ghcr.io/acornsoftlab/kore-board.frontend:latest
```

* clean-up
```
$ docker stop frontend backend metrics-scraper
```

## Kubernetes

### Prerequites  : create a kubeconfig configmap `kore-board-kubeconfig`

* Create
```
$ kubectl create ns kore
$ kubectl create configmap kore-board-kubeconfig --from-file=config=${HOME}/.kube/config -n kore
```

* Modify

```
$ kubectl create configmap kore-board-kubeconfig --from-file=config=${HOME}/.kube/config --dry-run -o yaml | kubectl apply  -n kore -f -
```

### Installation using Yaml

* Installation
```
$ kubectl apply -f kuberntes/recommended.yaml
```

* clean-up
```
$ kubectl delete -f kuberntes/recommended.yaml
```

### Installation using Helm Chart


* Installation

```
$ helm install -n kore kore-board ./kuberntes/helm-chart/ \
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

## Custom Token 

* backend `--token` 옵션 활용해 custom token 지정

* examples

```
# 사용자 token 값 생성 및 token 파일에 저장

$ TOKEN="$(pwd)/token"
$ rm "${TOKEN}"
$ echo -n "$(openssl rand -hex 32)" > "${TOKEN}"
$ cat "${TOKEN}" && echo ""

# backend 실행 시 volumn mount 하고 시작옵션 --token 에 해당 파일 지정 
$ docker run --rm -d --name backend \
    -v "${KUBECONFIG}:/app/.kube/config"\
    -v "${TOKEN}:/app/token"\
    --link metrics-scraper:metrics-scraper\
    ghcr.io/acornsoftlab/kore-board.backend:latest --kubeconfig=/app/.kube/config --token=/app/token --metrics-scraper-url=http://metrics-scraper:8000
```

## Developments

### Helm Chart

* output yaml

```
$ helm template --debug -n kore kore-board ./kuberntes/helm-chart/ \
  --set frontend.service.type=NodePort \
  --set frontend.service.nodePort=30080
```

# apply test

```
$ kubectl create ns kore
$ helm template --debug -n kore kore-board ./kuberntes/helm-chart/ \
  --set frontend.service.type=NodePort \
  --set frontend.service.nodePort=30080 \
  | kubectl apply -n kore -f -
```


* dry-run

```
$ helm install --dry-run --debug -n kore kore-board ./kuberntes/helm-chart/ \
  --set frontend.service.type=NodePort \
  --set frontend.service.nodePort=30080
```



* Packaging helm-chart

```
$ helm package helm-chart/    # packaging  (tgz file 생성)
$ helm repo index .           # index.yaml 파일 생성/수정
```
