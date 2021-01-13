# Install on Kubernetes


## Install Istio 

* Istio v1.7.3
* Addons :Promethus, Kiali
* https://istio.io/v1.7/docs/setup/install/istioctl/
* Ingress-gateway를 `NodePort:32080 ` 로 오픈하고 Kiali와 연동


### Istio download

```
$ curl -L https://istio.io/downloadIstio | ISTIO_VERSION=1.7.3 TARGET_ARCH=x86_64 sh -
$ cd istio-1.7.3
$ chmod +x bin/istioctl
$ cp bin/istioctl /usr/local/bin
```

### Install Istio

* 아래와 같이 default profile 에서 수정할 ingressGateways 스펙을 조회 (해당 스펙 참조)
```
$ istioctl profile dump default --config-path components.ingressGateways
$ istioctl profile dump default --config-path values.gateways.istio-ingressgateway
```

* install - Ingress-gateway 를 NodePort (32080)로  override 설치

```
$ istioctl install --set profile=default -f - <<EOF
apiVersion: install.istio.io/v1alpha1
kind: IstioOperator
spec:
  components:
    ingressGateways:
    - enabled: true
      k8s:
        service:
          ports:
          - name: status-port
            port: 15021
            targetPort: 15021
          - name: http2
            port: 80
            targetPort: 8080
            nodePort: 32080
          - name: https
            port: 443
            targetPort: 8443
            nodePort: 32443
          - name: tls
            port: 15443
            targetPort: 15443
      name: istio-ingressgateway
  values:
    gateways:
      istio-ingressgateway:
        type: NodePort
EOF
```


### Install Addons : Promethus & Kiali
* https://kiali.io/documentation/latest/quick-start/

* Install Promethus & Kiali 
```
$ kubectl apply -n istio-system -f samples/addons/prometheus.yaml
$ kubectl apply -n istio-system -f samples/addons/kiali.yaml
```

* Expose Kiali UI
```
$ kubectl apply -n istio-system -f - <<EOF
apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  name: default-gateway
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
  name: kiali
spec:
  hosts:
  - "*"
  gateways:
  - default-gateway
  http:
  - match:
    - uri:
        prefix: /kiali
    route:
    - destination:
        port:
          number: 20001
        host: kiali
EOF
```

* Promethus Admin API Enable
  * https://prometheus.io/docs/prometheus/latest/querying/api/
  * `--web.enable-admin-api` 옵션 추가


## Install Acoronsoft-dashboard

* Careate configmap `acornsoft-dashboard-kubeconfig`

```
$ kubectl create configmap acornsoft-dashboard-kubeconfig --from-file=${HOME}/.kube/config -n acornsoft-dashboard

# if update configmap
$ kubectl create configmap acornsoft-dashboard-kubeconfig --from-file=${HOME}/.kube/config --dry-run -o yaml | kubectl apply  -n acornsoft-dashboard -f -
```


* Deploy

```
$ kubectl apply -f kuberntes/recommended.yaml
```


### NodePort
* Frontend: 30080
* Backend : 30081
* Dashboard backend : 30090

### Verify

* Web UI : http://<server>:30080/
* Kiali UI : http://<server>:32080/kiali
* Kiali embedding UI (kiosk mode) : http://<server>:32080/kiali/console/graph/namespaces/?kiosk=true


## Install Acoronsoft-dashboard - "in-cluster" mode
> Install for Single Cluster

```
# create namespace

$ NAMESPACE="dashboard"
$ kubectl create ns ${NAMESPACE}


# install metrics-scraper

$ kubectl run metrics-scraper -n ${NAMESPACE}\
  --image=ghcr.io/acornsoftlab/acornsoft-dashboard.metrics-scraper:latest --port=8000\
  -- --db-file=metrics.db
$ kubectl expose pod metrics-scraper -n ${NAMESPACE} --port=8000 --name=metrics-scraper


# install backend

$ kubectl run backend -n ${NAMESPACE}\
  --image=ghcr.io/acornsoftlab/acornsoft-dashboard.backend:latest --port=3001
$ kubectl expose pod backend -n ${NAMESPACE} --name=backend --type='NodePort' --port=3001


# install kubernetes-dashoard

$ kubectl create role acornsoft-dashboard -n ${NAMESPACE} --resource=* --verb=*
$ kubectl create rolebinding acornsoft-dashboard -n ${NAMESPACE} --role=acornsoft-dashboard --serviceaccount=${NAMESPACE}:default
$ kubectl create secret generic kubernetes-dashboard-csrf -n ${NAMESPACE} --from-literal="csrf="
$ kubectl create clusterrolebinding acornsoft-dashboard --clusterrole=cluster-admin --serviceaccount=${NAMESPACE}:default

$ kubectl run dashboard -n ${NAMESPACE}\
  --image=ghcr.io/acornsoftlab/acornsoft-dashboard.dashboard:latest --port=9090\
  -- --namespace=${NAMESPACE} --sidecar-host=http://metrics-scraper:8000

$ kubectl expose pod dashboard -n ${NAMESPACE}  --name=dashboard --type='NodePort' --port=9090


# install frontend

$ BACKEND_PORT="$(kubectl get svc/backend -n ${NAMESPACE} -o jsonpath="{.spec.ports[0].nodePort}")"
$ DASHBOARD_PORT="$(kubectl get svc/dashboard -n ${NAMESPACE} -o jsonpath="{.spec.ports[0].nodePort}")"

$ kubectl run frontend -n ${NAMESPACE} --image=ghcr.io/acornsoftlab/acornsoft-dashboard.frontend:latest  --port=3000\
  --env="BACKEND_PORT=${BACKEND_PORT}"\
  --env="DASHBOARD_PORT=${DASHBOARD_PORT}"
$ kubectl expose pod frontend -n ${NAMESPACE} --name=frontend --type='NodePort' --port=3000


# open in your browser

$ echo "http://<end-point ip>:$(kubectl get svc/frontend -n ${NAMESPACE} -o jsonpath="{.spec.ports[0].nodePort}")"
```

* Clean-up

```
$ kubectl delete -n ${NAMESPACE} pod/backend pod/dashboard pod/frontend pod/metrics-scraper
$ kubectl delete -n ${NAMESPACE} service/backend service/dashboard service/frontend service/metrics-scraper
$ kubectl delete -n ${NAMESPACE} role/acornsoft-dashboard rolebinding/acornsoft-dashboard
$ kubectl delete -n ${NAMESPACE} secret/kubernetes-dashboard-csrf
$ kubectl delete clusterrolebinding/acornsoft-dashboard
$ kubectl delete ns ${NAMESPACE}
```


