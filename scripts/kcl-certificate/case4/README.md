# Scenario #4
> 라우팅 정책 반영 소요시간 (목표 : 평균 2초 이하)

* 단위서비스에 라우팅 정책을 반영한 후 재가동없이 실시간 평균 반영시간을 측정
* 정책 반영 성공율 100% 기준

## 측정방법

* 측정용 스크립트(어플리케이션)을 활용하여 측정 및 결과확인
* 시작시점 : 시험대상 서비스에서 현재 라우팅 서비스 확인하고 이에 따라 정책 A 또는 B를 반영한 후 해당 시점 측정
* 종료시점 : 시험대상 서비스에서 정책 반영에 따른 결과값 변경 여부 확인하고 해당 시점 측정
* 수행시간 = 종료시점 - 시작시점
* 총 1000회 측정하여 평균 반영시간을 계산


## Prerequites

```
$ kubectl apply -n bookinfo -f bookinfo/networking/destination-rule-all.yaml
$ kubectl apply -n bookinfo -f bookinfo/networking/virtual-service-all-v1.yaml

$ kubectl apply -n bookinfo -f - <<EOF
apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  name: bookinfo-gateway
spec:
  selector:
    istio: ingressgateway
  servers:
  - port:
      number: 80
      name: http
      protocol: HTTP
    hosts:
    - "*"
---
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: productpage
spec:
  hosts:
  - "*"
  gateways:
  - bookinfo-gateway
  http:
  - route:
    - destination:
        host: productpage
        subset: v1
EOF
```

## Policy A
> 리뷰 STAR 미 호출 정책 (STAR NOT_FOUND)


```
$ kubectl apply -n bookinfo -f - <<EOF
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: reviews
spec:
  hosts:
  - reviews
  http:
  - route:
    - destination:
        host: reviews
        subset: v1
EOF
```

## Policy B
> 리뷰 STAR 값 호출 정책 (STAR FOUND)

```
$ kubectl apply -n bookinfo -f - <<EOF
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: reviews
spec:
  hosts:
  - reviews
  http:
  - route:
    - destination:
        host: reviews
        subset: v3
EOF
```


## Verify
 
* http://101.55.69.105:32080/productpage
* CSS Selector : `body > div.container-fluid > div:nth-child(2) > div:nth-child(2) > blockquote:nth-child(2) > font > span:nth-child(1)` class 속성값으로 확인
* 공백(NOT_FOUND) 이면 정책 A 상태, 아니면 정책 B 상태
  


## Clean-up

```
$ kubectl delete -n bookinfo gateway/bookinfo-gateway
$ kubectl delete -n bookinfo -f bookinfo/networking/destination-rule-all.yaml
$ kubectl delete -n bookinfo -f bookinfo/networking/virtual-service-all-v1.yaml
```
