#!/bin/bash

datetime=`date "+%Y-%m-%d %H:%M:%S"`

echo $datetime

echo "【构建】 ${1}"

tag=ghcr.io/opencmit/pangee-cluster

echo
echo "【构建 web】"

cd web
pnpm install
pnpm build

echo \{\"version\":\"${1}\",\"buildDate\":\"$datetime\"\} > ./dist/version.json
cd ..

echo
echo "【构建 镜像】"

docker buildx build --platform linux/amd64,linux/arm64 -t $tag:$1 -t $tag:latest --load .

echo "【构建 成功】 请执行下面的指令推送到镜像仓库"

echo "docker push $tag:$1"
echo "docker push $tag:latest"

echo $datetime