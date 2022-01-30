---
title: 高可用
description: Kuboard Spray 安装的 Kubernetes 集群如何实现高可用
---

# 高可用

> 本文主要内容参考自： [kubespray 实现高可用的方式](https://kubespray.io/#/docs/ha-mode)

Kuboard-Spray 基于 kubespray 实现 Kubernetes 集群安装，在 Kubernetes 集群中，如下组件需要提供高可用的访问端口：
* etcd 集群
* kube-apiserver 服务

后面的章节分别阐述基于 Kuboard-Spray 安装 k8s 集群时，这两个组件实现高可用的方式。

## ETCD

### ETCD 集群的高可用

[ETCD集群的高可用](https://etcd.io/docs/v3.5/faq/#what-is-failure-tolerance) 由 etcd 集群本身保证。只要集群中的成员数量不低于规定的最小数量，集群就可以正常运作。如果因为网络故障导致集群成员的失联，etcd 可以在网络恢复以后自动恢复到正常的状态，并且通过 Raft 算法保证集群的一致性；如果因为集群节点的机器临时故障（例如断电）使某个 etcd 集群成员失效，在该集群成员重新启动后，该成员将读取日志中故障发生之前的交易，并重新加入集群；如果集群中的某个节点硬件发生永久性的故障，则可以直接将其从集群中移除，并加入新的节点以替换该故障节点。

对于 ETCD 集群，建议在集群中提供奇数个节点，下表显示了不同的节点数量时 ETCD 集群可以容忍的错误节点数量：

|  集群节点数   | Majority  | 最大容错数 |
|  ----  | ----  |----|
| 1  | 1 |  0 |
| 2  | 2 |  0 |
| 3  | 2 |  1 |
| 4  | 2 |  1 |
| 5  | 3 |  2 |
| 6  | 3 |  2 |
| 7  | 4 |  3 |
| 8  | 4 |  3 |
| 9  | 5 |  4 |

从表中可以看出，偶数节点集群（例如8节点集群）相比奇数节点集群（例如7节点集群）并不能带来更大的容错能力。相反，在出现网络分裂（network partition）时，奇数节点数的集群可以始终确保分割后的子集中，有一个子集中的节点数量大于另一个子集中的节点数量，并在网络恢复时，可以无争议地以节点数量大的集群子集中的数据作为有效数据。

### KuboardSpray 部署的 ETCD

在使用 KuboardSpray 部署 K8S 时，如果有多个 ETCD 节点，则这些节点就已经按照 ETCD 集群的方式做好配置。并且所有 ETCD 节点的访问端口都被配置到了 ETCD 客户端（例如 kube-api-server）中的 etcd 连接参数中。

## kube-apiserver

API-Server 是一个无状态服务，可以通过负载均衡，将请求轮发到集群中的任意一个 kube-apiserver 节点。Kuboard Spray （基于 kubespray 实现）在每个非控制节点上配置了一个 nginx 反向代理，作为负载均衡器使用。这种做法被称作 localhost 负载均衡。这种做法相比于使用一个独立的负载均衡器来给 kube-apiserver 分发请求，主要的优势在于简单和便捷，尤其当用户不便于提供外部的负载均衡器或者虚拟IP时非常有效。使用 KuboardSpray 安装 K8S 集群时，所有的非控制节点上的 kubelet / kube-proxy 组件都将通过此 localhost 负载均衡访问 kube-apiserver。

如下图所示：
![本地负载均衡](./ha-mode.assets/loadbalancer_localhost.png)

master 节点（控制节点）上的 kubelet / kube-proxy / kube-scheduler / kube-controller-manager 等组件都将直接访问 kube-apiserver 的 `localhost:8080` 端口；而非控制节点上的 kubelet / kube-proxy 组件，则将访问该节点 nginx 反向代理的 `localhost:443` 端口，并由 nginx 反向代理进一步将请求轮发到某一个 master 节点（控制节点）上的 kube-apiserver 的 `:6443` 端口。

::: tip 外置的负载均衡

* 如果您直接使用 kubespray 部署 K8S 集群，也可以为 kube-apiserver 配置一个外置的负载均衡器。但是KuboardSpray 暂时未开放此配置。

:::
