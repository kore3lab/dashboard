# Public Certificate (KCL)

## Introduction
* Kubernetes v1.16.13
* Istio v1.7.3
* [Install Kubernetes & Istio ](../install/)

### Scenario
1. 초당 최대 토폴로지 시각화 Element 수
1. 장애발생 분석추적 성공률 (목표 : 99% 이상)
1. 프록시 경유 지연시간 (목표 : 50ms 이하)
1. 라우팅 정책 반영 소요시간 (목표 : 평균 2초 이하)
1. 보안 인증 반영 소요시간 (목표 : 평균 2초 이하)

### Node Port

```
acornsoftlab-kore3   backend                     NodePort       10.103.5.120     <none>        3001:30081/
acornsoftlab-kore3   dashboard                   NodePort       10.106.39.66     <none>        9090:30090/
acornsoftlab-kore3   frontend                    NodePort       10.96.142.69     <none>        3000:30080/
guestbook            frontend                    NodePort       10.97.113.100    <none>        80:30594/
istio-system         istio-ingressgateway        NodePort       10.101.168.108   <none>        15021:31033/TCP,80:32080/TCP,443:32443/TCP,15443:30742/
istio-system         prometheus                  NodePort       10.99.226.2      <none>        9090:32090/
sock-shop            front-end                   NodePort       10.105.226.230   <none>        80:30001/
```

## Test

```
$ ./run.sh <testcase no> init    # 준비
$ ./run.sh <testcase no> exec    # 시작
$ ./run.sh <testcase no> result  # 결과
$ ./run.sh <testcase no> clean   # Clean-up
$ ./run.sh <testcase no> dev     # 개발화면
```

|시험                                 |목표         |init   |exec |result |clean  |참고                         |               |
|---                                  |---          |:-----:|:---:|:-----:|:-----:|---                          |---            |
|초당 최대 토폴로지 시각화 Element 수 |30 ea/s 이상 |O      |     |       |       |http://101.55.69.105:30080/  |web            |
|장애발생 분석추적 성공률             |99% 이상     |O      |O    |O      |       |kubectl get po -n sock-shop  |jmeter, 25분   |
|프록시 경유 지연시간                 |50ms 이하    |       |O    |       |       |kubectl get po -n latency    |fortio, 2분    |
|라우팅 정책 반영 소요시간            |2s 이하      |O      |O    |O      |O      |kubectl get po -n bookinfo   |jmeter, 4~6분  |
|보안 인증 반영 소요시간              |2s 이하      |O      |O    |O      |O      |kubectl get po -n bookinfo   |jmeter, 7~12분 |
