# Acornsoft Kubernetes Dashboard

## Preparation

아래 소프트웨어 설치 및 $PATH 변수에 추가 필요

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

* `subtree`로 구성된 소스 중 사용되지 않는 파일을 제외하도록 지정하고 소스 트리를 업데이트합니다.

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
!/src/app/metrics-scraper
/src/app/metrics-scraper/hack/build.sh
/src/app/metrics-scraper/pkg
/src/app/metrics-scraper/vendor
/src/app/metrics-scraper/Dockerfile
/src/app/metrics-scraper/go.*
/src/app/metrics-scraper/*.go
!/src/app/metrics-scraper/**/*_test.go
EOF

# 트리 업데이트
$ git read-tree HEAD -m -u
```


* Install dependencies (frontend, kubernetes-dashboard, graph)

```
# frontend
$ npm i

# dashboard
$ cd dashboard
$ npm i

# graph
$ cd ../src/app/graph
$ npm i
```

* Run

```
$ cd ../../..
$ npm run start
```

* Verified running

```
# frontend
$ curl http://localhost:3000/

# backend
$ curl http://localhost:3001/healthy

# dashboard
$ curl http://localhost:9090/api/v1/pod

# metrics scraper
$ curl http://localhost:8000/api/v1

# graph
http://localhost:3002/topology.html
```

## NPM 

* Develop
```
$ npm run start                 # 실행 (backend, dashboard, frontend, metrics-scraper)
$ npm run start:frontend        # frontend 실행
$ npm run start:backend         # backend 실행
$ npm run start:dashboard       # kubernetes dashboard backend 실행
$ npm run start:metrics-scraper # metrics scraper 실행
$ npm run start:graph           # graph 실행
$ npm run start:graph:backend   # 그래프 개발 실행 (backend + graph)
```

* Build
```
$ npm run build:frontend      # frontend 빌드 (using on docker build)
$ npm run build:graph         # 그래프 빌드 frontend 에 변경된 최신 그래프 적용시 사용
$ npm run run                 # frontend container 에서 nuxt 실행 (docker image entrypoint) 
```

* Containerization
  * default tag : "latest" (.npmrc 파일 참조)
  * `--acornsoft-dashboard:docker_image_tag=<value>` 옵션으로 latest 대신 tag 지정 가능

```
# docker build
$ npm run docker:build:frontend
$ npm run docker:build:backend
$ npm run docker:build:dashboard
$ npm run docker:build:metrics-scraper

# docker push
$ npm run docker:push:frontend    
$ npm run docker:push:backend
$ npm run docker:push:dashboard
$ npm run docker:push:metrics-scraper

# docker build & push
$ npm run docker:build:push:frontend    
$ npm run docker:build:push:backend
$ npm run docker:build:push:dashboard
$ npm run docker:build:push:metrics-scraper

# all (frontend, backend, dashboard)
$ npm run docker:build        # build
$ npm run docker:build:push   # build & push

# tag 를 옵션으로 지정하는 예
$ npm run docker:build:backend --acornsoft-dashboard:docker_image_tag=v0.2.0
```

### Using ports
* 3000 : front-end
* 3001 : backend (restful-api)
* 3002 : graph 개발
* 9090 : kubernetes dashboard backend
* 8000 : metrics scraper


## Deployment


### Deploy on Docker

```
$ docker run --rm -d\
    --name metrics-scraper -p 8000:8000\
    -v ${HOME}/.kube/config:/app/.kube/config\
    ghcr.io/acornsoftlab/acornsoft-dashboard.metrics-scraper:v0.1.1\
    --kubeconfig=/app/.kube/config --db-file=metrics.db

$ docker run --rm -d\
    --name backend -p 3001:3001\
    -v ${HOME}/.kube/config:/app/.kube/config\
    ghcr.io/acornsoftlab/acornsoft-dashboard.backend:v0.1.1\
    --kubeconfig=/app/.kube/config

$ docker run --rm -d\
    --name dashboard -p 9090:9090\
    -v ${HOME}/.kube/config:/app/.kube/config\
    --link metrics-scraper:metrics-scraper\
    ghcr.io/acornsoftlab/acornsoft-dashboard.dashboard:v0.1.1\
    --kubeconfig=/app/.kube/config\
    --sidecar-host=http://metrics-scraper:8000

$ docker run --rm -d\
    --name frontend -p 3000:3000\
    -e BACKEND_PORT="3001"\
    -e DASHBOARD_PORT="9090"\
    ghcr.io/acornsoftlab/acornsoft-dashboard.frontend:v0.1.1

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
> Backend rest-api

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

