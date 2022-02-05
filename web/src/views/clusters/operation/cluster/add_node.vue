<i18n>
en:
  add_nodes_desc: You are going to add the following nodes into the cluster.
  requiresAtLeastOneOnlineNode: All nodes pending to add are offline.
  requiresAllControlNodeOnline: "{node} is offline. Remove it, or bring it back."
  requireAtLeastOneControlPlane: Requires at least one control_plane.
  requiresOddEtcdCount: Etcd count should be an odd number.
  requiresKubeNodeCount: Requires at lease one kube_node.

  requiresAllEtcdNodeOnline: "All Etcd nodes must be online, currently, we cannot reach node {node}"
zh:
  add_nodes_desc: 将要添加以下节点
  requiresAtLeastOneOnlineNode: 至少需要一个在线的待添加节点
  requiresAllControlNodeOnline: "安装集群或者添加控制节点时，所有控制节点必须在线。请启动 {node} 或者将其移除。"
  requireAtLeastOneControlPlane: 至少需要一个控制节点
  requiresOddEtcdCount: ETCD 节点的数量必须为奇数
  requiresKubeNodeCount: 至少要有一个在线的工作节点

  requiresAllEtcdNodeOnline: "所有 ETCD 节点必须在线，当前 {node} 不在线"
</i18n>

<template>
  <el-form-item prop="action" :rules="rules">
    <div class="form_description">
      {{ $t('add_nodes_desc') }}
    </div>
    <template v-for="(item, key) in pendingAddNodes" :key="'h' + key">
      <el-tag v-if="pingpong[item.name] && pingpong[item.name].ping === 'pong'" style="margin-top: 10px; margin-right: 10px;" :disabled="pingpong[item.name].ping !== 'pong'" :label="item.name">
        <span class="app_text_mono" style="font-size: 14px;">{{item.name}}</span>
      </el-tag>
    </template>
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
              if (_this.cluster.inventory.all.hosts[node].kuboardspray_node_action === 'add_node') {
                temp += node + ','
              }
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
  },
  components: { },
  mounted () {
  },
  methods: {
    populateRequest() {
      let _this = this
      let temp = ''
      for (let node of _this.pendingAddNodes) {
        if (_this.pingpong[node.name] && _this.pingpong[node.name].ping === 'pong') {
          temp += node.name + ','
        }
      }
      return { nodes_to_add: trimMark(temp) }
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
