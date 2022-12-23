# Authenticate API access with serviceaccount

if `client-certificate` is not available on KUBECONFIG


## Prerequisites

* `kubectl` and runnable

## How-to

* Install kore-board
  * See [Installation](./installation.md) page.

* Create a secret for serviceaccount `kore-board`
```
$ kubectl apply -f - <<EOF
apiVersion: v1
kind: Secret
metadata:
  name: kore-board-secret
  namespace: kore
  annotations:
    kubernetes.io/service-account.name: kore-board
type: kubernetes.io/service-account-token
EOF
```

* Grant `cluster-admin` permission to the serviceaccount

```
$ kubectl create clusterrolebinding kore-board-binding --clusterrole=cluster-admin --serviceaccount=kore:kore-board
```

* Create and verify the KUBECONFIG file  (ca.crt, token)

```
$ echo -e "apiVersion: v1
kind: Config
clusters:
- cluster:
    certificate-authority-data: $(kubectl get secret kore-board-secret -n kore -o jsonpath='{.data.ca\.crt}')
    server: $(kubectl config view -o jsonpath="{.clusters[?(@.name==\"$(kubectl config view -o jsonpath="{.contexts[?(@.name==\"$(kubectl config current-context)\")].context.cluster}")\")].cluster.server}")
  name: token-cluster
contexts:
- context:
    cluster: token-cluster
    user: token-user
  name: admin
current-context: token
users:
- name: token-user
  user:
    token: $(kubectl get secret kore-board-secret -n kore -o jsonpath='{.data.token}' | base64 --decode)
" > kubeconfig-token.yaml

$ kubectl get sa kore-board -n kore --kubeconfig="$(pwd)/kubeconfig-token.yaml"
$ kubectl get secret kore-board-secret -n kore --kubeconfig="$(pwd)/kubeconfig-token.yaml"
$ kubectl get nodes --kubeconfig="$(pwd)/kubeconfig-token.yaml"
```

### GKE Autopilot

* Configure kubconfig with `gcloud container` command

```
$ export KUBECONFIG="kubeconfig-gke-autopilot.yaml"
$ gcloud container clusters get-credentials autopilot-cluster-1 --region asia-northeast3 --project kore-project
$ kubectl get nodes
```

* Set gcloud config when edit **permission error** occurs

```
ERROR: (gcloud.container.clusters.get-credentials) get-credentials requires edit permission on ....
```

```
$ gcloud config set container/use_client_certificate False
```

* [Authenticate API access with serviceaccount](#how-to)
