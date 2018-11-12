#!/bin/bash
set -e

cd ../ipfs
GOOS=linux GOARCH=amd64 go build -o ../docker/ipfs_linux

cd ../docker

# IMAGE_NAME=dockerfordanny/udfs-test:latest
# echo "start build image ${IMAGE_NAME}"
# docker build . -t ${IMAGE_NAME}
# echo "------------------------------------\n"


# echo "start push image ${IMAGE_NAME}"
# docker push ${IMAGE_NAME}



