# Back-End
> Backend rest-api

* backend restful api 
* language :  go-lang 1.15
* web frameworks : gin
* client-go 주요 참고 소스 
  * https://github.com/kubernetes/api/blob/master/core/v1/types.go
  * https://github.com/kubernetes/apimachinery/blob/master/pkg/apis/meta/v1/meta.go
  * 리스트 조회: https://github.com/kubernetes/client-go/blob/master/listers/core/v1 
* Web Terminal 주요 참고 소스
  * https://github.com/fanux/fist/tree/master/terminal
  * https://github.com/rancher: 소스 코드는 없는 듯. 
  * https://github.com/du2016/web-terminal-in-go
  * https://github.com/lf1029698952/kube-webshell 


## Run

```
$ npm run start:backend
```

* Arguments

|이름                   |기본값                 |설명                                                                                       |
|---                    |---                    |---                                                                                        |
|--kubeconfig           |                       |kubeconfig 파일 위치                                                                       |
|--log-level            |info                   |로그 레벨(panic,fatal,error,warning,info,debug,trace) https://github.com/sirupsen/logrus)  |
|--metrics-scraper-url  |http://localhost:8000  |metrics-scraper api url                                                                    |


* 환경변수 (env)

|이름       |기본값 |설명                 |
|---        |---    |---                  |
|KUBECONFIG |       |kubeconfig 파일 위치 |


