---
---

# 制作资源包

## 资源包的代码也是开源的

[kuboard-spray-resource](https://github.com/eip-work/kuboard-spray-resource)

## 资源包的版本编码规则

资源包的版本号由如下几部分组成：
* `spray-` 前缀 + kubespray的版本号，例如 `spray-v2.18.0-2`，指定了资源包路径中 `/kuboard-spray/resource/content/3rd/kubespray` 目录在 `https://github.com/eip-work/kubespray` 仓库的 Tag；
* `k8s-` 前缀 + k8s 的版本号，例如 `k8s-v1.23.1`，指定了资源包所支持的 Kubernetes 版本；
* 资源包自己的版本号 + 资源包对应的 CPU 架构，例如 `v1.2-amd64`：
  * 对于不同的 kubespray 版本号 + k8s 版本号组合，可能出现相同的资源包版本号；
  * 但是对于相同的 kubespray 版本号 + k8s 版本号组合，资源包版本号应该是唯一的；

以上三个部分由下划线 `_` 分隔，组成一个完整的资源包版本号，例如： `spray-v2.18.0-2_k8s-v1.23.1_v1.2-amd64`。

## 资源包制作方法

* 预计2022年2月提供