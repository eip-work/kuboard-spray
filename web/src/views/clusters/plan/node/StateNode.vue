<i18n>
en:
  label: Node Info
  description: Node info retrived from kubernetes cluster.
zh:
  label: 节点信息
  description: 从 K8S 获取到的节点信息

</i18n>

<template>
  <ConfigSection v-model:enabled="enabled" :label="$t('label')" :description="$t('description', {nodeName: nodeName})" disabled anti-freeze>
    <el-form label-width="160px" class="app_form_mini">
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
    }
  },
  components: { StateNodeCondition },
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

