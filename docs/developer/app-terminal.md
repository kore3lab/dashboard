# Web-Terminal

## 개요

### 전체 주요 프로세스

![termianl process](http://www.plantuml.com/plantuml/proxy?cache=no&src=https://raw.githubusercontent.com/acornsoftlab/dashboard/master/docs/developer/app-terminal.puml)

### Linux `unshare` 명령을 활용

* 네임스페이스 단위로 고립 처리
* unshare `--fork, --proc` 옵션을 통해서 메인 프로세스를 기준으로 서브 프로세스 (local command) 단위 네임스페이스 생성
* unshare `--mount-proc --mount` 옵션을 통해서 프로세스의 네임스페이스로 마운트될 파일 시스템 구성 (kubectl 등이 설치된 파일 시스템)


## 개발

* 개발 제약 사항
  * 리눅스 `unshare` 명령을 활용하기 때문에 컨테이너 환경에서만 동작
  * "소스 개발 → 컨테이너 빌드 및 실행 → 디버깅 → 소스 수정 → 컨테이너 빌드 및 실행 → ..." 절차로 개발
* 컨테이너 빌드 및 실행은 docker-compose 활용
* 디버깅은 VSCode의 "Docker:Debug in Container"을 활용 (컨테이너 리모트 연결)


### 컨테이너 디버깅 환경 설정
> using VSCode

1. "VSCode → Main Menu → Run → Add Configuration → Docker:Debug in Container" 선택
1. 아래와 같이 입력

```
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch remote",
      "type": "go",
      "request": "attach",
      "mode": "remote",
      "remotePath": "",
      "port": 5555,
      "host": "127.0.0.1",
      "showLog": true,
      "trace": "log"
    }
  ]
}
```

### 이미지 빌드 및 실행
> using docker-compose

```
 $ docker-compose -f src/app/terminal/dc-debug.yaml up --build
```

### 디버깅 (활용 예시)

* `./src/app/terminal/terminal-main.go` 파일 `func healthy(w http.ResponseWriter, r *http.Request) {` 에 break-point 지정
* 디버깅 시작 
  * VSCode → Main Menu → Run → Debug 선택
  * 또는 "F5" 키 

*  콘솔을 열고 아래와 같이 break-point 이벤트 발생

```
$ curl http://localhost:3003/healthy
```

* VSCode에서 지정된 break-point 로 포커싱되는 것을 확인
* "F10", "F11" 키를 활용하여 디버깅 실행


