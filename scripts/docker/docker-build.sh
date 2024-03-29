#!/bin/sh
# -----------------------------------------------------------------
# usage
if [ "$#" -lt 3 ]; then 
	echo "docker-build.sh <module name - backend/frontend/all> <command> <tag>"
	echo "    ./docker-build.sh backend build v0.1.0"
	exit 0; 
fi

build() {
	if [[ "${MODULE}" == *"frontend"* ]];			then docker build --tag ${FRONTEND}:${TAG} ${ROOT};							fi
	if [[ "${MODULE}" == *"backend"* ]];			then docker build --tag ${BACKEND}:${TAG} ${ROOT}/src/app/backend;				fi
	if [[ "${MODULE}" == *"metrics-scraper"* ]];	then docker build --tag ${SCRAPER}:${TAG} ${ROOT}/src/app/metrics-scraper;	fi
	if [[ "${MODULE}" == *"terminal"* ]];	        then docker build --tag ${TERMINAL}:${TAG} ${ROOT}/src/app/terminal;	fi
}
push() {
	if [[ "${MODULE}" == *"frontend"* ]];			then docker push ${FRONTEND}:${TAG};	fi
	if [[ "${MODULE}" == *"backend"* ]];			then docker push ${BACKEND}:${TAG};	fi
	if [[ "${MODULE}" == *"metrics-scraper"* ]];	then docker push ${SCRAPER}:${TAG};		fi
	if [[ "${MODULE}" == *"terminal"* ]];	        then docker push ${TERMINAL}:${TAG};		fi
}

MODULE="$1"
COMMAND="$2"
TAG="$3"
GROUP="$4"

if [ "${GROUP}" == "" ];	then GROUP="ghcr.io/kore3lab";	fi
PROJECT="kore-board"
ROOT="$(pwd)"
FRONTEND="${GROUP}/${PROJECT}.frontend"
BACKEND="${GROUP}/${PROJECT}.backend"
SCRAPER="${GROUP}/${PROJECT}.metrics-scraper"
TERMINAL="${GROUP}/${PROJECT}.terminal"

echo "Module  : ${MODULE}"
echo "Command : ${COMMAND}"
echo "Tag     : ${TAG}"

if [[ "${COMMAND}" == *"build"* ]];		then build;		fi
if [[ "${COMMAND}" == *"push"* ]];		then push;		fi