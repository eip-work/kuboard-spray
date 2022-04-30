---
---

# 重置密码

本文描述了如何重置 KuboardSpray 中 `admin` 用户的密码。请按照如下步骤操作：

* 执行如下命令，进入 kuboard-spray 容器的命令行界面：
  ```sh
  docker exec -it kuboard-spray /bin/bash
  ```

* 在 kuboard-spray 容器中执行如下命令：
  ```sh
  ./kuboard-spray-admin reset-password
  ```

  输出结果如下所示：
  ```log
  设置日志级别为 info
  2022/04/30 22:41:39.801 |  info | try to reset-password
  2022/04/30 22:41:39.803 |  info | 已将 admin 的密码重置为 Kuboard123
  ```

  至此，您已经成功重置了 `admin` 用户的密码，重置后密码为 `Kuboard123`