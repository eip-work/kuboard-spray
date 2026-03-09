<template>
  <div class="pending_node app_text_mono">
    {{ node.ansible_host }}
    <el-tooltip >
      <template #content>
        <div class="app_text_mono">
          {{ nodeName }}
        </div>
        <div class="app_text_mono">
          {{ roles }}
        </div>
      </template>
      <div class="pending_node_roles">
        <div v-if="roles['kube_control_plane']" :class="`kube_control_plane${roles['kube_control_plane'] ? '_enabled' : ''}`">M</div>
        <div v-if="roles['kube_node']" :class="`kube_node${roles['kube_node'] ? '_enabled' : ''}`">W</div>
        <div v-if="roles['etcd']" :class="`etcd${roles['etcd'] ? '_enabled' : ''}`">E</div>
      </div>
    </el-tooltip>
  </div>
</template>

<script>
export default {
  props: {
    nodeName: { type: String, required: true },
    node: { type: Object, required: true },
    inventory: { type: Object, required: true }
  },
  data() {
    return {

    }
  },
  computed: {
    roles () {
      let result = {}
      for (let role in this.inventory.all.children.target.children.k8s_cluster.children) {
        for (let n in this.inventory.all.children.target.children.k8s_cluster.children[role].hosts) {
          if (n === this.nodeName) {
            result[role] = true
          }
        }
      }
      for (let n in this.inventory.all.children.target.children.etcd.hosts) {
        if (n === this.nodeName) {
          result['etcd'] = true
        }
      }
      return result
    }
  }
}
</script>

<style lang="scss" scoped>
.pending_node {
  border: solid 1px var(--el-border-color-light);
  border-radius: 5px;
  padding: 0 0 0 10px;
  background-color: var(--el-color-white);
  display: flex;
  gap: 5px;
  .pending_node_roles {
    display: flex;
    div {
      padding: 0 3px;
    }
    .kube_control_plane {
      background-color: var(--el-color-primary-light-9);
      color: var(--el-text-color-placeholder);
    }
    .kube_control_plane_enabled {
      background-color: var(--el-color-primary);
      font-weight: bolder;
      color: var(--el-color-white);
    }
    .kube_node {
      background-color: var(--el-color-success-light-9);
      color: var(--el-text-color-placeholder);
    }
    .kube_node_enabled {
      background-color: var(--el-color-success);
      font-weight: bolder;
      color: var(--el-color-white);
    }
    .etcd {
      background-color: var(--el-color-warning-light-9);
      color: var(--el-text-color-placeholder);
    }
    .etcd_enabled {
      background-color: var(--el-color-warning);
      font-weight: bolder;
      color: var(--el-color-white);
    }
  }
}
</style>