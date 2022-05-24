# API Specification

## Backend
> kubernetes@in-cluster : in-cluster name
> kubernetes-admin@kubernetes : kubeconfig default context name


|URL Pattern                        |Method |설명                                     |
|---                                |---    |---                                      |
|/api/clusters                      |GET    |k8s cluster context 리스트 조회          |
|/api/clusters/:cluster/topology    |GET    |토플로지 그래프 조회                     |
|/api/clusters/:cluster/dashboard   |GET    |Dashboard 데이터 조회                    |

* Examples

```
$ curl -X GET http://localhost:3001/api/contexts
$ curl -X GET http://localhost:3001/api/clusters/kubernetes@in-cluster/graph/topology
$ curl -X GET http://localhost:3001/api/clusters/kubernetes@in-cluster/dashboard
```


### Backend kubernetes Raw API
> 멀티 클러스터를 지원하는 Kubernetes API Proxy API


#### Kubernetes API
* [Kubernetes API Concepts](https://kubernetes.io/docs/reference/using-api/api-concepts/)
* [OepnAPI spec.](https://raw.githubusercontent.com/kubernetes/kubernetes/master/api/openapi-spec/swagger.json)
* Kubernetes API 에서 제공하는 resource 와 resource의 api-group은 `kubectl api-resources -o wide` 으로 resource와 apiGroup 조회 가능

```
$ kubectl api-resources -o wide

# CRD 경우
$ kubectl get crd
$ kubectl get crd virtualservices.networking.istio.io -o jsonpath="{.spec.group}"
```

#### URL Pattern

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


#### Apply APIs
> Create a resource

|URL Pattern            |Method |설명                     |
|---                    |---    |---                      |
|/raw/clusters/:cluster |POST   |Applay                   |

* Example

```
$ curl -X POST -H "Content-Type: application/json" http://localhost:3001/raw/clusters/kubernetes@in-cluster -d @- <<EOF
{
    "apiVersion": "v1",
    "kind": "Namespace",
    "metadata": {
        "name": "test-namespace"
    }
}
EOF
```

#### Update APIs
> Update a resource

|URL Pattern            |Method |설명                     |비고             |
|---                    |---    |---                      |---              |
|/raw/clusters/:cluster |PATCH  |Update                   |                 |


* Reuqest Header `Content-Type` 으로 patch 방식 선택 (
  * patch 방식에 대한 이해 - [Update API Objects in Place Using kubectl patch](https://kubernetes.io/docs/tasks/manage-kubernetes-objects/update-api-object-kubectl-patch/)
  * [JSON Merge Patch (RFC 7396)](https://tools.ietf.org/html/rfc7386)  : `Content-Type : application/merge-patch+json`
  * [JSON Patch (RFC 6902)](https://tools.ietf.org/html/rfc6902) : `Content-Type : application/json-patch+json`
  * Strategic merge patch : `Content-Type : application/strategic-merge-patch+json`


* JSON merge patch example

```
$ curl -X PATCH -H "Content-Type: application/merge-patch+json" http://localhost:3001/raw/clusters/kubernetes@in-cluster/api/v1/namespaces/default/pods/busybox -d @- <<EOF
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
$ curl -X PATCH -H "Content-Type: application/json-patch+json" http://localhost:3001/raw/clusters/kubernetes@in-cluster/api/v1/namespaces/default/pods/busybox -d @- <<EOF
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



#### Resources 분류

* Resources는 다음과 같이 4가지로 분류 가능

|No.  |Core   |Namespaced |Resources                                                  |
|:---:|:---:  |:---:      |---                                                        |
|1    |O      |O          |Pod, Service, PersistentVolumeClaim, ...                   |
|2    |O      |X          |Namespace, PersistentVolume, ...                           |
|3    |X      |O          |Deployment, DaemonSet, PodMetrics, Role, RoleBinding, ...  |
|4    |X      |X          |NodeMetrics, ClusterRole, ClusterRoleBinding, ...          |



#### Core Resources APIs

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


#### apiGrouped Resource APIs

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

#### CRUD examples

```
# Create
$ curl -X POST -H "Content-Type: application/json" http://localhost:3001/raw/clusters/kubernetes@in-cluster -d @- <<EOF
{
    "apiVersion": "v1",
    "kind": "Namespace",
    "metadata": {
        "name": "test-namespace"
    }
}
EOF

# Get
$ curl -X GET http://localhost:3001/raw/clusters/kubernetes@in-cluster/api/v1/namespaces/test-namespace

# Update
$ curl -X PATCH -H "Content-Type: application/merge-patch+json" http://localhost:3001/raw/clusters/kubernetes@in-cluster/api/v1/namespaces/test-namespace  -d @- <<EOF
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
$ curl -X DELETE http://localhost:3001/raw/clusters/kubernetes@in-cluster/api/v1/namespaces/test-namespace

# List
$ curl -X GET http://localhost:3001/raw/clusters/kubernetes@in-cluster/api/v1/namespaces
```


## Metrics-Scraper

|URL Pattern                                                 |Method |설명                       |
|---                                                         |---    |---                        |
|/api/v1/clusters/:cluster                                   |GET    | summary metrics  조회     |
|/api/v1/clusters/:cluster/nodes/:node                       |GET    | Node metrics 조회         |
|/api/v1/clusters/:cluster/namespaces/:namespaces/pods/:pod  |GET    | Pod metrics 조회          |

* 변수
  * `:cluster` : Kubeconfig context name
  * `:node` :  Node name
  * `:metrics` : `cpu` or `memory`
  * `:pod` : Pod name

* Examples

```
$ curl -X GET http://localhost:8000/api/v1/clusters/kubernetes@in-cluster
$ curl -X GET http://localhost:8000/api/v1/clusters/kubernetes@in-cluster/nodes/vm-live-01
$ curl -X GET http://localhost:8000/api/v1/clusters/kubernetes@in-cluster/namespaces/default/pods/busybox
```


## Web-Terminal

* 아래 설명에서 사용하는 변수는 다음과 같습니다.
  * `clusters` : Kubeconfig context name
  * `namespaces` : Resource namespace
  * `pods` : Pod name
  * `containers` : Container name
  * `termtype` : terminal type(cluster/pod/container) 

* prefix : /api/terminal/clusters/{CLUSTER}

|URL Pattern                                                                  |Method |설명                                 |
|---                                                                          |---    |---                                  |
|termtype/{TERMTYPE}                                                          |GET    |Web terminal 접속토큰 요청(kubectl)  |
|namespaces/{NAMESPACE}/pods/{POD}/termtype/{TERMTYPE}                        |GET    |Web terminal 접속토큰 요청(pod)      |
|namespaces/{NAMESPACE}/pods/{POD}/containers/{CONTAINER}/termtype/{TERMTYPE} |GET    |Web terminal 접속토큰 요청(container)|

* anothers

|URL Pattern      |Method |설명                                 |
|---              |---    |---                                  |
|/api/terminal/ws |GET    |Web terminal websocket 접속요청      |
|/api/v1/config   |PATCH  |kubeconfig refresh event from backend|


