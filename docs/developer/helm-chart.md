# Packaging Helm-chart 

* validation

```
$ helm lint scripts/install/kubernetes/helm-chart
```

* output to yaml

```
$ helm template -n kore kore-board scripts/install/kubernetes/helm-chart \
  --set frontend.service.type=NodePort \
  --set frontend.service.nodePort=30080
```

* apply test

```
$ kubectl create ns kore
$ helm template --debug -n kore kore-board scripts/install/kubernetes/helm-chart \
  --set frontend.service.type=NodePort \
  --set frontend.service.nodePort=30080 \
  | kubectl apply -n kore -f -
```


* dry-run

```
$ helm install --dry-run -n kore kore-board scripts/install/kubernetes/helm-chart \
  --set frontend.service.type=NodePort \
  --set frontend.service.nodePort=30080
```

* packaging helm-chart

```
$ helm package scripts/install/kubernetes/helm-chart  -d scripts/install/kubernetes   # packaging  (tgz file 생성)
$ helm repo index scripts/install/kubernetes                                         # index.yaml 파일 생성/수정
```
