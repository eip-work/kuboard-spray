<i18n>
en:
  confirmDelete: Are you sure to delete this node ?
  noRemoveOffline: cannot remove node when cluster is unreachable
  addNodeFirst: Please complete or cancel the deletion of nodes first.
  cannot_remove_last_etcd: Cannot remove the last etcd node
  cannot_remove_last_master: Cannot remove the last kube_control_plane
zh:
  confirmDelete: 您是否确认要删除此节点？
  noRemoveOffline: 集群不在线，不能删除节点
  addNodeFirst: 请先完成或取消节点添加动作，再执行节点删除动作。
  cannot_remove_last_etcd: 不能删除最后一个 ETCD 节点
  cannot_remove_last_master: 不能删除最后一个控制节点
</i18n>

<template>
  <div class="node_wrapper">
    <div class="status" v-if="name !== 'localhost' && name !== 'bastion'">
      <el-button v-if="cluster.type === 'gap'" type="primary" circle icon="el-icon-document-checked" :loading="pingpong_loading"></el-button>
      <el-button v-else-if="pendingAction === 'remove_node'" type="danger" circle icon="el-icon-delete" :loading="pingpong_loading"></el-button>
      <el-button v-else-if="pendingAction === 'add_node'" type="warning" circle icon="el-icon-plus" :loading="pingpong_loading"></el-button>
      <template v-else-if="k8sNode">
        <el-button v-if="k8sNodeStatus === 'Ready'" type="success" circle :plain="onlineNodes[name] === undefined" icon="el-icon-check" :loading="pingpong_loading"></el-button>
        <el-button v-else type="danger" circle :plain="onlineNodes[name] === undefined" icon="el-icon-moon-night" :loading="pingpong_loading"></el-button>
      </template>
    </div>
    <div class="delete_button" v-if="!hideDeleteButton && editMode !== 'view'">
      <el-button v-if="pendingAction === 'remove_node'" icon="el-icon-check" type="success" circle @click="cancelDelete"></el-button>
      <el-popconfirm v-else icon="el-icon-info" icon-color="red" :title="$t('confirmDelete')" @confirm="deleteNode" placement="right-start">
        <template #reference>
          <el-button icon="el-icon-delete" type="danger" circle @submit.prevent.stop></el-button>
        </template>
      </el-popconfirm>
    </div>
    <div :class="nodeClass">
      <div class="app_text_mono" style="font-weight: bold;">
        <span v-if="name !== 'localhost' && name !== 'bastion'">{{ name }}</span>
        <span v-else>{{ $t('obj.' + name) }}</span>
      </div>
      <div class="app_text_mono" v-if="inventory.all.hosts[name]">
        {{ inventory.all.hosts[name].ansible_host || (name !== 'localhost' && name !== 'bastion' ? '???.???.???.???' : '') }}
      </div>
      <div :class="role + ' role app_text_mono'" v-for="(_, role) in roles" :key="'r' + role">{{ roleName(role) }}</div>
      <slot></slot>
    </div>
  </div>
</template>

<script>

