#!/bin/bash
set -e

all=$*

if [ -z "${all}" ]; then
    echo No Args provided
    echo Terminal will exit.
    sleep 1
    exit 1
fi

if [[ $all == ERROR:* ]]; then
    echo ${all}
    sleep 1
    exit 1
fi

# Set terminal type 
echo "export TERM=xterm" >> /root/.bashrc

# Set completion
echo "source /usr/share/bash-completion/bash_completion" >> /root/.bashrc
echo 'source <(kubectl completion bash)' >> /root/.bashrc
echo 'complete -F __start_kubectl k' >> /root/.bashrc

# Mount Namespace 구성 (https://man7.org/linux/man-pages/man1/unshare.1.html)
# --fork : 지정된 프로그램을 직접 실행이 아닌 자식 프로세스로 포크해서 처리
# --pid : PID 기준 네임스페이스 파일을 지정하면 바인딩 마운트를 통해서 영구적인 네임스페이스 생성, 영구 네임스페이스 구성시 --fork 옵션을 지정해야 한다.
# --mount-proc : 프로그램을 시작하기 전에 /proc 파일 시스템을 마운트한다.
# --mount : 지정한 파일의 내용으로 마운트 네임스페이스 생성
unshare --fork --pid --mount-proc --mount /opt/k3webterminal/config-filesystem.sh ${all}