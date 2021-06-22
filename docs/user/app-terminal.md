# Web-Terminal

## Run

* Arguments

|이름                   |기본값 |설명                                                                                       |
|---                    |---    |---                                                                                        |
|--kubeconfig           |       |kubeconfig 파일 위치                                                                       |
|--log-level            |debug  |로그 레벨(panic,fatal,error,warning,info,debug,trace) https://github.com/sirupsen/logrus)  |
|--corsonoff            |on     |CORS(Cross-Origin Resource Sharing) on/off (defaults to on(blocked by CORS))               |


* 환경변수 (env)

|이름       |기본값 |설명                 |
|---        |---    |---                  |
|KUBECONFIG |       |kubeconfig 파일 위치 |


### 실행 여부 확인

* check healthy

```
$ curl http://localhost:3003/healthy

"healthy"
```

* check websocket
  * "HTTP/1.1 101" 리턴 여부 확인

```
$ curl --include --no-buffer \
  --header "Connection: Upgrade" \
  --header "Upgrade: websocket" \
  --header "Sec-WebSocket-Key: x2pbk53raexyotmhf7ug" \
  --header "Sec-WebSocket-Version: 13" \
  http://localhost:3003/api/terminal/ws

HTTP/1.1 101 Switching Protocols
....
```

## API

* 아래 설명에서 사용하는 변수는 다음과 같습니다.
  * `clusters` : Kubeconfig context name
  * `namespaces` : Resource namespace
  * `pods` : Pod name
  * `containers` : Container name
  * `termtype` : terminal type(cluster/pod/container) 

* prefix : /api/terminal/clusters/{CLUSTER}

|URL Pattern                                                                  |Method |설명                                 |
|---                                                                          |---    |---                                  |
|termtype/{TERMTYPE}                                                          |GET    |Web terminal 접속토큰 요청(kubectl)  |
|namespaces/{NAMESPACE}/pods/{POD}/termtype/{TERMTYPE}                        |GET    |Web terminal 접속토큰 요청(pod)      |
|namespaces/{NAMESPACE}/pods/{POD}/containers/{CONTAINER}/termtype/{TERMTYPE} |GET    |Web terminal 접속토큰 요청(container)|

* anothers

|URL Pattern      |Method |설명                                 |
|---              |---    |---                                  |
|/api/terminal/ws |GET    |Web terminal websocket 접속요청      |
|/api/v1/config   |PATCH  |kubeconfig refresh event from backend|


## 클라이언트 개발

### 개발용 서버 컨테이너 실행
* kubeconfig 는 `${HOME}/.kube/config` 적용

```
$ docker run --rm -d --privileged -p 3003:3003 --name terminal \
    -v "${HOME}/.kube:/app/.kube"\
    ghcr.io/acornsoftlab/kore-board.terminal:latest --kubeconfig=/app/.kube/config --corsonoff=off
```

### 클라이언트 활용 (예제 코드)

* 명령어 `kubectl config view` 를 실행시키는 간단한 클라이언트 예제
* kubeconfig context 는 'apps-06' 로 가정
* 주요 구현 로직
  * Token 조회
  * socket connection 오픈
  * Token 전달
  * Screen 사이즈 전달
  * 명령어(`kubectl config view`) 전달
  * 수신 데이터 출력
* 아래 예제 코드 참조
```
<!DOCTYPE html>
<head>
<script src="https://unpkg.com/axios/dist/axios.min.js"></script>
<script>
  let context = "apps-06"

  axios.get(`http://localhost:3003/api/terminal/clusters/${context}/termtype/cluster`)  // get a token
    .then(resp=> {
      let token = resp.data.Token;

      let socket = new WebSocket("ws://localhost:3003/api/terminal/ws");
      socket.onopen = (e) => {
        console.log("[open]");
        socket.send(JSON.stringify({ Arguments: "", AuthToken: token,})); // authentication (using a token)
        socket.send("3" + JSON.stringify({ columns: 100, rows: 100 }));   // set screnn size
        socket.send("1kubectl config view\n");                            // execute shell
      };
      socket.onmessage = (e) =>  {
        if(e.data[0] == "1") {
          document.write(atob(e.data.slice(1)).replaceAll("\n","<br>"));  // receive data 
        }
      };
    })
</script>
</head>
</html>
```
