# kubeconfig Configuration

## Introduction

* strategy (kubeconfig-provider)
  * file : update kubeconfig as file
  * configmap : update kubeconfig as configmap


## How to apply
  * apply the feature as a startup parameter.

```
kubeconfig=
kubeconfig=/app/.kube/config
kubeconfig=strategy=configmap,configmap=kore-board-kubeconfig,namespace=kore,filename=config
```


### file-strategy
> using volumne (persistent or)

```
spec:
  containers:
    - name: backend
      image: ghcr.io/acornsoftlab/kore-board.backend:latest
      args:
        - --kubeconfig=/app/.kube/config
      volumeMounts:
        - mountPath: /app/.kube
          name: kubeconfig-volume
    volumes:
      - name: kubeconfig-volume
```


### configmap-strategy
> using configmap

* create a configmap

```
$ kubectl create configmap kore-board-kubeconfig  -n kore --from-file=config=${HOME}/.kube/config
```


```
spec:
  containers:
    - name: backend
      image: ghcr.io/acornsoftlab/kore-board.backend:latest
      args:
        - --kubeconfig=strategy=configmap,configmap=kore-board-kubeconfig,namespace=kore,filename=config
```

