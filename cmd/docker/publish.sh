#!/bin/bash
set -e

cd ../udfs
GOOS=linux GOARCH=amd64 go build -o ../docker/udfs_linux

cd ../docker

 IMAGE_NAME=dockerfordanny/udfs-test
 echo "start build image ${IMAGE_NAME}"
 docker build . -t ${IMAGE_NAME}
 echo "------------------------------------\n"


 echo "start push image ${IMAGE_NAME}"
 docker push ${IMAGE_NAME}



