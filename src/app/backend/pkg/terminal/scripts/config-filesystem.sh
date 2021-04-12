#!/bin/bash
set -e

if [ "${WELCOME_BANNER}" ]; then
    echo ${WELCOME_BANNER}
fi

arg1=$1
arg2=$2
arg3=$3
arg4=$4
arg5=$5

# echo $arg1
# echo $arg2
# echo $arg3
# echo $arg4
# echo $arg5

#####################################
# 네임스페이스에 마운트할 정보 구성
#####################################

# 마운트 경로 구성
mkdir -p /nonexistent
mount -t tmpfs -o size=${SESSION_STORAGE_SIZE} tmpfs /nonexistent
cd /nonexistent

# 셀 정보 구성
cp /root/.bashrc ./
 echo 'source /opt/kubectl-aliases/.kubectl_aliases' >> .bashrc
echo 'PS1="> "' >> .bashrc

# kubectl을 위한 구성
mkdir -p .kube

export HOME=/nonexistent

echo $arg2| base64 -d > .kube/config

if [ ${KUBECTL_INSECURE_SKIP_TLS_VERIFY} == "true" ];then
    {
        clusters=`kubectl config get-clusters | tail -n +2`
        for s in ${clusters[@]}; do
            {
                echo `kubectl config set-cluster ${s} --insecure-skip-tls-verify=true` > /dev/null 2>&1
                echo `kubectl config unset clusters.${s}.certificate-authority-data` > /dev/null 2>&1
            } || {
                echo err > /dev/null 2>&1
            }
        done
    } || {
        echo err > /dev/null 2>&1
    }
fi

# 권한 구성
chown -R nobody:nogroup .kube

# 마운트 경로 export
export TMPDIR=/nonexistent

# 터미널로 처리될 최종 명령어 구성
# POD으로 연결할때
if [ "${arg1}" != "cluster" ];then
{
    #컨테이너 정보가 있으면 넘겨준다
    if [ -n "${arg5}" ] && [ "${arg5}" != "undefined" ];then
      container="--container ${arg5}"
    fi
    
    # exec su -s /bin/bash nobody -c "kubectl exec --stdin --tty ${arg3} ${container} -- /bin/bash" 
    /usr/bin/kubectl exec --kubeconfig .kube/config --stdin --tty --namespace=${arg3} ${arg4} ${container} -- /bin/sh
}
else
    exec su -s /bin/bash nobody
fi