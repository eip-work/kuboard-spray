ghcr.io/opencmit/pangee-cluster-devcontainer:latest

镜像使用如下命令构建：

```sh
docker build \
  --build-arg HTTP_PROXY=http://192.168.3.109:7890 \
  --build-arg HTTPS_PROXY=http://192.168.3.109:7890 \
  --build-arg NO_PROXY=localhost,127.0.0.1,.local,192.168.0.0/16 \
  -t ghcr.io/opencmit/pangee-cluster-devcontainer:latest .
```