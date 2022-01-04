# 版本记录

## v1.0.0-alpha.2

**优化**

* 重构 package.yaml 结构，不兼容 alpha.1 的 /data 目录
* 优化 ssh 超时时间的设置，避免部分情况下等候时间过长
* 增加 git / nerdctl 工具
  
**问题修复**

* [#3](https://github.com/eip-work/kuboard-spray/issues/3) 启用跳板机设置后再禁用跳板机，安装时提示找不到 bastion
* 检查版本更新时的问题