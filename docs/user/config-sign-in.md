# Sign-in Configuration

## Introduction

* strategy
  * cookie : managed login true/false cookie
  * local : using JWT mechanism, access-token and refresh-token issuing local,  access-key & refresh-key properties required

* supported secret
  * static-user : compare static username/password
  * static-token : compare static token string
  * basic-auth : compare kubernetes's basic-auth secret ( username, password )
  * service-account-token : compare kubernetes's service-account-token secret

* login schema
  * user : username, password
  * token : string


## How to apply
  * apply the feature as a startup parameter.

```
auth=strategy=<cookie/local>,secret=<supported secret>,<prop1=value1>,<prop2=value2>
```

* example

```
kind: Deployment
apiVersion: apps/v1
metadata:
  labels:
    app: kore-board
    kore.board: backend
  name: backend
  namespace: kore
spec:
...
    spec:
      containers:
        - name: backend
          image: ghcr.io/acornsoftlab/kore-board.backend:latest
          args:
            - --metrics-scraper-url=http://metrics-scraper:8000
            - --log-level=info
            - --auth=strategy=cookie,secret=static-token,token=acornsoft
```

### static-token

* define static token string

```
spec:
  containers:
    - name: backend
      image: ghcr.io/acornsoftlab/kore-board.backend:latest
      args:
        - --auth=strategy=<cookie/local>,secret=static-token,token=<token-string>
```


### static-user

* define static username, password

```
spec:
  containers:
    - name: backend
      image: ghcr.io/acornsoftlab/kore-board.backend:latest
      args:
        - --auth=strategy=<cookie/local>,secret=static-user,username=<username>,password=<password>
```

### basic-auth

* create a 'basic-auth' secret
```
$ kubectl apply -f - <<EOF
apiVersion: v1
kind: Secret
metadata:
  name: secret-basic-auth
type: kubernetes.io/basic-auth
stringData:
  username: admin
  password: t0p-Secret
EOF
```

* using volumn mount

```
spec:
  containers:
    - name: backend
      image: ghcr.io/acornsoftlab/kore-board.backend:latest
      args:
        - --auth=strategy=<cookie/local>,secret=basic-auth,dir=/var/user
        ...
      volumeMounts:
      - name: user-vol
        mountPath: "/var/user"
    volumes:
    - name: user-vol
      secret:
        secretName: secret-basic-auth
```

### service-account-token 

```
spec:
  containers:
    - name: backend
      image: ghcr.io/acornsoftlab/kore-board.backend:latest
      args:
        - --auth=strategy=<cookie/local>,access-key=<access-token-secret>,refresh=<refresh-token-secret>,secret=service-account-token
```

* get a token-string (ex. query from serviceaccount `kore-board`)

```
$ SECRET="$(kubectl get sa -n kore -l app=kore-board -o jsonpath='{.items[0].secrets[0].name}')"
$ echo "$(kubectl get secret ${SECRET} -n kore -o jsonpath='{.data.token}' | base64 --decode)"
```

* get a token-string (ex. query from current namespace's serviceaccount `default`)

```
$ SECRET="$(kubectl get sa default -o jsonpath='{..secrets[0].name}')"
$ echo "$(kubectl get secret ${SECRET} -o jsonpath='{.data.token}' | base64 --decode)"
```

* input token-string in your browser login page
