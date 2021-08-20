# Configure startup parameters

## Backend

* parameters

|이름                   |기본값                                               |설명                                                                                           |
|---                    |---                                                  |---                                                                                            |
|--log-level            |debug                                                |로그 레벨(trace, debug, info, warning, error, fatal, panic) https://github.com/sirupsen/logrus |
|--kubeconfig           |                                                     |kubeconfig 파일 위치                                                                           |
|--metrics-scraper-url  |http://localhost:8000                                |metrics-scraper api url                                                                        |
|--terminal-url         |http://localhost:3003                                |terminal api url                                                                               |
|--auth                 |strategy=cookie,secret=static-token,token=kore3lab   |인증처리방식 설정                                                                              |


* 환경변수 (env)

|이름                 |기본값                                               |설명                     |
|---                  |---                                                  |---                      |
|LOG_LEVEL            |debug                                                |"--log-level"            |
|KUBECONFIG           |                                                     |"--kubeconfig"           |
|METRICS_SCRAPER_URL  |http://localhost:8000                                |"--metrics-scraper-url"  |
|TERMINAL_URL         |http://localhost:3003                                |"--terminal-url "        |
|AUTH                 |strategy=cookie,secret=static-token,token=kore3lab   |"--auth"                 |


* Configuration of authentication
  * See [Sign-in configuration](./config-sign-in.md) page

* Configuration of kubeconfig
  * See [Kkuberconfig configuration](./config-kubeconfigs.md) page



## Metrics-Scraper

* parameters

|이름                 |기본값               |설명                       |
|---                  |---                  |---                        |
|--log-level          |                     |로그 레벨                  |
|--kubeconfig         |kubeconfig 파일 위치 |                           |
|--db-file            |/tmp/metrics.db      |sqllite database file path |
|--metric-resolution  |1m0s                 |metrics 수집 주기          |
|--metric-duration    |15m0s                |metrics 적산값 유지 기간   |
|--namespace          |                     |                           |

* 환경변수 (env)

|이름       |기본값 |설명                 |
|---        |---    |---                  |
|KUBECONFIG |       |kubeconfig 파일 위치 |


## Web-Terminal

* parameters

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