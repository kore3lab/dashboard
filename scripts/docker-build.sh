#!/bin/sh
# -----------------------------------------------------------------
# usage
if [ "$#" -lt 3 ]; then 
	echo "docker-build.sh <module name - backend/frontend/dashboard/all> <command> <tag>"
	echo "    ./docker-build.sh backend build v0.1.0"
	exit 0; 
fi

build() {
	if [[ "${MODULE}" == *"frontend"* ]];	then docker build --tag ${FRONTEND}:${TAG} ${ROOT};					fi
	if [[ "${MODULE}" == *"backend"* ]];	then docker build --tag ${BACKEND}:${TAG} ${ROOT}/src/app/backend;	fi
	if [[ "${MODULE}" == *"dashboard"* ]];	then docker build --tag ${DASHBOARD}:${TAG} ${ROOT}/dashboard;		fi
}
push() {
	if [[ "${MODULE}" == *"frontend"* ]];	then docker push ${FRONTEND}:${TAG};	fi
	if [[ "${MODULE}" == *"backend"* ]];	then docker push ${BACKEND}:${TAG};		fi
	if [[ "${MODULE}" == *"dashboard"* ]];	then docker push ${DASHBOARD}:${TAG};	fi
}
# latest() {
# 	if [[ "${MODULE}" == *"frontend"* ]];	then docker tag ${FRONTEND}:${TAG} ${FRONTEND}:latest;		docker push ${FRONTEND}:latest;		fi
# 	if [[ "${MODULE}" == *"backend"* ]];	then docker tag ${BACKEND:${TAG}} ${BACKEND}:latest;		docker push ${BACKEND}:latest;		fi
# 	if [[ "${MODULE}" == *"dashboard"* ]];	then docker tag ${DASHBOARD}:${TAG} ${DASHBOARD}:latest;	docker push ${DASHBOARD}:latest;	fi
# }

GROUP="acornsoftlab/"
PROJECT="acornsoft-dashboard"
ROOT="$(pwd)"
FRONTEND="${GROUP}${PROJECT}.frontend"
BACKEND="${GROUP}${PROJECT}.backend"
DASHBOARD="${GROUP}${PROJECT}.dashboard"

MODULE="$1"
COMMAND="$2"
TAG="$3"

echo "Module  : ${MODULE}"
echo "Command : ${COMMAND}"
echo "Tag     : ${TAG}"

if [[ "${COMMAND}" == *"build"* ]];		then build;		fi
if [[ "${COMMAND}" == *"push"* ]];		then push;		fi