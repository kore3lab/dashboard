# Getting started

## Preparation

아래 소프트웨어 설치 및 $PATH 변수에 추가 필요

* Curl 7+
* Git 2.13.2+
* Docker 1.13.1+ ([installation manual](https://docs.docker.com/engine/installation/linux/docker-ce/ubuntu/))
* Golang 1.13.9+ ([installation manual](https://golang.org/dl/))
* Node.js 12+ and npm 6+ ([installation with nvm](https://github.com/creationix/nvm#usage))
* Gulp.js 4+ ([installation manual](https://github.com/gulpjs/gulp/blob/master/docs/getting-started/1-quick-start.md))

## Configuration of development environment

* clone

```
$ git clone https://github.com/kore3lab/dashboard.git
$ cd dashboard
```

* Installation dependencies (frontend,graph)

```
# graph
$ cd ../src/app/graph
$ npm i

$ cd ../../..
$ npm i
```

* Run

```
$ npm run start
```

* verify the running

```
# frontend
$ curl http://localhost:3000/

# backend
$ curl http://localhost:3001/healthy

# metrics scraper
$ curl http://localhost:8000/api/v1

# terminal
$ curl http://localhost:3003/healthy
```

### Using ports
* 3000 : front-end
* 3001 : backend (restful-api)
* 3002 : graph 개발
* 3003 : terminal
* 8000 : metrics scraper


## Build & Run


* Develop
```
$ npm run start                 # 실행 (backend, frontend, metrics-scraper)
$ npm run start:frontend        # frontend 실행
$ npm run start:backend         # backend 실행
$ npm run start:metrics-scraper # metrics scraper 실행
$ npm run start:terminal        # container 에서 terminal 실행 
$ npm run start:graph           # graph 실행
$ npm run start:graph:backend   # 그래프 개발 실행 (backend + graph)
```

* Build
```
$ npm run build:frontend      # frontend 빌드 (using on docker build)
$ npm run build:graph         # 그래프 빌드 frontend 에 변경된 최신 그래프 적용시 사용
$ npm run run                 # frontend container 에서 nuxt 실행 (docker image entrypoint) 
$ npm run build:terminal      # terminal 빌드 (using on docker build)
```

* Containerization
  * default tag : "latest" (.npmrc 파일 참조)
  * `--kore-board:docker_image_tag=<value>` 옵션으로 latest 대신 tag 지정 가능

```
# docker build
$ npm run docker:build:frontend
$ npm run docker:build:backend
$ npm run docker:build:metrics-scraper
$ npm run docker:build:terminal

# docker push
$ npm run docker:push:frontend    
$ npm run docker:push:backend
$ npm run docker:push:metrics-scraper
$ npm run docker:push:terminal

# docker build & push
$ npm run docker:build:push:frontend    
$ npm run docker:build:push:backend
$ npm run docker:build:push:metrics-scraper
$ npm run docker:build:push:terminal

# all (frontend, backend)
$ npm run docker:build        # build
$ npm run docker:build:push   # build & push

# tag 를 옵션으로 지정하는 예
$ npm run docker:build:backend --kore-board:docker_image_tag=v0.2.0
```


## Debugging for "in-cluster" environment 

```
# 이전 token & ca 파일 삭제
$ sudo rm /var/run/secrets/kubernetes.io/serviceaccount/token /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
$ sudo mkdir -p /var/run/secrets/kubernetes.io/serviceaccount/

# namespace & serviceaccount 이름 지정
$ export NAMESPACE="kore"
$ export SERVICE_ACCOUNT="kore-board"

# token & ca 파일 생성
$ echo -n "$(kubectl get secret -n ${NAMESPACE} $(kubectl get sa ${SERVICE_ACCOUNT} -n ${NAMESPACE} -o jsonpath={.secrets..name} | cut -f1 -d ' ') -o jsonpath='{$.data.token}' | base64 --decode)" | sudo tee /var/run/secrets/kubernetes.io/serviceaccount/token
$ echo -n "$(kubectl config view --raw=true -o jsonpath='{.clusters[0].cluster.certificate-authority-data}' | base64 --decode)" | sudo tee /var/run/secrets/kubernetes.io/serviceaccount/ca.crt

# set api-server endpoint
$ export KUBERNETES_SERVICE_HOST=$(kubectl config view --raw=true -o jsonpath='{.clusters[0].cluster.server}' |  awk -F/ '{print $3}' |  awk -F: '{print $1}')
$ export KUBERNETES_SERVICE_PORT=$(kubectl config view --raw=true -o jsonpath='{.clusters[0].cluster.server}' |  awk -F/ '{print $3}' |  awk -F: '{print $2}')

# create a empty kubeconfig file
$ rm "$(pwd)/.tmp/config"
$ export KUBECONFIG="$(pwd)/.tmp/config"
```

