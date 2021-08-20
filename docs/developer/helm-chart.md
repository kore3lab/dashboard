# Packaging Helm-chart 

* validation

```
$ helm lint scripts/install/kuberntes/helm-chart
```

* output to yaml

```
$ helm template -n kore kore-board scripts/install/kuberntes/helm-chart \
  --set frontend.service.type=NodePort \
  --set frontend.service.nodePort=30080
```

* apply test

```
$ kubectl create ns kore
$ helm template --debug -n kore kore-board scripts/install/kuberntes/helm-chart \
  --set frontend.service.type=NodePort \
  --set frontend.service.nodePort=30080 \
  | kubectl apply -n kore -f -
```


* dry-run

```
$ helm install --dry-run -n kore kore-board scripts/install/kuberntes/helm-chart \
  --set frontend.service.type=NodePort \
  --set frontend.service.nodePort=30080
```

* packaging helm-chart

```
$ helm package scripts/install/kuberntes/helm-chart  -d scripts/install/kuberntes   # packaging  (tgz file 생성)
$ helm repo index scripts/install/kuberntes                                         # index.yaml 파일 생성/수정
```
