
# KuboardSpray

基于 [kubespray](https://github.com/kubernetes-sigs/kubespray) 提供图形化的 K8S 集群离线安装、维护工具。

## 快速安装

找一台不低于1核2G，不少于10G剩余磁盘空间，已经安装好 docker 的服务器，执行如下指令，即可完成 KuboardSpray 的安装：

``` sh
docker run -d \
  --restart=unless-stopped \
  --name=kuboard-spray \
  -p 80:80/tcp \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v /root/kuboard-spray-data:/data \
  eipwork/kuboard-spray:latest-amd64
  # 如果抓不到这个镜像，可以尝试一下这个备用地址：
  # swr.cn-east-2.myhuaweicloud.com/kuboard/kuboard-spray:latest-amd64
```

在浏览器地址栏中输入 `http://这台机器的IP地址`，输入默认密码 `Kuboard123`，即可登录 kuboard-spray 界面，剩下的事情，在界面上操作一下，您就会啦。如果有困难，试试这篇文档 [使用 KuboardSpray 安装 Kubernetes 集群](https://kuboard.cn/install/install-k8s.html)

## 自制资源包

资源包的项目地址在 [kuboard-spray-resource](https://github.com/eip-work/kuboard-spray-resource)，您可以查看一下该项目，尝试自己制作资源包。

## 社区

对此项目感兴趣的同学，请添加本项目的 Star 以后，扫码加入群聊（提供 star 截图才会被拉入群聊哦！）

<p>
  <img src="https://addons.kuboard.cn/downloads/qr_code_kuboard-spray.jpg" style="width: 150px; height: 150px;"/>
</p>