# Contributor Guide

## Branch 전략

### Git-Flow
![Basic Git Flow For Making Open Source Contributions on GitHub](https://dnncommunity.org/DesktopModules/Blog/BlogImage.ashx?TabId=65&ModuleId=454&Blog=1&Post=1470&w=1140&h=400&c=0&key=289a2e46-efbd-471c-830d-ccfdd93d46ea)

* [Basic Git Flow For Making Open Source Contributions on GitHub](https://dnncommunity.org/blogs/Post/1470/Basic-Git-Flow-For-Making-Open-Source-Contributions-on-GitHub)
* [우린 Git-flow를 사용하고 있어요](https://woowabros.github.io/experience/2017/10/30/baemin-mobile-git-branch-strategy.html)

### Summary

* 개발은 원본(acornosftlab/dashboard) repository를 fork 받아 로컬 repository 에서 수행하고 "Pull Reuqest"를 통해 리뷰 프로세스 수행
* 리뷰가 완료되면 원본(acornosftlab/dashboard) repository의 branch에 merge 하여 반영
* Feature 개발 : `develop` branch에서 로컬 branch를 생성&개발 후 원본 repository의 `develop` branch에 merge
* Hotfix 개발 : `master` branch에서 로컬 branch를 생성&개발 후 원본 repository의 `develop` 와 `master` branch에 merge
* Release 개발 : `develop` branch에서 로컬 branch를 생성&개발 후 원본 repository의 `masater` 와 `develop` branch에 merge
* Release/Hotfix 완료 후 태깅처리하여 릴리즈 출시
* 태깅 및 `master` branch로의 merge 작업은 관리자(owner) 권한 사용자가 수행

## Contributing

### Fork & Clone

* 자신의 repository에 "Fork" 받는다.

```
$ git clone https://github.com/<github id>/dashboard.git
$ cd dashboard/

$ git remote -v

origin	https://github.com/<github id>/dashboard.git (fetch)
origin	https://github.com/<github id>/dashboard.git (push)
```

원격 저장소(프로젝트의 원래 저장소)를 추가한다.

```
$ git remote add upstream https://github.com/acornsoftlab/dashboard.git

$ git remote -v

origin	https://github.com/<github id>/dashboard.git (fetch)
origin	https://github.com/<github id>/dashboard.git (push)
upstream	https://github.com/acornsoftlab/dashboard.git (fetch)
upstream	https://github.com/acornsoftlab/dashboard.git (push)
```

### 개발 준비

`develop` branch로 이동
```
$ git checkout develop

$ git branch

* develop
  master
```

* 리모트 최신소스를 forked repository `develop` branch와 동기화 (rebase)
```
$  git fetch upstream
$  git rebase upstream/develop
```


* 작업할 개발 branch를 생성하고, 해당 branch로 이동

```
$ git branch <branch name>
$ git checkout <branch name>

$ git branch 

  develop
* <branch name>
  master
```

 * [README.md](https://github.com/acornsoftlab/dashboard/blob/master/README.md)에서 "Getting started" 섹션의 clone 이후 내용을 참고하여 개발환경 구성



### 개발 완료 후 Pull Request 생성
> [Pull Request](https://github.com/itnpeople/k8s.docs/blob/master/git.md) 섹션 참조


* **필수** 개발된 로컬 branch의 commit을 업무에 맞도록 정리(squash) 작업 수행
  * 수월한 충돌 해결을 미연에 방지 효과
  * [커밋 합치기(squash)](https://github.com/itnpeople/k8s.docs/blob/master/git.md) 섹션 참조

* **필수** 리모트 최신소스를 forked repository `develop` branch와 동기화(rebase)

```
$ git fetch upstream
$ git rebase upstream/develop
```

* 충돌이 발생할 경우 충돌 해결 

* noConflict 이후 개발 branch에 force 푸쉬
```
$ git push origin +<branch name>
```

* 브라우저 에서 `http://github.com/<github id>/dashboard` 을 열고 "Pull Request" 생성

### Commit message convention (Pull Request)

* 커밋 메시지는 아래와 같이 `구분 : 제목` 형태로 작성, 이슈번호가 없을 경우는 `(#이슈번호)`는 생략

```
구분 : 제목 (#이슈번호)
```

* 구분
  * feat: a new feature
  * fix: a bug fix
  * docs: changes to documentation
  * style: formatting, etc; no code change
  * refactor: refactoring production code
  * test: adding tests, refactoring test; no production code change
  * chore: updating build tasks, package manager configs, etc; no production code change
