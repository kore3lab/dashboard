#!/bin/bash
set -e

if [ "${WELCOME_BANNER}" ]; then
    echo ${WELCOME_BANNER}
fi

# get options:
while (( "$#" )); do
    case "$1" in
        --inclustermode)
            if [ -n "$2" ] && [ ${2:0:1} != "-" ]; then
                ARG_INCLUSTER_MODE=$2
                shift 2
            else
                echo "Error: Argument for $1 is missing" >&2
                exit 1
            fi
            ;;
        --kubeconfig)
            if [ -n "$2" ] && [ ${2:0:1} != "-" ]; then
                ARG_KUBECONFIG=$2
                shift 2
            else
                echo "Error: Argument for $1 is missing" >&2
                exit 1
            fi
            ;;
        --kubetoken)
            if [ -n "$2" ] && [ ${2:0:1} != "-" ]; then
                ARG_KUBETOKEN=$2
                shift 2
            else
                echo "Error: Argument for $1 is missing" >&2
                exit 1
            fi
            ;;   
        --termtype)
            if [ -n "$2" ] && [ ${2:0:1} != "-" ]; then
                ARG_TERM_TYPE=$2
                shift 2
            else
                echo "Error: Argument for $1 is missing" >&2
                exit 1
            fi
            ;;  
        --container)
            if [ -n "$2" ] && [ ${2:0:1} != "-" ]; then
                ARG_CONTAINER=$2
                shift 2
            else
                echo "Error: Argument for $1 is missing" >&2
                exit 1
            fi
            ;; 
        --pod)
            if [ -n "$2" ] && [ ${2:0:1} != "-" ]; then
                ARG_POD=$2
                shift 2
            else
                echo "Error: Argument for $1 is missing" >&2
                exit 1
            fi
            ;;
        --namespace)
            if [ -n "$2" ] && [ ${2:0:1} != "-" ]; then
                ARG_NAMESPACE=$2
                shift 2
            else
                echo "Error: Argument for $1 is missing" >&2
                exit 1
            fi
            ;; 
        --cluster)
            if [ -n "$2" ] && [ ${2:0:1} != "-" ]; then
                CLUSTER=$2
                shift 2
            else
                echo "Error: Argument for $1 is missing" >&2
                exit 1
            fi
            ;;       
        --help)
            echo "Usage:  $0 -i <inclustermode> [options]" >&2
            echo "        -i | --inclustermode  %  (set inclustermode to ...)" >&2
            echo "        -k | --kubeconfig  %  (set kubeconfig of ...)" >&2
            echo "        -t | --kubetoken     (set token of ...)" >&2
            exit 0
            ;;
        --*) # unsupported flags
            echo "Error: Unsupported flag: $1" >&2
            echo "$0 -h for help message" >&2
            exit 1
            ;;
        *)
            echo "Error: Arguments with not proper flag: $1" >&2
            echo "$0 -h for help message" >&2
            exit 1
            ;;
    esac
done
# echo "===parsed command line option==="
# echo " - inclustermode: ${ARG_INCLUSTER_MODE}"
# echo " - kubeconfig: ${ARG_KUBECONFIG}"
# echo " - kubetoken: ${ARG_KUBETOKEN}"
# echo " - termmode: ${ARG_TERM_TYPE}"
# echo " - container: ${ARG_CONTAINER}"
# echo " - pod: ${ARG_POD}"
# echo " - namespace: ${ARG_NAMESPACE}"



#####################################
# 네임스페이스에 마운트할 정보 구성
#####################################

# 마운트 경로 구성
mkdir -p /nonexistent
mount -t tmpfs -o size=${SESSION_STORAGE_SIZE} tmpfs /nonexistent
cd /nonexistent

# 셀 정보 구성
#cp /root/.bashrc ./
# Set terminal type 
echo "export TERM=xterm" >> .bashrc

# Set completion
echo "source /usr/share/bash-completion/bash_completion" >> .bashrc
echo 'source <(kubectl completion bash)' >> .bashrc
echo 'complete -F __start_kubectl k' >> .bashrc
echo 'source /opt/kubectl-aliases/.kubectl_aliases' >> .bashrc
echo 'PS1="kubectl@dashboard> "' >> .bashrc

# kubectl을 위한 구성
mkdir -p .kube

export HOME=/nonexistent

if [ "${ARG_INCLUSTER_MODE}" == "true" ]; then
    echo `kubectl config set-credentials webkubectl-user --token=${ARG_KUBETOKEN}` > /dev/null 2>&1
    echo `kubectl config set-cluster kubernetes --server=${ARG_KUBECONFIG}` > /dev/null 2>&1
    echo `kubectl config set-context kubernetes --cluster=kubernetes --user=webkubectl-user` > /dev/null 2>&1
    echo `kubectl config use-context kubernetes` > /dev/null 2>&1
else
    echo $ARG_KUBECONFIG| base64 -d > .kube/config
fi

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
if [ "${ARG_TERM_TYPE}" == "cluster" ];then
    exec su -s /bin/bash nobody
elif [ "${ARG_TERM_TYPE}" == "pod" ];then
    /usr/bin/kubectl exec --kubeconfig .kube/config --stdin --tty --namespace=${ARG_NAMESPACE} ${ARG_POD} -- /bin/bash || /usr/bin/kubectl exec --kubeconfig .kube/config --stdin --tty --namespace=${ARG_NAMESPACE} ${ARG_POD} -- /bin/sh || echo "remote shell is not supported"
elif [ "${ARG_TERM_TYPE}" == "container" ];then
    /usr/bin/kubectl exec --kubeconfig .kube/config --stdin --tty --namespace=${ARG_NAMESPACE} ${ARG_POD} --container ${ARG_CONTAINER} -- /bin/bash || /usr/bin/kubectl exec --kubeconfig .kube/config --stdin --tty --namespace=${ARG_NAMESPACE} ${ARG_POD} --container ${ARG_CONTAINER} -- /bin/sh || echo "remote shell is not supported"
else
    echo "term type argument error"
    exit 1    
fi