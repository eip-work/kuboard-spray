#!/bin/bash

datetime=`date "+%Y-%m-%d %H:%M:%S"`
echo $datetime
echo "【构建】 ${1}"

tag=eipwork/kuboard-spray
tag_backup=swr.cn-east-2.myhuaweicloud.com/kuboard/kuboard-spray

echo
echo "【构建 web】"

arch="amd64"
if [ $(uname -m) != "x86_64" ]; then
  arch="arm64"
fi


cd web
yarn install
yarn build

echo \{\"version\":\"${1}-${arch}\",\"buildDate\":\"$datetime\"\} > ./dist/version.json
cd ..

echo
echo "【构建 server】"
rm -f ./server/kuboard-spray
docker run --rm -v ${PWD}:/usr/src/kuboard-spray \
  -v ~/temp/build-temp/pkg:/go/pkg \
  -w /usr/src/kuboard-spray/server golang:1.16.5-buster \
  sh -c "export GOPROXY=https://goproxy.io,direct && go build kuboard-spray.go"

ls -lh ./server/kuboard-spray

echo
echo "【构建 镜像】"

docker build -f Dockerfile -t $tag:$1-${arch} .

echo "【构建 成功】"
echo $tag:$1-${arch}

docker tag $tag:$1-${arch} $tag:latest-${arch}
docker tag $tag:$1-${arch} $tag_backup:$1-${arch}
docker tag $tag:$1-${arch} $tag_backup:latest-${arch}

if [ "$2" == "" ]
then

  docker push $tag:$1-${arch}
  docker push $tag:latest-${arch}
  docker push $tag_backup:$1-${arch}
  docker push $tag_backup:latest-${arch}

else

  echo
  echo "【稍后推送镜像：】"

  echo "#!/bin/bash" > push.sh
  echo "docker push $tag:$1-${arch}" >> push.sh
  echo "docker push $tag:latest-${arch}" >> push.sh
  echo "docker push $tag_backup:$1-${arch}" >> push.sh
  echo "docker push $tag_backup:latest-${arch}" >> push.sh

fi

echo $datetime