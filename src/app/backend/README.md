# Backend

## Raw

* [Kubernetes API Concepts](https://kubernetes.io/docs/reference/using-api/api-concepts/)


### 샘플 호출 예제

```
# Root URL 정의
$ export RootUrl="http://localhost:3001/api/_raw/$(k config current-context)"


# core API
$ curl ${RootUrl}//v1/pods
$ curl ${RootUrl}//v1/namespaces 
$ curl ${RootUrl}//v1/nodes
$ curl ${RootUrl}/apps/v1/deployments 

# core API (in namespace)
$ curl ${RootUrl}//v1/namespaces/kube-system/pods
$ curl ${RootUrl}/apps/v1/namespaces/kube-system/deployments

# CRD 
$ curl ${RootUrl}/networking.istio.io/v1alpha3/virtualservices
$ curl ${RootUrl}/networking.istio.io/v1alpha3/namespaces/bookinfo/virtualservices
```

* apiGroups
```
$ kubectl api-resources -o wide

$ kubectl get crd
$ kubectl get crd virtualservices.networking.istio.io -o jsonpath="{.spec.group}"
```