export default {
  props: {
    name: { type: String, required: true },
    cluster: { type: Object, required: true },
    active: { type: Boolean, required: false, default: false },
    hideDeleteButton: { type: Boolean, required: false, default: false },
    pingpong: { type: Object, required: false, default: () => {return {}} },
    pingpong_loading: { type: Boolean, required: false, default: false },
  },
  inject: ['editMode', 'isClusterInstalled', 'isClusterOnline', 'pendingAddNodes', 'onlineNodes'],
  computed: {
    inventory: {
      get () { return this.cluster.inventory },
      set () {}
    },
    pendingAction () {
      if (this.inventory.all.hosts[this.name] && this.inventory.all.hosts[this.name].kuboardspray_node_action) {
        return this.inventory.all.hosts[this.name].kuboardspray_node_action
      }
      return undefined
    },
    nodeClass () {
      let result = 'node'
      if (this.active) {
        result += ' active'
      }
      if (this.onlineNodes[this.name] && this.onlineNodes[this.name].ping === 'pong') {
        result += ' online_node'
      }
      if (this.pendingAction === 'remove_node') {
        result += ' delete_node'
      }
      if (this.name !== 'localhost' && this.name !== 'bastion') {
        if (this.pingpong[this.name] === undefined) {
          result += ' unknown_status'
        } else if (this.pingpong[this.name].ping === 'pong') {
          result += ' online_node'
        } else {
          result += ' offline_node'
        }
      }
      if (this.k8sNode && this.k8sNodeStatus === 'Unknown') {
        result += ' error_node'
      }
      if (this.cluster.type === 'gap') {
        result += ' gap'
      }
      return result
    },
    roles () {
      let result = {}
      for (let role in this.inventory.all.children.target.children.k8s_cluster.children) {
        for (let n in this.inventory.all.children.target.children.k8s_cluster.children[role].hosts) {
          if (n === this.name) {
            result[role] = true
          }
        }
      }
      for (let n in this.inventory.all.children.target.children.etcd.hosts) {
        if (n === this.name) {
          result['etcd'] = true
        }
      }
      return result
    },
    k8sNode () {
      if (this.cluster.state && this.cluster.state.nodes) {
        return this.cluster.state.nodes[this.name]
      }
      return undefined
    },
    k8sNodeStatus () {
      if (this.k8sNode) {
        for (let index in this.k8sNode.status.conditions) {
          let con = this.k8sNode.status.conditions[index]
          if (con.type === 'Ready' && con.status === "True") {
            return 'Ready'
          }
        }
      }
      return 'Unknown'
    }
  },
  mounted () {
  },
  methods: {
    roleName(role) {
      if (role === 'etcd') {
        return this.inventory.all.children.target.children.etcd.hosts[this.name].etcd_member_name || this.$t('node.' + role)
      } else {
        return this.$t('node.' + role)
      }
    },
    deleteNode () {
      if (this.isClusterInstalled && !this.isClusterOnline) {
        this.$message.error(this.$t('noRemoveOffline'))
        return
      }
      if (this.pendingAddNodes.length > 0 && this.onlineNodes[this.name] !== undefined) {
        this.$message.error(this.$t('addNodeFirst'))
        return
      }
      if (this.roles.etcd) {
        let count = 0
        for (let key in this.inventory.all.children.target.children.etcd.hosts) {
          let host = this.inventory.all.children.target.children.etcd.hosts[key]
          if (host.kuboardspray_node_action === undefined) {
            count ++
          }
        }
        if (count === 1) {
          this.$message.error(this.$t('cannot_remove_last_etcd'))
          return
        }
      }
      if (this.roles.kube_control_plane) {
        let count = 0
        for (let key in this.inventory.all.children.target.children.k8s_cluster.children.kube_control_plane.hosts) {
          let host = this.inventory.all.children.target.children.k8s_cluster.children.kube_control_plane.hosts[key]
          if (host.kuboardspray_node_action === undefined) {
            count ++
          }
        }
        if (count === 1) {
          this.$message.error(this.$t('cannot_remove_last_master'))
          return
        }
      }
      if (this.onlineNodes[this.name]) {
        this.inventory.all.hosts[this.name].kuboardspray_node_action = 'remove_node'
        if (this.inventory.all.children.target.children.k8s_cluster.children.kube_control_plane.hosts[this.name]) {
          this.inventory.all.children.target.children.k8s_cluster.children.kube_control_plane.hosts[this.name].kuboardspray_node_action = 'remove_node'
        }
        if (this.inventory.all.children.target.children.etcd.hosts[this.name]) {
          this.inventory.all.children.target.children.etcd.hosts[this.name].kuboardspray_node_action = 'remove_node'
        }
      } else {
        delete this.inventory.all.hosts[this.name]
        delete this.inventory.all.children.target.children.k8s_cluster.children.kube_control_plane.hosts[this.name]
        delete this.inventory.all.children.target.children.k8s_cluster.children.kube_node.hosts[this.name]
        delete this.inventory.all.children.target.children.etcd.hosts[this.name]
        this.$emit('deleted')
      }
    },
    cancelDelete () {
      delete this.inventory.all.hosts[this.name].kuboardspray_node_action
      if (this.inventory.all.children.target.children.k8s_cluster.children.kube_control_plane.hosts[this.name]) {
        delete this.inventory.all.children.target.children.k8s_cluster.children.kube_control_plane.hosts[this.name].kuboardspray_node_action
      }
      if (this.inventory.all.children.target.children.etcd.hosts[this.name]) {
        delete this.inventory.all.children.target.children.etcd.hosts[this.name].kuboardspray_node_action
      }
    }
  }
}
</script>

