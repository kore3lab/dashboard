#!/bin/sh
# -----------------------------------------------------------------
# usage
if [ "$#" -lt 3 ]; then 
	echo "docker-build.sh <module name - backend/frontend/dashboard/all> <command> <tag>"
	echo "    ./docker-build.sh backend build v0.1.0"
	exit 0; 
fi

build() {
	if [[ "${MODULE}" == *"frontend"* ]];			then docker build --tag ${FRONTEND}:${TAG} ${ROOT};							fi
	if [[ "${MODULE}" == *"dashboard"* ]];			then docker build --tag ${DASHBOARD}:${TAG} ${ROOT}/dashboard;				fi
	if [[ "${MODULE}" == *"metrics-scraper"* ]];	then docker build --tag ${SCRAPER}:${TAG} ${ROOT}/src/app/metrics-scraper;	fi
}
push() {
	if [[ "${MODULE}" == *"frontend"* ]];			then docker push ${FRONTEND}:${TAG};	fi
	if [[ "${MODULE}" == *"dashboard"* ]];			then docker push ${DASHBOARD}:${TAG};	fi
	if [[ "${MODULE}" == *"metrics-scraper"* ]];	then docker push ${SCRAPER}:${TAG};		fi
}

MODULE="$1"
COMMAND="$2"
TAG="$3"
GROUP="$4"

if [ "${GROUP}" == "" ];	then GROUP="ghcr.io/acornsoftlab";	fi
PROJECT="acornsoft-dashboard"
ROOT="$(pwd)"
FRONTEND="${GROUP}/${PROJECT}.frontend"
BACKEND="${GROUP}/${PROJECT}.backend"
SCRAPER="${GROUP}/${PROJECT}.metrics-scraper"


echo "Module  : ${MODULE}"
echo "Command : ${COMMAND}"
echo "Tag     : ${TAG}"

if [[ "${COMMAND}" == *"build"* ]];		then build;		fi
if [[ "${COMMAND}" == *"push"* ]];		then push;		fi