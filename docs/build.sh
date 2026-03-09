#!/bin/bash
datetime=`date +%Y%m%d-%H%M%S`
echo $datetime
tag=opencmit/pangee-cluster-docs

pnpm install
pnpm docs:build

docker build -t $tag:$datetime .

# docker push $tag:$datetime
# echo pushded $tag:$datetime


# docker run -d -p 8080:80 --name pangee-cluster-docs opencmit/pangee-cluster-docs:20251003-175809