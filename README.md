# Kore3 ServiceMesh Platform
[![GitHub release](https://img.shields.io/github/v/release/aconsoftlab/servicemesh.svg)](https://github.com/acornsoftlab/kore3/releases/latest)

## Getting started
> 개발환경 구성 방법을 설명

### Preparation

아래 소프트웨어가 설치 및 $PATH 변수에 추가 필요

* Curl 7+
* Git 2.13.2+
* Docker 1.13.1+ ([installation manual](https://docs.docker.com/engine/installation/linux/docker-ce/ubuntu/))
* Golang 1.13.9+ ([installation manual](https://golang.org/dl/))
    * Dashboard uses `go mod` for go dependency management, so enable it with running `export GO111MODULE=on`.
* Node.js 12+ and npm 6+ ([installation with nvm](https://github.com/creationix/nvm#usage))
* Gulp.js 4+ ([installation manual](https://github.com/gulpjs/gulp/blob/master/docs/getting-started/1-quick-start.md))


### Clone

```
$ git clone https://github.com/acornsoftlab/kore3.git
$ cd kore3
```

### Install 

* npm denpendy
  * 만일 루트 권한으로 실행한 다면 `--unsafe-perm flag` 옵션지정
  * kubernetes dashboard 도 npm install 수행 

```
$ npm ci
$ cd dashboard
$ npm ci
```


* `subtree` 구성된 kubernetes/dashboard 에서 필요한 파일만 선별하고 트리를 업데이트 한다.

```
# sparse checkout 옵션 지정
$ git config core.sparsecheckout true

# 대상 파일 지정
$ vi .git/info/sparse-checkout

/*
!/dashboard
/dashboard/aio/gulp/*.*
/dashboard/dist
/dashboard/src/app/backend
/dashboard/.babelrc
/dashboard/.gitignore
/dashboard/package*.*
/dashboard/go.*
/dashboard/gulpfile.babel.js


# 트리 업데이트
$ git read-tree HEAD -m -u
```

### Run

```
$ export KUBECONFIG=~/.kube/config  # kubeconfig 파일 지정

$ npm run dev:graph            # 그래프 개발 (backend, graph)
$ npm run dev:ui               # UI 개발 (backend, frontend)
$ npm run start                # 개발 실행 (backend, dashboard, graph)
$ npm run start:backend        # backend 개발 실행
$ npm run start:dashboard      # kubernetes dashboard backend 개발 실행
$ npm run start:graph          # graph 개발 실행
```

* Dashboard 개발 시  `dashboard-metrics-scraper` 실행 필요

```
$ kubectl port-forward  svc/dashboard-metrics-scraper -n kubernetes-dashboard 8000
```

### Validation

```
# frontend
$ curl http://localhost:3000/

# backend
$ curl http://localhost:3001/healthy


# graph
$ curl http://localhost:3002/graph/mesh.html
$ curl http://localhost:3002/graph/rbac.html
$ curl http://localhost:3002/graph/topology.html

# dashboard
$ curl http://localhost:9090/api/v1/pod
$ curl http://localhost:9090/api/v1/pod?context=play.getapps.run
$ curl http://localhost:9090/api/v1/pod?context=apps-05
```

### Ports
* 3000 : front-end
* 3001 : backend (restful-api)
* 3002 : graph 개발
* 9090 : kubernetes dashboard backend

## Developments
> 개발 문서


## Architecture
> 아키텍쳐 및 구성

### Frontend (TODO)

* Frameworks : [nuxtJS](https://ko.nuxtjs.org/guide/plugins/)
* Template & Markup
  * [AdminLTE](https://adminlte.io/)
  * [bootstrap-vue v2.0.0](https://bootstrap-vue.org/docs/components/dropdown)
  * [Bootstrap v4.3.1](https://getbootstrap.com/)

```
$ export BACKEND_PORT="3001"
$ export BDASHBOARD_PORT="9090"
$ export KIALI_ROOT_URL="http://localhost:20001"
$ npm run start:frontend
```

#### 환경변수 (env)

|변수명     |설명             |기본값                 |
|---          |---              |---                    |
|BACKEND_URL  |backend Root URL |http://localhost:3001  |


### Backend

* backend restful api 
* language :  go-lang 1.15
* web frameworks : gin
* client-go 주요 참고 소스 
  * https://github.com/kubernetes/api/blob/master/core/v1/types.go
  * https://github.com/kubernetes/apimachinery/blob/master/pkg/apis/meta/v1/meta.go
  * 리스트 조회: https://github.com/kubernetes/client-go/blob/master/listers/core/v1 

#### 환경변수 (env)

|변수명     |설명                 |기본값               |
|---        |---                  |---                  |
|KUBECONFIG |kubeconfig 파일 위치 |${HOME}/.kube/config |

#### 제공 API

|apis                             |이름                             |비고  |예                                                  |   |
|---                              |---                              |---   |---                                                 |---|
|/apis/contexts                   |k8s cluster context 리스트 조회  |      |http://localhost:3001/apis/context                  |   |
|/apis/cluster/:cluster/topology  |토플로지 그래프 조회             |      |http://localhost:3001/apis/cluster/apps-05/topology |   |


### Dashboard

* Kubernetes dashboard 프로젝트(https://github.com/kubernetes/dashboard) 활용
* https://github.com/kubernetes/dashboard repository 를  `subtree` 로 구성 

* `subtree` 구성 방법

```
$ mkdir aconsoftlab.kore3
$ cd aconsoftlab.kore3
$ echo "# kore3" >> README.md
$ git init
$ git add README.md
$ git commit -m "initialize"
$ git remote add origin git@github.com:acornsoftlab/kore3.git
$ git subtree add --prefix=dashboard https://github.com/kubernetes/dashboard.git master
```


## Containerization

* Clean-up docker cache

```
$ docker system prune
```

* Definition 

```
$ BACKEND="kore3.backend"
$ DASHBOARD="kore3.dashboard"
$ FRONTEND="kore3.frontend"
```



* Build
```
$ docker build --tag acornsoftlab/${BACKEND}:latest ./src/app/backend   # backend
$ docker build --tag acornsoftlab/${DASHBOARD}:latest ./dashboard       # kubernetes-dashboard

# frontend
$ npm run build:graph                                     # graph 변경사항 반영할 경우 실행
$ docker build --tag acornsoftlab/${FRONTEND}:latest .

$ docker images | grep acornsoftlab
```

* Run on docker
```
$ docker run --rm -d -p 3001:3001 -v ${HOME}/.kube/config:/app/.kube/config --name ${BACKEND} acornsoftlab/${BACKEND}:latest
$ docker run --rm -d -p 9090:9090 -v ${HOME}/.kube/config:/app/.kube/config --name ${DASHBOARD} acornsoftlab/${DASHBOARD}:latest
$ docker run --rm -d -p 3000:3000 --name ${FRONTEND} acornsoftlab/${FRONTEND}:latest

$ docker ps
```

* Verify
```
# backend
$ curl http://localhost:3001/healthy
"healthy"

# dashboard
$ curl http://localhost:9090/api/v1/pod

# frontend 
# oepn in browser : http://localhost:3000/
```


* Debug
```
$ docker logs ${BACKEND}
$ docker logs ${DASHBOARD}
$ docker logs ${FRONTEND}

$ docker exec -it  ${FRONTEND} sh
```

* Push
```
# login
$ docker login -u acornsoftlab -p <password>

$ docker push acornsoftlab/${BACKEND}:latest
$ docker push acornsoftlab/${DASHBOARD}:latest
$ docker push acornsoftlab/${FRONTEND}:latest

$ TAG="v0.1.0"
$ docker tag acornsoftlab/${BACKEND}:latest acornsoftlab/${BACKEND}:${TAG}
$ docker tag acornsoftlab/${DASHBOARD}:latest acornsoftlab/${DASHBOARD}:${TAG}
$ docker tag acornsoftlab/${FRONTEND}:latest acornsoftlab/${FRONTEND}:${TAG}

$ docker push acornsoftlab/${BACKEND}:${TAG}
$ docker push acornsoftlab/${DASHBOARD}:${TAG}
$ docker push acornsoftlab/${FRONTEND}:${TAG}

$ docker images | grep acornsoftlab
```

## 참고

[nuxtjs](https://ko.nuxtjs.org/)
[nuxtjs github](https://github.com/nuxt/nuxt.js/)
[패스트캠퍼스 Vue.js 수업 자료](https://joshua1988.github.io/vue-camp/textbook.html)