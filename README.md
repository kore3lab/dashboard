# Kore3 ServiceMesh Platform

## Getting start
개발환경 구성 방법을 설명

### Preparation

아래 소프트웨어가 설치 및 $PATH 변수에 추가 필요

* Curl 7+
* Git 2.13.2+
* Docker 1.13.1+ ([installation manual](https://docs.docker.com/engine/installation/linux/docker-ce/ubuntu/))
* Golang 1.13.9+ ([installation manual](https://golang.org/dl/))
    * Dashboard uses `go mod` for go dependency management, so enable it with running `export GO111MODULE=on`.
* Node.js 12+ and npm 6+ ([installation with nvm](https://github.com/creationix/nvm#usage))
* Gulp.js 4+ ([installation manual](https://github.com/gulpjs/gulp/blob/master/docs/getting-started/1-quick-start.md))


### 프로젝트 환경 설정

* 프로젝트 clone

```
git clone https://github.com/acornsoftlab/servicemesh.git
cd servicemesh
```

* npm denpendy 설치

```
npm ci
```

* 만일 루트 권한으로 실행한 다면 `--unsafe-perm flag` 지정

```
npm ci --unsafe-perm
```

* kubernetes 필요파일 선별 

```
git config core.sparsecheckout true
echo '/*' >> .git/info/sparse-checkout
echo '!/dashboard' >> .git/info/sparse-checkout
echo '/dashboard/aio/gulp/*.*' >> .git/info/sparse-checkout
echo '/dashboard/dist' >> .git/info/sparse-checkout
echo '/dashboard/src/app/backend' >> .git/info/sparse-checkout
echo '/dashboard/.babelrc' >> .git/info/sparse-checkout
echo '/dashboard/.gitignore' >> .git/info/sparse-checkout
echo '/dashboard/package*.*' >> .git/info/sparse-checkout
echo '/dashboard/go.*' >> .git/info/sparse-checkout
echo '/dashboard/gulpfile.babel.js' >> .git/info/sparse-checkout
```

* tree 업데이트

```
git read-tree HEAD -m -u
```

### Kubernetes dashboard 프로젝트 활용

* Kubernetes dashboard 프로젝트(https://github.com/kubernetes/dashboard) 를 백엔드 활용하기 위해 `subtree` 로 구성 

* 구성과정

```
mkdir aconsoftlab.servicemesh
cd aconsoftlab.servicemesh
echo "# servicemesh" >> README.md
git init
git add README.md
git commit -m "initialize"
git remote add origin git@github.com:acornsoftlab/servicemesh.git
git subtree add --prefix=dashboard https://github.com/kubernetes/dashboard.git master
```

* 빌드

```
npm run build:backend
```

* 실행

```
dashboard/dist/amd64/dashboard --kubeconfig ~/.kube/config
```

```
export KUBECONFIG=~/.kube/config
npm run start:backend
```

* 실행 확인

```
curl http://localhost:9090/api/v1/pod
curl http://localhost:9090/api/v1/pod?context=play.getapps.run
curl http://localhost:9090/api/v1/pod?context=apps-05
```