* 제약사항
  * Web Terminal (context/pod/container 터미널접속) 기능은 linux OS의 unshare기능을 이용하기 때문에 backend가 container 환경일때 동작한다
    아래 링크를 참조하여 container deploy 후 사용가능
    [Deploy on Docker](#deploy-on-docker)


## API

|URL Pattern                        |Method |설명                                     |
|---                                |---    |---                                      |
|/api/clusters                      |GET    |k8s cluster context 리스트 조회          |
|/api/clusters/:cluster/topology    |GET    |토플로지 그래프 조회                     |
|/api/topology                      |GET    |토플로지 그래프 조회 (default cluster)   |
|/api/clusters/:cluster/dashboard   |GET    |Dashboard 데이터 조회                    |
|/api/dashboard                     |GET    |Dashboard 데이터 조회 (default cluster)  |

* Examples

```
$ curl -X GET http://localhost:3001/api/clusters
$ curl -X GET http://localhost:3001/api/topology
$ curl -X GET http://localhost:3001/api/clusters/apps-05/dashboard
```


## Kubernetes Raw API
> 멀티 클러스터를 지원하는 Kubernetes API Proxy API


### Kubernetes API
* [Kubernetes API Concepts](https://kubernetes.io/docs/reference/using-api/api-concepts/)
* [OepnAPI spec.](https://raw.githubusercontent.com/kubernetes/kubernetes/master/api/openapi-spec/swagger.json)
* Kubernetes API 에서 제공하는 resource 와 resource의 api-group은 `kubectl api-resources -o wide` 으로 resource와 apiGroup 조회 가능

```
$ kubectl api-resources -o wide

# CRD 경우
$ kubectl get crd
$ kubectl get crd virtualservices.networking.istio.io -o jsonpath="{.spec.group}"
```

### URL Pattern

* URL Patter은 다음과 같이 URL Prefix : `/raw` 와 kubernetes-api URL로 구성
* kubernetes-api 는 각 resource 의 `metadata.selfLink` 참조 가능

```
/raw/<Kubernetes-api URL>
```

* 아래 설명에서 사용하는 변수는 다음과 같습니다.
  * `:cluster` : Kubeconfig context name
  * `:version` :  Resource spec. version
  * `:resource` : Resource name
  * `:apiGroups` : Resource groups 


### Apply APIs
> Create a resource

|URL Pattern            |Method |설명                     |
|---                    |---    |---                      |
|/raw/clusters/:cluster |POST   |Applay                   |
|/raw                   |POST   |Apply (default cluster)  |

* Example

```
$ curl -X POST -H "Content-Type: application/json" http://localhost:3001/raw -d @- <<EOF
{
    "apiVersion": "v1",
    "kind": "Namespace",
    "metadata": {
        "name": "test-namespace"
    }
}
EOF
```

### Update APIs
> Update a resource

|URL Pattern            |Method |설명                     |비고             |
|---                    |---    |---                      |---              |
|/raw/clusters/:cluster |PATCH  |Update                   |                 |
|/raw                   |PATCH  |Update (default cluster) |                 |


* Reuqest Header `Content-Type` 으로 patch 방식 선택 (
  * patch 방식에 대한 이해 - [Update API Objects in Place Using kubectl patch](https://kubernetes.io/docs/tasks/manage-kubernetes-objects/update-api-object-kubectl-patch/)
  * [JSON Merge Patch (RFC 7396)](https://tools.ietf.org/html/rfc7386)  : `Content-Type : application/merge-patch+json`
  * [JSON Patch (RFC 6902)](https://tools.ietf.org/html/rfc6902) : `Content-Type : application/json-patch+json`
  * Strategic merge patch : `Content-Type : application/strategic-merge-patch+json`


* JSON merge patch example

```
$ curl -X PATCH -H "Content-Type: application/merge-patch+json" http://localhost:3001/raw/api/v1/namespaces/default/pods/busybox -d @- <<EOF
{
    "metadata": {
        "labels": {
            "app": "busybox-merge"
        }
    }
}
EOF

$ kubectl get po busybox -n default  -o jsonpath="{.metadata.labels}"
```

* JSON patch example

```
$ curl -X PATCH -H "Content-Type: application/json-patch+json" http://localhost:3001/raw/api/v1/namespaces/default/pods/busybox -d @- <<EOF
[
    {
        "op": "replace", 
        "path": "/metadata/labels/app", 
        "value":"busybox-json"
    }
]
EOF

$ kubectl get po busybox -n default -o jsonpath="{.metadata.labels}"
```



### Resources 분류

* Resources는 다음과 같이 4가지로 분류 가능

|No.  |Core   |Namespaced |Resources                                                  |
|:---:|:---:  |:---:      |---                                                        |
|1    |O      |O          |Pod, Service, PersistentVolumeClaim, ...                   |
|2    |O      |X          |Namespace, PersistentVolume, ...                           |
|3    |X      |O          |Deployment, DaemonSet, PodMetrics, Role, RoleBinding, ...  |
|4    |X      |X          |NodeMetrics, ClusterRole, ClusterRoleBinding, ...          |



### Core Resources APIs

|URL                                                                        |Method |설명                             |
|---                                                                        |---    |---                              |
|/raw/clusters/:cluster/api/:version/:resource                              |GET    |non-namespaced 리소스 목록 조회  |
|/raw/clusters/:cluster/api/:version/:resource:/:name                       |GET    |non-namespaced 리소스 조회       |
|/raw/clusters/:cluster/api/:version/:resource:/:name                       |DELETE |non-namespaced 리소스 삭제       |
|/raw/clusters/:cluster/api/:version/:resource:/:name                       |PATCH  |non-namespaced 리소스 수정       |
|/raw/clusters/:cluster/api/:version/namespaces/:resource:                  |GET    |namespaced 리소스 목록조회       |
|/raw/clusters/:cluster/api/:version/namespaces/:namespace/:resource/:name  |GET    |namespaced 리소스 조회           |
|/raw/clusters/:cluster/api/:version/namespaces/:namespace/:resource/:name  |DELETE |N\namespaced 리소스 삭제         |
|/raw/clusters/:cluster/api/:version/namespaces/:namespace/:resource/:name  |PATCH  |N\namespaced 리소스 수정         |


### apiGrouped Resource APIs

|URL                                                                                  |Method |설명                             |
|---                                                                                  |---    |---                              |
|/raw/clusters/:cluster/apis/:apiGroup/:version/:resource                             |GET    |non-namespaced 리소스 목록 조회  |
|/raw/clusters/:cluster/apis/:apiGroup/:version/:resource:/:name                      |GET    |non-namespaced 리소스 조회       |
|/raw/clusters/:cluster/apis/:apiGroup/:version/:resource:/:name                      |DELETE |non-namespaced 리소스 삭제       |
|/raw/clusters/:cluster/apis/:apiGroup/:version/:resource:/:name                      |PATCH  |non-namespaced 리소스 수정       |
|/raw/clusters/:cluster/apis/:apiGroup/:version/namespaces/:resource:                 |GET    |namespaced 리소스 목록조회       |
|/raw/clusters/:cluster/apis/:apiGroup/:version/namespaces/:namespace/:resource/:name |GET    |namespaced 리소스 조회           |
|/raw/clusters/:cluster/apis/:apiGroup/:version/namespaces/:namespace/:resource/:name |DELETE |namespaced 리소스 삭제           |
|/raw/clusters/:cluster/apis/:apiGroup/:version/namespaces/:namespace/:resource/:name |PATCH  |namespaced 리소스 수정           |

### CRUD examples

```
# Create
$ curl -X POST -H "Content-Type: application/json" http://localhost:3001/raw -d @- <<EOF
{
    "apiVersion": "v1",
    "kind": "Namespace",
    "metadata": {
        "name": "test-namespace"
    }
}
EOF

# Get
$ curl -X GET http://localhost:3001/raw/api/v1/namespaces/test-namespace

# Update
$ curl -X PATCH -H "Content-Type: application/merge-patch+json" http://localhost:3001/raw/api/v1/namespaces/test-namespace  -d @- <<EOF
{
    "metadata": {
        "labels": {
            "istio-injection": "disabled"
        }
    }
}
EOF

# verify
$ kubectl  get ns/test-namespace -o jsonpath={.metadata.labels.istio-injection}

# Delete
$ curl -X DELETE http://localhost:3001/raw/api/v1/namespaces/test-namespace


# List
$ curl -X GET http://localhost:3001/raw/api/v1/namespaces
```


### Debugging for docker container enviroment
* Dockerfile Build 
  * 디버깅용 backend 이미지 빌드 
    /src/app/backend/docker-compose -f dc-debug.yaml up --build
  * vscode/goland 등 에디터 원격 디버깅 설정(remote port 5555)
  * 참고링크 [컨테이너 내부 Go 애플리케이션 디버깅하기][https://mingrammer.com/debugging-containerized-go-app/]
  
