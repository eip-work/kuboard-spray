---
sidebarDepth: 0
---

# PangeeCluster

中国移动开源的图形化的 K8S 集群离线安装、维护工具。

## 快速安装

### docker

找一台不低于 1 核 2G，不少于 10G 剩余磁盘空间，已经安装好 docker 的服务器，执行如下指令，即可完成 PangeeCluster 的安装：

```sh
docker run -d \
  --privileged \
  --restart=unless-stopped \
  --name=pangee-cluster \
  -p 9080:9080/tcp \
  -e TZ=Asia/Shanghai \
  -v ~/pangee-cluster-data:/data \
  opencmit/pangee-cluster:latest-amd64
  # 如果您是 arm64 环境，请将标签里的 amd64 修改为 arm64，例如 opencmit/pangee-cluster:latest-arm64
```

在浏览器地址栏中输入 `http://这台机器的IP地址:9080`，输入用户名 `admin`，默认密码 `PangeeCluster123`，即可登录 pangee-cluster 界面，剩下的事情，在界面上操作一下，您就会啦。如果有困难，试试这篇文档 [使用 PangeeCluster 安装 Kubernetes 集群](./guide/install-k8s.md)

::: tip 常见问题

- 如果您安装 Kubernetes 集群时，打算将 pangee-cluster 所在服务器也当做 kubernetes 集群的一个节点，请使用二进制方式运行 `pangee-cluster-bin`，否则可能会因为容器引擎的配置导致冲突（运行 pangee-cluster 时的容器引擎，与 pangee-cluster 安装集群时在节点上安装的容器引擎）
- PangeeCluster 的信息保存在容器的 `/data` 路径，请将其映射到一个您认为安全的地方，上面的命令中，将其映射到了 `~/pangee-cluster-data` 路径；
- 只要此路径的内容不受损坏，重启、升级、重新安装 Pangee-Cluster，或者将数据及 Pangee-Cluster 迁移到另外一台机器上，您都可以找回到原来的信息；

:::

### 二进制

```sh
# 创建数据目录
mkdir ~/pangee-cluster-bin-data

# 运行 pangee-cluster-bin
sudo ./pangee-cluster-bin \
  -e PANGEE_CLUSTER_PORT=9080 \
  -e TZ=Asia/Shanghai \
  -v ~/pangee-cluster-bin-data:/data &

# 以此方式运行时，建议执行下面的指令，并保持终端不退出，否则 pangee-cluster 会随终端关闭而退出
top

# 停止 pangee-cluster-bin
sudo pkill pangee-cluster-bin
```

在浏览器地址栏中输入 `http://这台机器的IP地址:9080`，输入用户名 `admin`，默认密码 `PangeeCluster123`，即可登录 pangee-cluster 界面，剩下的事情，在界面上操作一下，您就会啦。如果有困难，试试这篇文档 [使用 PangeeCluster 安装 Kubernetes 集群](./guide/install-k8s.md)

::: tip 常见问题

- pangee-cluster-bin 是使用 [dockerc](https://github.com/NilsIrl/dockerc) 将容器镜像编译成了二进制程序，不需要在运行环境中安装 docker。如果您安装 Kubernetes 集群时，打算将 pangee-cluster 所在服务器也当做 kubernetes 集群的一个节点，请使用二进制方式运行 `pangee-cluster-bin`，否则可能会因为容器引擎的配置导致冲突（运行 pangee-cluster 时的容器引擎，与 pangee-cluster 安装集群时在节点上安装的容器引擎）
- PangeeCluster 的信息保存在容器的 `/data` 路径，请将其映射到一个您认为安全的地方，上面的命令中，将其映射到了 `~/pangee-cluster-data` 路径；
- 只要此路径的内容不受损坏，重启、升级、重新安装 Pangee-Cluster，或者将数据及 Pangee-Cluster 迁移到另外一台机器上，您都可以找回到原来的信息；

:::

## 自制资源包

Pangee-Cluster 将定期提供最新版本的资源包，可以在 pangee-cluster 资源包管理界面中查到，如果您是离线环境，也可以在 [https://github.com/opencmit/pangee-cluster-resource-package](https://github.com/opencmit/pangee-cluster-resource-package/)找到最新的资源包。您也可以参考该资源包项目自己定制资源包。