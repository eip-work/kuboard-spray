---
sidebarDepth: 0
---

# 加速执行

KuboardSpray 已经尽可能地优化了 ansible 的执行速度，想要执行得快一些，有如下建议：

* KuboardSpray 与集群的目标机器部署在同样的内网环境
* 如果可能，不使用跳板机 / 堡垒机
* 如果必须使用跳板机 / 堡垒机，则，为目标机器配置 ssh private key 登录，而不是使用密码



