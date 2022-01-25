<i18n>
en:
  roles: Node Role
  label: Node Info
  description: Node info retrived from kubernetes cluster.
zh:
  roles: 节点角色
  label: 节点信息
  description: 从集群在线获取到的节点信息

</i18n>

<template>
  <ConfigSection v-model:enabled="enabled" :label="$t('label')" :description="$t('description', {nodeName: nodeName})" disabled anti-freeze>
    <el-form label-width="160px" class="app_form_mini">
      <el-form-item :label="$t('roles')" label-width="100px">
        <div class="info" style="display: flex; flex-wrap: wrap;">
          <template v-if="node.k8s_node">
            <NodeRoleTag v-if="node.k8s_node.metadata.labels['node-role.kubernetes.io/control-plane'] !== undefined" role="kube_control_plane" style="margin-right: 10px;" enabled read-only></NodeRoleTag>
            <NodeRoleTag v-if="isKubeNode" role="kube_node" enabled read-only></NodeRoleTag>
          </template>
          <template v-if="node.etcd_member">
            <NodeRoleTag role="etcd" enabled read-only>{{node.etcd_member.name}}</NodeRoleTag>
          </template>
        </div>
      </el-form-item>
      <template v-if="node.k8s_node">

        <el-form-item label="base info" label-width="100px">
          <div class="info">
            <FieldString :holder="node.k8s_node.spec" fieldName="podCIDR"></FieldString>
            <FieldString v-for="(addr, index) in node.k8s_node.status.addresses" :key="'a' + index" :holder="addr" :label="addr.type" fieldName="address"></FieldString>
            <FieldString :holder="node.k8s_node.status.nodeInfo" fieldName="containerRuntimeVersion"></FieldString>
            <FieldString :holder="node.k8s_node.status.nodeInfo" fieldName="kubeletVersion"></FieldString>
            <FieldString :holder="node.k8s_node.status.nodeInfo" fieldName="osImage"></FieldString>
          </div>
        </el-form-item>
        <el-form-item label="allocatable" label-width="100px">
          <div class="info">
            <FieldString :holder="node.k8s_node.status.allocatable" fieldName="cpu"></FieldString>
            <FieldString :holder="node.k8s_node.status.allocatable" fieldName="memory"></FieldString>
            <FieldString :holder="node.k8s_node.status.allocatable" fieldName="ephemeral-storage"></FieldString>
            <FieldString :holder="node.k8s_node.status.allocatable" fieldName="pods"></FieldString>
          </div>
        </el-form-item>
        <el-form-item label="conditions" label-width="100px">
          <div class="info">
            <StateNodeCondition v-for="(condition, index) in node.k8s_node.status.conditions" :key="'condition' + index" :condition="condition"></StateNodeCondition>
          </div>
        </el-form-item>
      </template>
    </el-form>
  </ConfigSection>
</template>

<script>
import StateNodeCondition from './StateNodeCondition.vue'
import NodeRoleTag from '../common/NodeRoleTag.vue'

export default {
  props: {
    cluster: { type: Object, required: true },
    nodeName: { type: String, required: true },
  },
  data() {
    return {

    }
  },
  inject: ['onlineNodes'],
  computed: {
    enabled: {
      get () {return true},
      set () {}
    },
    node: {
      get () {return this.onlineNodes[this.nodeName] || {}},
      set () {}
    },
    isKubeNode () {
      let node = this.onlineNodes[this.nodeName]
      if (node !== undefined && node.k8s_node !== undefined) {
        node = node.k8s_node
        for (let key in node.spec.taints) {
          let taint = node.spec.taints[key]
          if (taint.effect === "NoSchedule" && taint.key === "node-role.kubernetes.io/master") {
            return false
          }
        }
        return true
      }
      return false
    }
  },
  components: { StateNodeCondition, NodeRoleTag },
  mounted () {
  },
  methods: {

  }
}
</script>

<style scoped lang="scss">
.info {
  background-color: var(--el-background-color-base);
  margin-bottom: 10px;
  padding: 5px 20px;
}
</style>

