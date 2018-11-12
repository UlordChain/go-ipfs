#!/bin/bash

UDFS_PATH=$1
if [ "$UDFS_PATH" == "" ]; then
	echo "必须指明UDFS工作目录"
	exit 1
fi

IMAGE_NAME=dockerfordanny/udfs-test:latest

# rm old container images
docker rmi -f ${IMAGE_NAME} > /dev/null
docker stop udfs > /dev/null
docker rm -f udfs > /dev/null

# run container with name udfs
docker run -p 4001:4001 -v ${UDFS_PATH}:/root/.ipfs --privileged=true --name udfs -d ${IMAGE_NAME}