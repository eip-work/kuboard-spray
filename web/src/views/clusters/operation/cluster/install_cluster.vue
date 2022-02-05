<i18n>
en:
  install_cluster: Install / Setup K8S Cluster
  requiresAtLeastOneOnlineNode: All nodes are offline.
  requiresAllControlNodeOnline: "{node} is offline. Remove it, or bring it back."
  requireAtLeastOneControlPlane: Requires at least one control_plane.
  requiresOddEtcdCount: Etcd count should be an odd number.
  requiresKubeNodeCount: Requires at lease one kube_node.
  requiresAllEtcdNodeOnline: "All Etcd nodes must be online, currently, we cannot reach node {node}"
  nodesToIncludeDesc: Include nodes
  re_install_cluster: "This is going to repeat the installation steps on the following nodes, generally speaking, you don't need to do _this."

  skip_downloads: Skip Downloads
zh:
  install_cluster: 安装 / 设置集群
  requiresAtLeastOneOnlineNode: 至少需要一个在线的节点
  requiresAllControlNodeOnline: "安装集群或者添加控制节点时，所有控制节点必须在线。请启动 {node} 或者将其移除。"
  requireAtLeastOneControlPlane: 至少需要一个控制节点
  requiresOddEtcdCount: ETCD 节点的数量必须为奇数
  requiresKubeNodeCount: 至少要有一个在线的工作节点
  requiresAllEtcdNodeOnline: "所有 ETCD 节点必须在线，当前 {node} 不在线"
  nodesToIncludeDesc: 包含节点
  re_install_cluster: 此操作将在如下节点上重复执行一次集群安装的动作，通常您不需要这么做。

  skip_downloads: 跳过下载
</i18n>

<template>
  <el-form-item prop="action" :rules="rules">
    <el-alert v-if="isClusterInstalled" type="warning" :title="$t('install_cluster')"
      class="app_margin_bottom" :closable="false">
      {{$t('re_install_cluster')}}
    </el-alert>
    <el-form-item :label="$t('nodesToIncludeDesc')">
      <template v-for="(item, key) in hosts" :key="'h' + key">
        <el-tag v-if="pingpong[key] && pingpong[key].ping === 'pong'" style="margin-bottom: 10px; margin-right: 10px;" effect="dark">
          <span class="app_text_mono" style="font-size: 14px;">{{key}}</span>
        </el-tag>
      </template>
    </el-form-item>
    <el-form-item :label="$t('skip_downloads')" v-if="isClusterInstalled">
      <el-switch v-model="formRef.install_cluster.skip_downloads"></el-switch>
    </el-form-item>
  </el-form-item>
</template>

<script>
import { trimMark } from './utils'

export default {
  props: {
    cluster: { type: Object, required: true },
    name: { type: String, required: true },
    form: { type: Object, required: true },
    pingpong: { type: Object, required: false, default: () => { return {} } },
    pingpong_loading: { type: Boolean, required: false, default: false },
    offlineNodes: { type: Array, required: false, default: () => { return [] } },
    hosts: { type: Object, required: true },
  },
  data() {
    let _this = this
    return {
      rules: [{
        validator (rule, value, callback) {
          // 至少有一个被添加的节点在线
          let temp = ''
          for (let node in _this.cluster.inventory.all.hosts) {
            if (_this.pingpong[node] && _this.pingpong[node].ping === 'pong') {
              temp += node + ','
            }
          }
          if (temp === '') {
            return callback(_this.$t('requiresAtLeastOneOnlineNode'))
          }

          // 至少有一个工作节点
          let countKubeNodeCount = 0
          for (let node in _this.cluster.inventory.all.children.target.children.k8s_cluster.children.kube_node.hosts) {
            if (_this.pingpong[node] && _this.pingpong[node].ping === 'pong') {
              countKubeNodeCount ++
            }
          }
          if (countKubeNodeCount === 0) {
            return callback(_this.$t('requiresKubeNodeCount'))
          }


          // 所有的控制节点必须在线
          let controlPlaneCount = 0
          for (let controlPlane in _this.cluster.inventory.all.children.target.children.k8s_cluster.children.kube_control_plane.hosts) {
            if (_this.pingpong[controlPlane] === undefined || _this.pingpong[controlPlane].ping !== 'pong') {
              return callback(_this.$t('requiresAllControlNodeOnline', { node: controlPlane }))
            }
            controlPlaneCount ++
          }
          if (controlPlaneCount === 0) {
            return callback(_this.$t('requireAtLeastOneControlPlane'))
          }
          // 所有的 etcd 节点必须在线
          let etcdCount = 0
          for (let etcd in _this.cluster.inventory.all.children.target.children.etcd.hosts) {
            if (_this.pingpong[etcd] === undefined || _this.pingpong[etcd].ping !== 'pong') {
              return callback(_this.$t('requiresAllEtcdNodeOnline, { node: etcd }'))
            }
            etcdCount ++
          }
          if (etcdCount % 2 == 0) {
            return callback(_this.$t('requiresOddEtcdCount'))
          }
          return callback()
        },
        trigger: 'change',
      }]
    }
  },
  inject: ['isClusterInstalled', 'isClusterOnline', 'pendingRemoveNodes', 'pendingAddNodes'],
  computed: {
    formRef: {
      get () { return this.form },
      set () {}
    }
  },
  components: { },
  beforeMount () {
    this.formRef.install_cluster = { skip_downloads: this.isClusterInstalled }
  },
  methods: {
    populateRequest() {
      let temp = ''
      for (let node in this.hosts) {
        if (this.offlineNodes.indexOf(node) < 0) {
          temp += node + ','
        }
      }
      return { nodes: trimMark(temp), skip_downloads: this.form.install_cluster.skip_downloads }
    },
  }
}
</script>

<style scoped lang="css">
.form_description {
  font-size: 12px;
  color: #aaa;
  max-width: 700px;
}
</style>
