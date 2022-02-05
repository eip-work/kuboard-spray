#!/bin/bash


echo "TASK [检查镜像是否已下载： $1] ****"
image="${1%:*}"
version="${1#*:}"
echo "      image  : ${image}"
echo "      tag    : ${version}"

tag=$(docker images --format "{{.Repository}}:{{.Tag}}" | grep "$1")

if [ "$tag" == "$1" ]; then
  echo -e "      \033[32m[ 镜像已存在，无需重复下载。]\033[0m "
else
  echo ""
  echo "TASK [拉取镜像： $1] ******"
  echo "      镜像大小大概 1GB 左右，根据您的网速不同，需要等候的时间不等。"
  echo "      如果想要查看下载进度，您可以在运行 kuboard-spray 的服务器上执行以下命令。"
  echo -e "     \033[34m docker pull $1 \033[0m"
  echo ""

  docker pull $1

  if [ $? -ne 0 ]; then
    echo -e "      \033[31m\033[01m\033[05m[ 拉取镜像失败。]\033[0m "
    echo -e "      如果您需要离线下载，请参考 https://kuboard.cn/support/kuboard-spray"
    exit
  fi
  echo -e "      \033[32m[ 拉取镜像成功。]\033[0m "
fi

echo ""
echo "TASK [加载资源包到本地] ********************************************************"
docker run -d --rm --name ${version} $1 sleep 3600

dataDir=$(pwd)
dataDir=${dataDir%/*}/data

rm -rf "${dataDir}/resource/${version}/content/"
docker cp "${version}:/kuboard-spray/resource/content/" "${dataDir}/resource/${version}/content/"

if [ $? -ne 0 ]; then
  echo -e "      \033[31m\033[01m\033[05m[ 加载资源包失败。]\033[0m "
  exit
fi

echo -e "      \033[32m[ 加载成功。]\033[0m "

echo ""

echo "TASK [执行清理动作] ************************************************************"

docker stop ${version}


docker rmi $1

echo -e "      \033[32m[ 清理完成。]\033[0m "

echo ""

echo "PLAY RECAP *********************************************************************"
echo "kuboardspray : ok=3    unreachable=0    failed=0"