<style scoped lang="scss">
.node_wrapper{
  .status {
    height: 0px;
    text-align: left;
  }
  .delete_button {
    height: 0px;
    text-align: right;
    display: none;
  }
  .node {
    border: solid 1px var(--el-border-color-light);
    border-radius: 5px;
    width: 120px;
    margin: 5px;
    padding: 15px;
    font-size: 13px;
    cursor: pointer;
    .role {
      font-size: 12px;
      margin-top: 5px;
      color: var(--el-text-color-primary);
      padding: 5px 10px;
      border-radius: 5px;
      width: calc(100% - 20px);
      text-align: center;
    }
    .kube_control_plane {
      background-color: var(--el-color-primary-light-9);
    }
    .kube_node {
      background-color: var(--el-color-success-lighter);
    }
    .etcd {
      background-color: var(--el-color-warning-lighter);
    }
  }
  .node:hover {
    border-color: var(--el-color-primary);
    background-color: var(--el-color-primary-light-9);
    .role {
      color: white;
    }
    .kube_control_plane {
      background-color: var(--el-color-primary);
    }
    .kube_node {
      background-color: var(--el-color-success);
    }
    .etcd {
      background-color: var(--el-color-warning);
    }
  }
  .node.active {
    border-color: var(--el-color-warning);
    background-color: var(--el-color-warning-lighter);
    .role {
      color: white;
    }
    .kube_control_plane {
      background-color: var(--el-color-primary);
    }
    .kube_node {
      background-color: var(--el-color-success);
    }
    .etcd {
      background-color: var(--el-color-warning);
    }
  }
  .node.delete_node {
    text-decoration-line: line-through;
    border-color: var(--el-color-danger) !important;
  }
  .node.delete_node.active {
    background-color: var(--el-color-danger-lighter) !important;
  }
  .node.unknown_status {
    border-block-style: dashed;
    border-inline-style: dashed;
  }
  .node.online_node {
    border-color: var(--el-color-success);
    border-block-style: solid;
    border-inline-style: solid;
  }
  .node.online_node.active {
    background-color: var(--el-color-success-lighter);
  }
  .node.offline_node {
    border-block-style: dashed;
    border-inline-style: dashed;
    background-color: var(--el-color-info-lighter);
    border-color: var(--el-color-info);
    .role {
      background-color: var(--el-color-info-light);
    }
  }
  .node.offline_node.active {
    .role {
      opacity: 0.5;
      background-color: var(--el-color-info);
    }
    border-color: var(--el-color-warning);
  }
  .node.gap {
    background-color: var(--el-color-primary-light-9) !important;
    border-color: var(--el-color-primary);
  }
  .node.gap.active {
    background-color: var(--el-color-primary-light-6) !important;
  }
  .node.error_node {
    .role {
      opacity: 0.5;
      background-color: var(--el-color-info);
    }
  }
  .node.error_node.active {
    border-color: var(--el-color-danger) !important;
    background-color: var(--el-color-danger-lighter) !important;
  }
}

.node_wrapper:hover .delete_button {
  display: block;
}

</style>
