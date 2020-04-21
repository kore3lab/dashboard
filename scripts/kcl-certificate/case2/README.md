# Scenario #2
> 장애발생 분석추적 성공률 (목표 : 99% 이상)

## 측정방법
* 시험대상 서비스에 임의의 장애를 주입하고 장애 분석결과를 모니터링하여 성공여부 확인
* "장애제거-요청발생-결과확인-장애주입-요청발생-결과확인" 6단계를 통해 실험을 진행
  1. 장애제거 : 장애 주입을 제거
  1. 요청발생 : 클라이언트 요청을 발생 3.결과확인 : 장애분석 조회 스크립트를 통해 장애 없음을 확인
  1. 장애주입 : 장애발생 클라이언트 스크립트 실행
  1. 요청발생 : 클라이언트 요청을 발생 6.결과확인 : 장애분석 조회 스크립트를 실행하여 장애발생 여부 확인
* 장애주입단계부터 결과확인단계까지 10초이내 장애결과 조회여부에 따라 성공/실패 판단
* 총 100회 측정하여 평균 성공률을 계산



## Prerequites

```
$ kubectl label namespace sock-shop istio-injection-
$ kubectl label namespace sock-shop istio-injection=enabled
$ kubectl apply -n sock-shop -f https://raw.githubusercontent.com/itnpeople/k8s.docs/master/demo/yaml/sock-shop.yaml

$ kubectl apply -n sock-shop -f - <<EOF
apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  name: sock-shop-gateway
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
  name: sock-shop
spec:
  hosts:
  - "sock-shop.acornsoft.io"
  gateways:
  - sock-shop-gateway
  http:
  - route:
    - destination:
        host: front-end
        port:
          number: 80
EOF
```

```
$ curl -H "HOST: sock-shop.acornsoft.io"  http://101.55.69.105:32080/
```


## Status OK
> 정상상태

```
$ kubectl apply -n sock-shop -f - <<EOF
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: sock-shop
spec:
  hosts:
  - "sock-shop.acornsoft.io"
  gateways:
  - sock-shop-gateway
  http:
  - route:
    - destination:
        host: front-end
        port:
          number: 80
EOF
```


## Status Abort
> abort injection 주입


```
$ kubectl apply -n sock-shop -f - <<EOF
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: sock-shop
spec:
  hosts:
  - "sock-shop.acornsoft.io"
  gateways:
  - sock-shop-gateway
  http:
  - fault:
      abort:
        httpStatus: 503
        percentage:
          value: 100
    route:
    - destination:
        host: front-end
        port:
          number: 80
EOF
```


## Clean-up

```
$ kubectl delete -n bookinfo gateway/bookinfo-gateway
$ kubectl delete -n bookinfo vs/productpage
$ kubectl delete -n bookinfo -f bookinfo/networking/destination-rule-all.yaml
```
