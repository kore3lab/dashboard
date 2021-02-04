---
layout: post
slug: "k3lab-dashboard-frontend-readme"
title: "K3lab Dashboard 프론트엔드 설명"
description: "K3lab Dashboard 프론트엔드 설명"
categories: [Vue.js, BootStrapVue]
tags: [Vue.js, BootStrapVue]
redirect_from:
  - /2021/2/01/
---


> **용어 및 개념 정리 정리**
> - **Micin**
> 믹스인은 **특정 기능(행위)만을 담당하는 클래스**로, 단독 사용(Stand-alone use)이 아닌 다른 클래스에 탑재되어 사용될 목적으로 작성된 **(조각) 클래스**를 의미합니다.

# K3lab Dashboard 프론트엔드

## 디렉토리 구성

Nuxt JS 의 기본 구조를 따른다.

- /assets
  - css, image, font와 같은 리소스들을 포함한다.
  
- /components
  - 애플리케이션에서 사용될 컴포넌트들을 포함하며 해당 경로에 위치된 컴포넌트들은 Nuxt.js의 비동기 데이터 함수인 asyncData또는 fetch를 사용할 수 없다.

- /layouts
  - 애플리케이션 전체에 대한 레이아웃을 포함한다. 기본으로 default.vue가 생성되어 있을 것이고 상황에 맞게 layout을 생성할 수 있다. 해당 디렉토리는 이름을 변경할 수 없다.
  - ./components
    - 레이아웃을 구성하는 컴포넌트 이다. (예: footer.vue, aside.vue ...)
  
- /pages
  - 실제 애플리케이션의 페이지 구성을 포함하며 해당 디렉토리의 구조에 따라 router가 자동으로 생성된다. 해당 디렉토리는 이름을 변경할 수 없다.

- /plugins
  - 애플리케이션에 바인딩 될 외부 혹은 내부 plugins를 포함한다.plugins는 애플리케이션이 인스턴스 화 되기 전에 실행하며 전역적으로 구성 요소를 등록하고 함수 또는 상수를 삽입할 수 있다.
  - mix

- /static
  - 해당 디렉토리는 정적인 파일들을 포함한다. 구성에 따라서 html, Javascript 파일도 포함 시킬 수 있다. 해당 디렉토리는 이름을 변경할 수 없다.

- /store
  - 애플리케이션에서 사용될 vuex store 파일들을 포함한다. 기본적으로 비활성화 상태이고 store 디렉토리에 index.js 파일을 작성하면 store가 활성화된다. 구성에 따라서 모듈 형태의 store를 형성할 수 있다. 해당 디렉토리는 이름을 변경할 수 없다.

## 현재 진행 상황

- Kubernetes Dashboard의 backend를 제거
- 표준 정의를 위한 샘플 페이지 생성 - pod.list.vue
- 나머지 페이지는 표준 정의가 된 후 값 대입 예정
- pull request 해보기

## Git Workflow - fork & feature-branch work flow 사용

- 1. 중앙 레파지토리를 포크 한다. - 메인 프로젝트
  - github > fork
  
- 2. 개발자 레파지토리로 포크된 프로젝트를 clone 한다.
  - git clone <GitHub_Username> / Project_Name
    ```text
    Clone 명령어를 통해 (아래의 명령어를 포함한 작업이 수행된다.)
    1. 디렉토리를 만들고
    2. 디렉토리로 들어가고 나서 git init 명령으로 빈 Git 저장소를 만든다.
    3. 입력한 URL을 origin이라는(기본값)􀀁이름의 리모트로 추가하고(git􀀁remote􀀁add)
    4. git􀀁fetch 명령으로 리모트 저장소에서 데이터를 가져온다.
    ```

- 3. 중앙 원격을 upstream(관용적 사용)이라는 이름의 리모트로 추가한다.
  - git remote add upstream [중앙 원격 url]
  - git remote add origin [내 원격 url] - clone 하면 기봊적으로 생성 된다.
  - git remote -v // 설정한 remote 목록을 확인

- 4. 내 로컬에서 새로운 작업을 하는 경우 master or develop 에서 그 작업에 대한 브랜치를 생성 후 이동한다.
  - git branch feature/login      //새로운 branch 생성
  - git checkout feature/login    // 해당 branch로 작업 위치 이동
  - git checkout -b feature/login // 위의 두 명령으로 합한 것

