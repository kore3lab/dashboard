# Scenario #3
> 프록시 경유 지연시간 (목표 : 50ms 이하)

* https://github.com/istio/tools/tree/release-1.7/perf/benchmark

* benchmark tool clone

```
$ git clone https://github.com/istio/tools.git
$ git checkout release-1.7
$ cd tools/perf/benchmark
```

* 어플리케이션 설치

```
export DNS_DOMAIN=local
export NAMESPACE=latency
export INTERCEPTION_MODE=REDIRECT
export ISTIO_INJECT=true
./setup_test.sh
```

* python 환경 구성

```
$ cd runner
$ pipenv --three
$ sudo pip3 install pandas bokeh
$ cd ..
```

* 실행 
  * 초당 100회 간격(qps) 으로 100초 동안(duration) 16개 쓰레드(conn)로 테스트 실행
  * 총 1000회

```
$ python3 ./runner/runner.py --conn 16 --qps 100 --duration 100 --bothsidecar
```

* 결과 예제

```
-------------- Running in both mode --------------
kubectl --namespace latency exec fortioclient-7f7f6bf785-9dgj5  -- fortio load  -jitter=False -c 10 -qps 10 -t 100s -a -r 0.00005   -httpbufferkb=128 -labels 52ea5c6d_qps_10_c_10_1024_mixer_both http://fortioserver:8080/echo?size=1024
Defaulting container name to captured.
Use 'kubectl describe pod/fortioclient-7f7f6bf785-9dgj5 -n latency' to see all of the containers in this pod.
02:05:50 I fortio_main.go:167> Not using dynamic flag watching (use -config to set watch directory)
Fortio 1.7.1 running at 10 queries per second, 8->8 procs, for 1m40s: http://fortioserver:8080/echo?size=1024
02:05:50 I httprunner.go:82> Starting http test for http://fortioserver:8080/echo?size=1024 with 10 threads at 10.0 qps
Starting at 10 qps with 10 thread(s) [gomax 8] for 1m40s : 100 calls each (total 1000)
02:07:30 I periodic.go:558> T008 ended after 1m40.001536143s : 100 calls. qps=0.99998463880597
02:07:30 I periodic.go:558> T003 ended after 1m40.002269183s : 100 calls. qps=0.9999773086849075
02:07:30 I periodic.go:558> T009 ended after 1m40.002948007s : 100 calls. qps=0.9999705207990489
02:07:30 I periodic.go:558> T005 ended after 1m40.002958677s : 100 calls. qps=0.9999704141053511
02:07:30 I periodic.go:558> T007 ended after 1m40.002955923s : 100 calls. qps=0.9999704416437223
02:07:30 I periodic.go:558> T001 ended after 1m40.002985765s : 100 calls. qps=0.9999701432414526
02:07:30 I periodic.go:558> T006 ended after 1m40.002993135s : 100 calls. qps=0.9999700695458589
02:07:30 I periodic.go:558> T000 ended after 1m40.002998682s : 100 calls. qps=0.9999700140791824
02:07:30 I periodic.go:558> T002 ended after 1m40.002967317s : 100 calls. qps=0.9999703277104709
02:07:30 I periodic.go:558> T004 ended after 1m40.003005942s : 100 calls. qps=0.9999699414835416
Ended after 1m40.003033098s : 1000 calls. qps=9.9997
Sleep times : count 990 avg 1.0074259 +/- 0.0005228 min 1.005650874 max 1.009220424 sum 997.351662
Aggregated Function Time : count 1000 avg 0.0024579484 +/- 0.0005064 min 0.000691148 max 0.004311082 sum 2.45794836
# range, mid point, percentile, count
>= 0.000691148 <= 0.0007 , 0.000695574 , 0.10, 1
> 0.0009 <= 0.001 , 0.00095 , 0.20, 1
> 0.001 <= 0.00125 , 0.001125 , 2.70, 25
> 0.00125 <= 0.0015 , 0.001375 , 9.70, 70
> 0.0015 <= 0.00175 , 0.001625 , 10.10, 4
> 0.00175 <= 0.002 , 0.001875 , 15.30, 52
> 0.002 <= 0.00225 , 0.002125 , 21.90, 66
> 0.00225 <= 0.0025 , 0.002375 , 39.70, 178
> 0.0025 <= 0.003 , 0.00275 , 95.30, 556
> 0.003 <= 0.0035 , 0.00325 , 97.80, 25
> 0.0035 <= 0.004 , 0.00375 , 99.20, 14
> 0.004 <= 0.00431108 , 0.00415554 , 100.00, 8
# target 50% 0.00259263
# target 75% 0.00281745
# target 90% 0.00295234
# target 99% 0.00392857
# target 99.9% 0.0042722
Sockets used: 10 (for perfect keepalive, would be 10)
Jitter: false
Code 200 : 1000 (100.0 %)
Response Header Sizes : count 1000 avg 249 +/- 0 min 249 max 249 sum 249000
Response Body/Total Sizes : count 1000 avg 1273 +/- 0 min 1273 max 1273 sum 1273000
All done 1000 calls (plus 10 warmup) 2.458 ms avg, 10.0 qps
Successfully wrote 3015 bytes of Json data to 2020-11-12-020550_52ea5c6d_qps_10_c_10_1024_mixer_both.json
```