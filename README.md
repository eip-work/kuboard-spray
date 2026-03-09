# PangeeCluster

基于 [kubespray](https://github.com/kubernetes-sigs/kubespray) 提供图形化的 K8S 集群离线安装、维护工具。

## 快速安装

PangeeCluster 支持以下两种安装方式，请选择其中一种进行安装。

### docker 安装
**使用 docker 安装时，该机器不可用于集群节点**
找一台不低于 1 核 2G，不少于 20G 剩余磁盘空间，已经安装好 docker 的服务器，执行如下指令，即可完成 PangeeCluster 的安装：

```sh
docker run -d \
  --restart=unless-stopped \
  --name=pangee-cluster \
  -p 8080:8080/tcp \
  -e PANGEE_CLUSTER_PORT=8080 \
  -e TZ=Asia/Shanghai \
  -v ~/pangee-cluster-data:/data \
  ghcr.io/opencmit/pangee-cluster:latest
```

在浏览器地址栏中输入 `http://这台机器的IP地址:端口号`，输入用户名 `admin`，默认密码 `PangeeCluster123`，即可登录 pangee-cluster 界面，剩下的事情，在界面上操作一下，您就会啦。如果有困难，试试这篇文档 <a href="./guide/install-k8s.html" target="_blank">使用 PangeeCluster 安装 Kubernetes 集群</a>

### 可执行文件
下载 GitHub release 页面的 pangee-cluster-bin 文件，执行如下指令，即可启动 PangeeCluster：

```sh
/path/to/pangee-cluster-bin \
  -e TZ=Asia/Shanghai \
  -e PANGEE_CLUSTER_PORT=8080 \
  -v /path/to/data:/data
```

在浏览器地址栏中输入 `http://这台机器的IP地址:端口号`，输入用户名 `admin`，默认密码 `PangeeCluster123`，即可登录 pangee-cluster 界面，剩下的事情，在界面上操作一下，您就会啦。如果有困难，试试这篇文档 <a href="./guide/install-k8s.html" target="_blank">使用 PangeeCluster 安装 Kubernetes 集群</a>

**常见问题**

- 导入资源包时，可能会碰到 `no such file or directory` 或者 `permission denied` 之类的错误提示，通常是因为您开启了 SELinux，导致 pangee-cluster 不能读取映射到容器 `/data` 的路径
- pangee-cluster 所在机器不能当做 K8S 集群的一个节点，因为安装过程中会重启集群节点的容器引擎，这会导致 pangee-cluster 被重启掉。

## 配置开发测试环境

[配置开发测试环境](./docs/setup-dev/dev.md)

