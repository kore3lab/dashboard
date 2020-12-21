# Acornsoft Kubernetes Dashboard

## Preparation

아래 소프트웨어가 설치 및 $PATH 변수에 추가 필요

* Curl 7+
* Git 2.13.2+
* Docker 1.13.1+ ([installation manual](https://docs.docker.com/engine/installation/linux/docker-ce/ubuntu/))
* Golang 1.13.9+ ([installation manual](https://golang.org/dl/))
    * Dashboard uses `go mod` for go dependency management, so enable it with running `export GO111MODULE=on`.
* Node.js 12+ and npm 6+ ([installation with nvm](https://github.com/creationix/nvm#usage))
* Gulp.js 4+ ([installation manual](https://github.com/gulpjs/gulp/blob/master/docs/getting-started/1-quick-start.md))

## Getting started

* clone

```
$ git clone https://github.com/acornsoftlab/dashboard.git
$ cd dashboard
```

* `subtree` 구성된 kubernetes/dashboard 에서 사용하지 않는 파일을 제외하도록 지정하고 소스 트리를 업데이트

```
# sparse checkout 옵션 지정
$ git config core.sparsecheckout true

# 대상 파일 지정
$ cat <<EOF> .git/info/sparse-checkout
/*
!/dashboard
/dashboard/Dockerfile
/dashboard/.npmrc
/dashboard/.gitignore
/dashboard/aio/gulp/*.*
/dashboard/i18n
/dashboard/src/app/backend
/dashboard/.babelrc
/dashboard/.gitignore
/dashboard/package*.*
/dashboard/go.*
/dashboard/gulpfile.babel.js
EOF

# 트리 업데이트
$ git read-tree HEAD -m -u
```


* npm install

```
# frontend
$ npm i

# dashboard
$ cd dashboard
$ npm i

# graph
$ cd src/app/graph
$ npm i
```

* run

```
$ npm run start
```


* validation

```
# frontend
$ curl http://localhost:3000/

# backend
$ curl http://localhost:3001/healthy

# dashboard
$ curl http://localhost:9090/api/v1/pod

# graph
http://localhost:3002/topology.html
```

* Dashboard 개발 시  `dashboard-metrics-scraper` 실행 필요

```
$ kubectl port-forward  svc/dashboard-metrics-scraper -n kubernetes-dashboard 8000
```

## NPM 

* 개발
```
$ npm run start               # 실행 (backend, dashboard, frontend)
$ npm run start:frontend      # frontend 실행
$ npm run start:backend       # backend 실행
$ npm run start:dashboard     # kubernetes dashboard backend 실행
$ npm run start:graph         # graph 실행
$ npm run start:graph:backend # 그래프 개발 실행 (backend + graph)
```

* 빌드
```
$ npm run build:frontend      # frontend 빌드 (using on docker build)
$ npm run build:graph         # 그래프 빌드 frontend 에 변경된 최신 그래프 적용시 사용
$ npm run run                 # frontend container 에서 nuxt 실행 (docker image entrypoint) 
```

* Containerization

```
# docker build
$ npm run docker:build:frontend
$ npm run docker:build:backend
$ npm run docker:build:dashboard

# docker push
$ npm run docker:push:frontend    
$ npm run docker:push:backend
$ npm run docker:push:dashboard

# docker build & push
$ npm run docker:build:push:frontend    
$ npm run docker:build:push:backend
$ npm run docker:build:push:dashboard

# all (frontend, backend, dashboard)
$ npm run docker:build        # build
$ npm run docker:build:push   # build & push
```


### Using ports
* 3000 : front-end
* 3001 : backend (restful-api)
* 3002 : graph 개발
* 9090 : kubernetes dashboard backend


## Deployment


### Deploy on Docker

```
$ docker run --rm -d -p 3001:3001 -v ${HOME}/.kube/config:/app/.kube/config --name backend acornsoftlab/acornsoft-dashboard.backend:v0.1.1
$ docker run --rm -d -p 9090:9090 -v ${HOME}/.kube/config:/app/.kube/config --name dashboard acornsoftlab/acornsoft-dashboard.dashboard:v0.1.1
$ docker run --rm -d -p 3000:3000 -e BACKEND_PORT="3001" -e DASHBOARD_PORT="9090" --name frontend acornsoftlab/acornsoft-dashboard.frontend:v0.1.1
$ docker ps
```

### Deploy on Kubernetes

[Install on Kubernetes](./scripts/install/README.md)



## Front-End
> Web UI

* Frameworks : [nuxtJS](https://ko.nuxtjs.org/guide/plugins/)
* Template & Markup
  * [AdminLTE](https://adminlte.io/)
  * [bootstrap-vue v2.0.0](https://bootstrap-vue.org/docs/components/dropdown)
  * [Bootstrap v4.3.1](https://getbootstrap.com/)


### Run

```
$ npm run start:frontend
```

* 환경변수 (env)

|변수명           |설명                             |기본값 |
|---              |---                              |---    |
|BACKEND_PORT     |backend 서비스 포트              |3001   |
|BDASHBOARD_PORT  |kubernetes-dashboard 서비스 포트 |9090   |
|KIALI_PORT       |kiali 서비스 포트                |20001  |

```
$ export BACKEND_PORT="3001"
$ export BDASHBOARD_PORT="9090"
$ export KIALI_PORT="20001"
$ npm run start:frontend
```

### 참조

[nuxtjs](https://ko.nuxtjs.org/)
[nuxtjs github](https://github.com/nuxt/nuxt.js/)
[패스트캠퍼스 Vue.js 수업 자료](https://joshua1988.github.io/vue-camp/textbook.html)

## Back-End
> Custom backend rest-api

* backend restful api 
* language :  go-lang 1.15
* web frameworks : gin
* client-go 주요 참고 소스 
  * https://github.com/kubernetes/api/blob/master/core/v1/types.go
  * https://github.com/kubernetes/apimachinery/blob/master/pkg/apis/meta/v1/meta.go
  * 리스트 조회: https://github.com/kubernetes/client-go/blob/master/listers/core/v1 

### Run

```
$ npm run start:backend
```

* 환경변수 (env)

|변수명     |설명                 |기본값               |
|---        |---                  |---                  |
|KUBECONFIG |kubeconfig 파일 위치 |${HOME}/.kube/config |


### 제공 API

|apis                               |이름                               |비고  |예                                                  |   |
|---                                |---                                |---   |---                                                 |---|
|/api/contexts                      |k8s cluster context 리스트 조회    |      |http://localhost:3001/apis/context                  |   |
|/api/cluster/:cluster/topology     |토플로지 그래프 조회               |      |http://localhost:3001/apis/cluster/apps-05/topology |   |
|/api/_raw/:context:/:apiGroups/*   |해당 context 의 kubernetes Raw API |      |아래 참조                                           |   |


### Raw API
> [Kubernetes API Concepts](https://kubernetes.io/docs/reference/using-api/api-concepts/)

* 패턴
  * `:context` : kubeconfig Context 이름
  * `:apiGroups` : api Groups `kubectl api-resources -o wide` 으로 조회 가능

```
/api/_raw/:context:/:apiGroups/*
```

* apiGroups
```
$ kubectl api-resources -o wide

# CRD 경우
$ kubectl get crd
$ kubectl get crd virtualservices.networking.istio.io -o jsonpath="{.spec.group}"
```


* Raw API 호출 예제

```
# Root URL 정의
$ export RootUrl="http://localhost:3001/api/_raw/$(kubectl config current-context)"


# core API (Core 그룹은 공백값)
$ curl ${RootUrl}//v1/pods
$ curl ${RootUrl}//v1/namespaces 
$ curl ${RootUrl}//v1/nodes
$ curl ${RootUrl}/apps/v1/deployments 

# core API (in namespace)
$ curl ${RootUrl}//v1/namespaces/kube-system/pods
$ curl ${RootUrl}/apps/v1/namespaces/kube-system/deployments

# CRD 
$ curl ${RootUrl}/networking.istio.io/v1alpha3/virtualservices
$ curl ${RootUrl}/networking.istio.io/v1alpha3/namespaces/bookinfo/virtualservices
```

### Dashboard
> Kubernetes dashbaord

* Kubernetes dashboard 프로젝트(https://github.com/kubernetes/dashboard) 활용
* https://github.com/kubernetes/dashboard repository 를  `subtree` 로 구성 

### `subtree` 구성 방법

```
$ mkdir aconsoftlab.dashboard
$ cd aconsoftlab.dashboard
$ echo "# Dashboard" >> README.md
$ git init
$ git add README.md
$ git commit -m "initialize"
$ git remote add origin git@github.com:acornsoftlab/dashboard.git
$ git subtree add --prefix=dashboard https://github.com/kubernetes/dashboard.git master
```

## Graph

### Build

* frontend와 연동되어 build하면 frontend 에 복사된다.

```
$ npm run build:graph
```

```
http://localhost:3002/topology.html   # topology graph
http://localhost:3002/mesh.html       # mesh graph (deprecated)
http://localhost:3002/rbac.html       # rbac graph
```
