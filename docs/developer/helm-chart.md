# Helm-chart development guide

* output yaml

```
$ helm template --debug -n kore kore-board ./kuberntes/helm-chart/ \
  --set frontend.service.type=NodePort \
  --set frontend.service.nodePort=30080
```

* apply test

```
$ kubectl create ns kore
$ helm template --debug -n kore kore-board ./kuberntes/helm-chart/ \
  --set frontend.service.type=NodePort \
  --set frontend.service.nodePort=30080 \
  | kubectl apply -n kore -f -
```


* dry-run

```
$ helm install --dry-run --debug -n kore kore-board ./kuberntes/helm-chart/ \
  --set frontend.service.type=NodePort \
  --set frontend.service.nodePort=30080
```



* Packaging helm-chart

```
$ helm package helm-chart/    # packaging  (tgz file 생성)
$ helm repo index .           # index.yaml 파일 생성/수정
```
