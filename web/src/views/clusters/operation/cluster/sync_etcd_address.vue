<i18n>
en:
  title: "Update apiserver's --etcd-servers param"
  requiresAllControlNodeOnline: "{node} is offline. Remove it, or bring it back."
  finishAddRemoveActionFirst: "Please finish/cancel {action} action on node {node} first."
  requireAtLeastOneControlPlane: Requires at least one control_plane.

  sync_etcd_address_desc: "Member list of etcd cluster changed, this task is going to update --etcd-servers param in /etc/kubernetes/manifests/kube-apiserver.yaml on all the remaining kube_control_plane nodes, to match the latest etcd member list."
zh:
  title: "更新 apiserver 中 etcd 连接参数"
  requiresAllControlNodeOnline: "安装集群或者添加控制节点时，所有控制节点必须在线。请启动 {node} 或者将其移除。"
  finishAddRemoveActionFirst: "请先完成或取消节点 {node} 的 {action} 操作"
  requireAtLeastOneControlPlane: 至少需要一个控制节点

  sync_etcd_address_desc: "ETCD 集群的成员列表已经发生变化，此操作将更新剩余控制节点上 /etc/kubernetes/manifests/kube-apiserver.yaml 文件中 --etcd-servers 的参数，以匹配最新的 etcd 成员列表。"
</i18n>

<template>
  <el-form-item prop="action" :rules="rules">
    <div class="form_description app_text_mono" style="margin-bottom: 10px;">
      {{ $t('sync_etcd_address_desc') }}
    </div>
    <div>
      <template v-for="(item, key) in cluster.inventory.all.children.target.children.k8s_cluster.children.kube_control_plane.hosts" :key="'kcp' + key">
        <el-tag v-if="pingpong[key] && pingpong[key].status === 'SUCCESS'" style="margin-right: 10px; margin-bottom: 10px;" effect="dark"
          :type="cluster.inventory.all.hosts[key].kuboardspray_node_action === undefined ? 'primary' : 'warning'">
          <span class="app_text_mono">{{key}}</span>
        </el-tag>
      </template>
    </div>
  </el-form-item>
</template>

<script>
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
          // 所有控制节点必须在线
          let controlPlaneCount = 0
          for (let controlPlane in _this.cluster.inventory.all.children.target.children.k8s_cluster.children.kube_control_plane.hosts) {
            if (_this.pingpong[controlPlane] === undefined || _this.pingpong[controlPlane].status !== 'SUCCESS') {
              return callback(_this.$t('requiresAllControlNodeOnline', { node: controlPlane }))
            }
            // 所有控制节点必须都已经完成添加/删除操作
            if (_this.cluster.inventory.all.hosts[controlPlane].kuboardspray_node_action !== undefined) {
              return callback(_this.$t('finishAddRemoveActionFirst', { node: controlPlane, action: _this.cluster.inventory.all.hosts[controlPlane].kuboardspray_node_action }))
            }
            controlPlaneCount ++
          }
          if (controlPlaneCount === 0) {
            return callback(_this.$t('requireAtLeastOneControlPlane'))
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
      return {}
    },
  }
}
</script>

<style scoped lang="scss">
.form_description {
  font-size: 12px;
  // color: var(--el-text-color-placeholder);
  color: #aaa;
  max-width: 700px;
}
</style>
