#!/usr/bin/env bash

which docker
if [[ $? -ne 0 ]] ; then
    echo "Error: docker not found."
    exit 1
fi

CURR_DIR=$(pwd)
VERSION=$1
OUT_DIR=${CURR_DIR}/_output

rm -rf ${OUT_DIR}
mkdir -p ${OUT_DIR}

PREFIX="wavefront-kubernetes-adapter-build"

echo "Building wavefront-kubernetes-adapter docker image..."
docker build -t ${PREFIX}:$VERSION -f ${CURR_DIR}/deploy/Dockerfile-build .

echo "Copying files from docker image..."
CONTAINER_ID=`docker run -d -t ${PREFIX}:${VERSION}`
docker cp ${CONTAINER_ID}:/go/src/github.com/wavefronthq/wavefront-kubernetes-adapter/_output/amd64  ${CURR_DIR}/_output/

echo "Cleaning up build container..."
docker stop ${CONTAINER_ID}
docker container rm ${CONTAINER_ID}
docker rmi $(docker images -f "dangling=true" -q)

echo "Done. Packages available under: ${CURR_DIR}/_output"