* Argument

|변수명     |설명                 |기본값               |
|---        |---                  |---                  |
|kubeconfig |kubeconfig 파일 위치 |                     |


* 환경변수 (env)

|변수명     |설명                 |기본값               |
|---        |---                  |---                  |
|KUBECONFIG |kubeconfig 파일 위치 |${HOME}/.kube/config |



### API

[Acornsoft Dashbard Backend](https://github.com/acornsoftlab/dashboard/blob/master/src/app/backend/README.md) 참조 


## Dashboard
> Wrapping Kubernetes-dashbaord

* Kubernetes dashboard 프로젝트(https://github.com/kubernetes/dashboard) 활용
* https://github.com/kubernetes/dashboard repository 를  `subtree` 로 구성 


### `subtree` 구성 방법

```
$ git subtree add --squash --prefix=dashboard https://github.com/kubernetes/dashboard.git master
```

### Run

```
$ npm run start:dashboard
```

* Argument

|변수명               |설명                                    |기본값                |
|---                  |---                                     |---                   |
|--kubeconfig         |kubeconfig 파일 위치                    |                      |
|--sidecar-host       |metrics 정보(metrics-scraper) 요청 URL  |http://localhost:8000 |


* 환경변수 (env)

|변수명     |설명                                     |기본값                                 |
|---          |---                                    |---                                    |
|KUBECONFIG   |kubeconfig 파일 위치                   |${HOME}/.kube/config                   |


## Metrics-Scraper
> Wrapping Kubernetes-sig dashbaord-metrics-scraper

* Kubernetes dashboard-metrics-scraper(https://github.com/kubernetes-sigs/dashboard-metrics-scraper) 활용
* https://github.com/kubernetes-sigs/dashboard-metrics-scraper repository 를  `subtree` 로 구성 


### `subtree` 구성 방법

```
$ git subtree add --squash --prefix=src/app/metrics-scraper https://github.com/kubernetes-sigs/dashboard-metrics-scraper.git master
```

### Run

```
$ npm run start:metrics-scraper
```

* Argument

|변수명               |설명                       |기본값           |
|---                  |---                        |---              |
|--kubeconfig         |kubeconfig 파일 위치       |                 |
|--db-file            |sqllite database file path |/tmp/metrics.db  |
|--metric-resolution  |metrics 수집 주기          |1m0s             |
|--metric-duration    |metrics 적산값 유지 기간   |15m0s            |
|--log-level          |로그 레벨                  |                 |
|--namespace          |                           |                 |

* 환경변수 (env)

|변수명     |설명                 |기본값               |
|---        |---                  |---                  |
|KUBECONFIG |kubeconfig 파일 위치 |${HOME}/.kube/config |


### API

|URL Pattern                                                                 |Method |설명                               |
|---                                                                         |---    |---                                |
|/api/v1/clusters/:cluster/nodes/:node/metrics/:metrics                      |GET    |클러스터 Node metrics 조회         |
|/api/v1/nodes/:node/metrics/:metrics                                        |GET    |default 클러스터 노드 metrics 조회 |
|/api/v1/clusters/:cluster/namespaces/:namespaces/pods/:pod/metrics/:metrics |GET    |클러스터 Pod metrics 조회          |
|/api/v1/namespaces/:namespaces/pods/:pod/metrics/:metrics                   |GET    |default 클러스터 Pod metrics 조회  |

* 변수
  * `:cluster` : Kubeconfig context name
  * `:node` :  Node name
  * `:metrics` : `cpu` or `memory`
  * `:pod` : Pod name

* Examples

```
$ curl -X GET http://localhost:8000/api/v1/clusters/apps-05/nodes/apps-114/metrics/cpu
$ curl -X GET http://localhost:8000/api/v1/nodes/apps-114/metrics/cpu
$ curl -X GET http://localhost:8000/api/v1/clusters/apps-06/namespaces/default/pods/dnsutils-797cbd6f5f-8sq8t/metrics/memory
$ curl -X GET http://localhost:8000/api/v1/namespaces/default/pods/dnsutils-797cbd6f5f-8sq8t/metrics/memory
```


## Graph
> D3.js based javascript graph library

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

## Link
* https://github.com/acornsoftlab/dashboard
* https://github.com/kubernetes/client-go
* https://github.com/kubernetes/dashboard
* https://bootstrap-vue.org/docs/components
* https://github.com/gin-gonic/gin