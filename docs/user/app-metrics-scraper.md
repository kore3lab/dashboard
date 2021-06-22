# Metrics-Scraper

## Run

* Arguments

|이름                 |기본값               |설명                       |
|---                  |---                  |---                        |
|--log-level          |                     |로그 레벨                  |
|--kubeconfig         |kubeconfig 파일 위치 |                           |
|--db-file            |/tmp/metrics.db      |sqllite database file path |
|--metric-resolution  |1m0s                 |metrics 수집 주기          |
|--metric-duration    |15m0s                |metrics 적산값 유지 기간   |
|--namespace          |                     |                           |

* 환경변수 (env)

|이름       |기본값 |설명                 |
|---        |---    |---                  |
|KUBECONFIG |       |kubeconfig 파일 위치 |


## API

|URL Pattern                                                                  |Method |설명                               |
|---                                                                          |---    |---                                |
|/api/v1/clusters/:cluster                                                    |GET    |클러스터 summary metrics  조회     |
|/api/v1/clusters/:cluster/nodes/:node/metrics/:metrics                       |GET    |클러스터 Node metrics 조회         |
|/api/v1/nodes/:node/metrics/:metrics                                         |GET    |default 클러스터 노드 metrics 조회 |
|/api/v1/clusters/:cluster/namespaces/:namespaces/pods/:pod/metrics/:metrics  |GET    |클러스터 Pod metrics 조회          |
|/api/v1/namespaces/:namespaces/pods/:pod/metrics/:metrics                    |GET    |default 클러스터 Pod metrics 조회  |

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