- 5. 나의 작업이 끝난 경우(해당 기능을 모두 구현한 경우) feature/login branch를 origin remote로 push한다.
  - commit 절차를 간략하게 정리후 push 한다. - 개발단계의 불필요한 commit 단계를 정리 한다.
    ```sh
    # git log 확인
    $ git log --pretty=oneline

    # git reset 사용 - squash all
    $ git reset $(git commit-tree HEAD^{tree} -m "A new start")

    # git rebase 사용 - log값을 얼마 없을때 유용
    # edit 화면에서 pick 부분을 squash로 변경 하면 합쳐진다.
    $ git rebase -i HEAD~2  

    # merge를 이용한 방법 - develop 에서 개발하고 staging 브랜치가 있을때 유용 한다.
    # develop 브랜치에서 feature/featureA 브랜치를 merge 합니다. 
    # 이 때 --squash 옵션을 주어 여러 커밋을 하나의 커밋으로 모아서 develop 브랜치에 새로 커밋을 만들어줍니다. 
    # 이렇게 하는 이유는 develop 브랜치에 한 feature 당 하나의 커밋, 한 bugfix 당 하나의 커밋을 유지하여 이후 staging 브랜치에 선택적으로 커밋을 옮길 때 수월하게 구분 할 수 있게 하기 위해서입니다.
    $ git checkout develop
    $ git merge --squash feature/featureA
    Updating cb0ab0e..0d9dd2e
    Fast-forward
    Squash commit -- not updating HEAD
      test.txt | 3 ++-
      1 file changed, 2 insertions(+), 1 deletion(-)
    $ git commit -m 'Merged in bugfix/remove-unused-code (pull request #1)’
    $ git push origin develop
    ```
  - git push 하면 origin remote(포크한 리파지토리)에도 똑같은 feature/login branch가 생긴다.

  - 6. pull request 요청을 한다.

- Git 커밋 메세지 정리
- 변경 사항을 원격으로 푸시하기 전에 항상 스쿼시하십시오.
  
  ```sh
  $ git rebase -i HEAD~4

  # 에디터 창에서 rebase 대상 커밋이 순서대로 표시됨(예전 커밋이 위에 있음)
  pick abcd
  pick 1234
  pick 34ef
  pick 87ad

  # 맨 위 커밋은 pick 으로 남기고 다른 커밋은 squash 나 fixup 으로 수정
  pick abcd
  squash 1234
  squash 34ef
  squash 87ad
  ```

## 표준 정의에 대한 생각

