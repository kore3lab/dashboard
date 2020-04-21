# Scenario #5
> 보안 인증 반영 소요시간 (목표 : 평균 2초 이하)

## 측정방법
* 측정용 스크립트(어플리케이션)을 활용하여 측정 및 결과확인
* 시작시점 : 시험대상 서비스로 http 호출 Status Code 200 여부 확인하고 보안인증정책 적용한 후 해당 시점 측정
* 종료시점 : 시험대상 서비스로 http 호출 Status Code 403 리턴 시점 측정
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
---
apiVersion: security.istio.io/v1beta1
kind: AuthorizationPolicy
metadata:
  name: productpage
  namespace: bookinfo
spec:
  selector:
    matchLabels:
      app: productpage
  action: ALLOW
  rules:
  - to:
    - operation:
        paths: ["/productpage"]
EOF
```

## Allow
> Apply allow authorization policy


```
$ kubectl apply -n bookinfo -f - <<EOF
apiVersion: security.istio.io/v1beta1
kind: AuthorizationPolicy
metadata:
  name: productpage
  namespace: bookinfo
spec:
  selector:
    matchLabels:
      app: productpage
  action: ALLOW
  rules:
  - to:
    - operation:
        paths: ["/productpage"]
EOF
```


## Deny
> Apply deny authorization policy


```
$ kubectl apply -n bookinfo -f - <<EOF
apiVersion: security.istio.io/v1beta1
kind: AuthorizationPolicy
metadata:
  name: productpage
  namespace: bookinfo
spec:
  selector:
    matchLabels:
      app: productpage
  action: DENY
  rules:
  - to:
    - operation:
        paths: ["/productpage"]
EOF
```

## Verify

* 200, 403 (forbidden)

```
$ curl -s http://101.55.69.105:32080/productpage -o /dev/null -w "%{http_code}"
```


## Clean-up

```
$ kubectl delete -n bookinfo authorizationpolicy/productpage
$ kubectl delete -n bookinfo gateway/bookinfo-gateway
$ kubectl delete -n bookinfo -f bookinfo/networking/destination-rule-all.yaml
$ kubectl delete -n bookinfo -f bookinfo/networking/virtual-service-all-v1.yaml
```

