#!/bin/bash

URL="http://101.55.69.105:32080"

init() {

	if [ "${CASE}" == "1" ]; then
		kubectl delete -n sock-shop -f https://raw.githubusercontent.com/itnpeople/k8s.docs/master/demo/yaml/sock-shop.yaml
	fi
	if [ "${CASE}" == "2" ]; then
		kubectl create ns sock-shop
		kubectl label namespace sock-shop istio-injection=enabled
		kubectl apply -n sock-shop -f https://raw.githubusercontent.com/itnpeople/k8s.docs/master/demo/yaml/sock-shop.yaml
		kubectl apply -n sock-shop -f ${NAME}/networking.yaml
		kubectl get po -n sock-shop -w
	fi
	if [ "${CASE}" == "4" ] || [ "${CASE}" == "5" ]; then
		kubectl apply -n bookinfo -f bookinfo/networking/destination-rule-all.yaml
		kubectl apply -n bookinfo -f bookinfo/networking/virtual-service-all-v1.yaml
		kubectl apply -n bookinfo -f ${NAME}/networking.yaml
		kubectl delete -n bookinfo authorizationpolicy productpage
	fi
}

verify() {
	if [ "${CASE}" == "2" ]; then
		curl -s -H "HOST: sock-shop.acornsoft.io"  ${URL}/  -o /dev/null -w "code:%{http_code}\n"
	fi
}

dev() {

	if [ "${CASE}" == "2" ] || [ "${CASE}" == "4" ] || [ "${CASE}" == "5" ]; then
		jmeter -t ${ROOT}/${NAME}/${NAME}.jmx
	fi
	if [ "${CASE}" == "3" ]; then
		export NAMESPACE=latency
		python3 ${ROOT}/${NAME}/runner/runner.py --conn 8 --qps 10 --duration 100 --bothsidecar
	fi

}

exec() {

	if [ "${CASE}" == "2" ] || [ "${CASE}" == "4" ] || [ "${CASE}" == "5" ]; then
		rm ${ROOT}/${NAME}/${NAME}.jtl
		jmeter -n -t ${ROOT}/${NAME}/${NAME}.jmx -l ${ROOT}/${NAME}/${NAME}.jtl
	fi
	if [ "${CASE}" == "3" ]; then
		export NAMESPACE=latency
		python3 ${ROOT}/${NAME}/runner/runner.py --conn 16 --qps 100 --duration 100 --bothsidecar
	fi

}


clean() {
	# if [ "${CASE}" == "2" ]; then
	# 	kubectl delete -n sock-shop -f ${NAME}/networking.yaml
	# 	kubectl delete -n sock-shop -f https://raw.githubusercontent.com/itnpeople/k8s.docs/master/demo/yaml/sock-shop.yaml
	# fi
	if [ "${CASE}" == "4" ] || [ "${CASE}" == "5" ]; then
		kubectl delete -n bookinfo -f bookinfo/networking/destination-rule-all.yaml
		kubectl delete -n bookinfo -f bookinfo/networking/virtual-service-all-v1.yaml
		kubectl delete -n bookinfo -f ${NAME}/networking.yaml
	fi
}

result() {
	if [ "${CASE}" == "2" ] || [ "${CASE}" == "4" ] || [ "${CASE}" == "5" ]; then
		jmeter -t ${ROOT}/result.jmx
	fi
}


# ------------------------------------------------------------------------------
if [ "$#" -lt 2 ]; then
	echo "not enough parameters (size=$#)"
	exit 0; 
fi

export JAVA_HOME="/Library/Java/JavaVirtualMachines/jdk1.8.0_201.jdk/Contents/Home"
export JVM_ARGS="-Xms1024m -Xmx1024m"
export ROOT="$(pwd)"
export CASE="$1"
export NAME="case${CASE}"

if [ "$2" == "init" ];		then init;		fi
if [ "$2" == "dev" ];		then dev;		fi
if [ "$2" == "exec" ];		then exec;		fi
if [ "$2" == "clean" ];		then clean;		fi
if [ "$2" == "verify" ];	then verify;	fi
if [ "$2" == "result" ]; 	then result;	fi