- 주석 처리 : JSDoc 기준으로 주석을 달아 준다.
  - 참고: [주석 - JAVASCRIPT.INFO](https://ko.javascript.info/comments)
  
- 메소드 이름 정책에 대한 생각
  - 카멜표기법을 사용 한다.
  - 메소드명 시작은 메소드의 성격을 나타낸다.
    - on : 액션
    - to : 변환 / 전달
    - get : 조회
      - 주어진 식별자에 대해 데이터를 가져 오는 데 사용됩니다. 
      - 주어진 식별자에 대한 개체를 찾을 수 없으면 null 개체를 반환합니다. 
      - 객체가 있는지 확실하지 않으면 get () 메서드를 사용하십시오. 
    - load : 조회
      - 주어진 식별자에 대해 데이터를 가져 오는 데 사용됩니다. 
      - 개체를 찾을 수 없음 예외가 발생합니다. 
      - 항상 프록시 객체를 반환하므로이 메서드는 객체를 지연로드합니다.
      - 객체가 존재한다고 확신하는 경우 load () 메서드를 사용하십시오. 

## 공통 함수 목록

/plugins/mixin-common-methods.js 에서 정의 한다.

- getElapsedTime(timestamp)
  ```js
  /**
   * timestamp를 day,hour,minute,second로 구분 봔환함
   *
   * @param {date} timestamp 변환할 date 값
   * @return {string} timestamp의 day/hour/minute/second 값으로 변환하여 반환함
   */
  ```

- getFormatMetrics(resource, metrics, decimals = 2)
  ```js
  /**
  * 리소스 메트릭 수집 값을 Formatting 한다.
  * CPU / milliCPU 또는 CPU의 1/1000 단위로 처리 할 때 "표현"됩니다. 나노 코어 / 나노 CPU는 CPU의 1/1000000000 (10 억분의 1)입니다.
  * Memory 단위는 이진접두어를 사용 한다.
  *
  * @param {string} resource 구분자 cpu/memory
  * @param {number} metrics 리소스 사용량 합계 값
  * @param {number} decimals 소수점 아래 자릿수. 기본값은 2 이다
  * @return {string} 리소스 합산 값의 단위를 변환한다.
  */
  ```

## 레이아웃 목록(/layouts)

### ./components/navbar.vue

네이스페이스 값을 드랍다운 형식으로 보여 준다.

- 페이지 메소드 목록
  - loadNamespaces()
    ```js
    /**
		 * 클러스터의 네임스페이스 목록을 가져 온다.
		 * 
		 * @async
		 * @function loadNamespaces
		 * @returns {Array} 네이스페이스 값을 [{valu, text}] 값으로 반환 한다.
		 */
    ```

## 페이지 목록

기본 컨셉은 kubectl로 볼 수 있는 값들을 항목으로 보여 준다.
리소스(cpu/memory) 사용량을 항목으로 보여 준다. 

### pod.list.vue

클러스터의 파드 목록을 리스트 한다.

- 리스트 항목
  ```js
  { key: "name", label: "이름", sortable: true },
  { key: "namespace", label: "네임스페이스", sortable: true  },
  { key: "ready", label: "Ready", sortable: true  },
  { key: "status", label: "상태", sortable: true  },
  { key: "restartCount", label: "재시작", sortable: true  },
  { key: "creationTimestamp", label: "생성시간" },
  { key: "nodeName", label: "노드", sortable: true  },
  { key: "usageCpu", label: "CPU 사용량", sortable: true  },
  { key: "usageMemory", label: "MEMORY 사용량", sortable: true  },
  ```

- kubuctl 정보
```sh
$ kubectl get pods -o wide
NAME                                READY   STATUS    RESTARTS   AGE     IP               NODE           NOMINATED NODE   READINESS GATES
nginx-deployment-7b74d6486c-km8ch   2/2     Running   0          25h     10.244.239.197   k8s-node-2-1   <none>           <none>
nginx-deployment-7b74d6486c-qgfj2   2/2     Running   0          25h     10.244.239.198   k8s-node-2-1   <none>           <none>
nginx-deployment-7b74d6486c-wp7w7   2/2     Running   0          25h     10.244.78.70     k8s-node-1-1   <none>           <none>
test-local-vol                      1/1     Running   14         3d22h   10.244.239.193   k8s-node-2-1   <none>           <none>
test-local-vol-1-1                  1/1     Running   13         3d22h   10.244.78.67     k8s-node-1-1   <none>           <none>
test-local-vol-1-2                  1/1     Running   13         3d22h   10.244.78.68     k8s-node-1-1   <none>           <none>
test-local-vol-1-3                  0/1     Pending   0          3d22h   <none>           <none>         <none>           <none>
test-local-vol-1-31                 0/1     Pending   0          3d22h   <none>           <none>         <none>           <none>

-- 메트릭 서버 api 사용
$ kubectl get --raw "/apis/metrics.k8s.io/v1beta1/namespaces//pods"

-- k8s apiserver api
$ kubectl get --raw "/api/v1/namespaces/default/pods/"
```

- 페이지 메소드 목록
  - query_All()
    ```js
    /**
     * 파드 리스트를 가져 와서 각 항목 값을 반환 한다.
     * 
     * @returns {object} 파드 리스트 항목 별로 값을 반환 한다.
     */
    ```

  - getMetrics()
    ```js
    /**
		 * 리소스 메트릭 값을 반환 한다.
		 * 
		 * @async
		 * @function getMetrics
		 * @returns {Promise<object>} 리소스(cpu/memory) 메트릭 값을 반환 한다.
		 */
    ```

  - toUsageHandler(resource, podName)
    ```js
    /**
		 * 메트릭 값을 리소스별로 단위 계산 후 반환 한다.
		 * 
		 * @param {string} resource 구분자 cpu/memory
		 * @param {string} podName 구분자 pod 이름으로 구분 한다.
		 * @return {string} 리소스 합산 값의 단위를 추가 해서 반환 한다.
		 */
    ```

# 참고 자료

- [믹스인 - JAVASCRIPT.INFO](https://ko.javascript.info/mixins)
- [주석 - JAVASCRIPT.INFO](https://ko.javascript.info/comments)
