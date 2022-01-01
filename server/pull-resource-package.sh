#!/bin/bash

echo "TASK [拉取镜像： $1] ******"
echo "      镜像大小大概 1GB 左右，根据您的网速不同，需要等候的时间不等。"
echo "      如果想要查看下载进度，您可以在运行 kuboard-spray 的服务器上执行以下命令。"
echo -e "     \033[34m docker pull $1 \033[0m"
echo ""

docker pull $1

if [ $? -ne 0 ]; then
  echo "拉取镜像失败"
  exit
fi

version="${1#*:}"

echo -e "      \033[32m[ 拉取镜像成功。]\033[0m "

echo ""
echo "TASK [加载资源包到本地] ********************************************************"
docker run -d --rm --name ${version} $1 sleep 3600

rm -rf "$(pwd)/data/resource/${version}/content/"
docker cp "${version}:/kuboard-spray/resource/content/" "$(pwd)/data/resource/${version}/content/"

echo -e "      \033[32m[ 加载成功。]\033[0m "

echo

echo "TASK [执行清理动作] ************************************************************"

docker stop ${version}


# docker rmi $1

echo -e "      \033[32m[ 清理完成。]\033[0m "

echo ""

echo "PLAY RECAP *********************************************************************"
echo "kuboardspray : ok=3    unreachable=0    failed=0"