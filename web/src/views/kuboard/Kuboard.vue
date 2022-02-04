<i18n>
en:
  hi: kuboard
zh:
  hi: kuboard
</i18n>

<template>
  <div>
    <div class="app_block_title">Kuboard</div>
    <div class="description">
      使用 Kuboard 可以管理 K8S 集群，Kuboard 与 Kuboard-Spray 的关系：
      <li>Kuboard-Spray 通过 ansible / ssh 命令直接操作集群节点的操作系统，执行 Kubernetes 集群的安装、添加节点、删除节点、备份、漏洞扫描等任务，使用 Kuboard-Spray 的用户通常是 Kubernetes 集群维护人员。</li>
      <li>Kuboard 通过 kubernetes 的 apiserver 执行集群的日常管理操作，例如：名称空间创建、Deployment 创建修改、Service 的创建修改等，Kuboard 支持细致的权限控制，可以使不同的团队共享使用同一个集群。</li>
      <p>
        通过此链接，详细了解 Kuboard。
        <el-link target="blank" href="https://kuboard.cn">
          <el-icon style="vertical-align: text-bottom;">
            <el-icon-link></el-icon-link>
          </el-icon>
          https://kuboard.cn
        </el-link> 
      </p>

    </div>
  </div>
</template>

<script>
import mixin from '../../mixins/mixin.js'

export default {
  mixins: [mixin],
  props: {
  },
  percentage () {
    return 100
  },
  breadcrumb () {
    return [
      { label: 'Kuboard 集成' }
    ]
  },
  computed: {
  },
  components: { },
  mounted () {
  },
  methods: {

  }
}
</script>

<style scoped lang="css">
.description {
  font-size: 13px;
  padding: 20px;
  line-height: 28px;
}
</style>
