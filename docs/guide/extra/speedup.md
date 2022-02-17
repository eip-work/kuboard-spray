---
sidebarDepth: 0
---

# 加速执行

KuboardSpray 已经尽可能地优化了 ansible 的执行速度，想要执行得快一些，有如下建议：

* KuboardSpray 与集群的目标机器部署在同样的内网环境
* [关闭 sshd 的 UseDNS](#关闭-ssh-的-usedns)
* 如果可能，不使用跳板机 / 堡垒机
* 如果必须使用跳板机 / 堡垒机，则，为目标机器配置 ssh private key 登录，而不是使用密码

## 关闭 ssh 的 UseDNS

对于所有的目标服务器，执行以下修改：

* 编辑 `/etc/ssh/sshd_config`
  ``` sh
  vi /etc/ssh/sshd_config
  ```
* 添加或修改
  ```ini
  UseDNS no
  ```
* 重启 `sshd` 服务
  ```sh
  systemctl restart sshd
  ```
## 使用最近的镜像仓库