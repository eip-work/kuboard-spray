#!/bin/bash

# $1: 文件夹路径
# $2: 文件名
# $3: 下载源 (github.com / gitee.com)
# $4: retries
# $5: download_architecture (amd64 / arm64 / amd64,arm64)
# $6: 是否使用 proxy (true / false)
# $7: proxy 地址

echo "文件夹路径: $1"
echo "文件名    : $2"
echo "下载源    : $3"
echo "重试次数  : $4"
echo "CPU架构   : $5"
echo "启用代理  : $6"
echo "代理地址  : $7"

echo "TASK [创建文件夹] ****"

mkdir -p $1/content

echo ""

echo "TASK [下载资源包] ****"

if [ "$3" != "null" ]; then
    # TODO: 修改地址
    wget -O $1/resource-pack.zip https://$3/Horan-Z/test/archive/refs/tags/$2
else
    echo "skip"
fi

echo ""

echo "TASK [解压资源包] ****"

unzip -o $1/resource-pack.zip -d $1/content

# 找到content下的套壳文件夹
inner_dir=$(find "$1/content" -mindepth 1 -maxdepth 1 -type d)

# 如果找到该目录，移动所有内容（包括隐藏文件和子目录）
if [ -n "$inner_dir" ]; then
    # 使用.*和*匹配所有文件（包括隐藏文件），并抑制通配符不匹配的错误
    mv "$inner_dir"/.* "$inner_dir"/* "$1/content/" 2>/dev/null
    rm -rf "$inner_dir"
fi

cp $1/content/package.yaml $1/package.yaml

echo ""

echo "TASK [执行资源下载脚本] ****"

bash $1/content/download-dependency.sh $4 $5 $6 $7

echo ""