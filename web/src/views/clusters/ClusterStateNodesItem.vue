<i18n>

</i18n>


<template>
  <div class="nodeInfo app_text_mono">
    <div style="display: flex;">
      <div style="min-width: 120px;">
        {{name}}
      </div>
      <div style="flex-grow: 1;">
        <template v-if="node.k8s_node">
          <el-tag type="primary" effect="dark" style="margin-left: 10px;" v-if="node.k8s_node.metadata.labels['node-role.kubernetes.io/control-plane'] !== undefined">{{$t('node.kube_control_plane')}}</el-tag>
          <el-tag type="success" effect="dark" style="margin-left: 10px;" v-if="isKubeNode">{{$t('node.kube_node')}}</el-tag>
        </template>
        <template v-if="node.etcd_member">
          <el-tag type="warning" effect="dark" style="margin-left: 10px; padding: 0 8px;">
            {{node.etcd_member.name}}
          </el-tag>
          <!-- <el-tag type="primary" effect="light" style="margin-left: 20px;">{{ node.etcd_member.clientURLs && node.etcd_member.clientURLs.length > 0 ? node.etcd_member.clientURLs[0] : '' }}</el-tag> -->
        </template>
        <template v-if="node.k8s_node">
          <template v-for="(addr, index) in node.k8s_node.status.addresses" :key="name + index">
            <el-tag type="primary" style="float: right;" v-if="addr.type === 'InternalIP'">{{addr.address}}</el-tag>
          </template>
        </template>
      </div>
    </div>
    <div v-if="node.etcd_member && !node.etcd_member.health.health" class="etcd_error">
      etcd endpoint unhealthy: {{node.etcd_member.health.endpoint}} {{node.etcd_member.health.error}}
    </div>
    <div>
      <div v-for="(condition, index) in node.k8s_node.status.conditions" :key="'condition' + index" style="margin-left: 130px;">
        <StateNodeCondition :condition="condition" hideSuccess></StateNodeCondition>
      </div>
    </div>
  </div>
</template>

<script>
import StateNodeCondition from './plan/node/StateNodeCondition.vue'

export default {
  props: {
    node: { type: Object, required: true },
    name: { type: String, required: true },
  },
  data() {
    return {

    }
  },
  computed: {
    isKubeNode () {
      let node = this.node
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
  components: { StateNodeCondition },
  mounted () {
  },
  methods: {

  }
}
</script>

<style scoped lang="scss">
.nodeInfo {
  padding: 10px 20px;
  margin-bottom: 10px;
  background-color: var(--el-color-info-light);
  .etcd_error {
    font-size: 12px;
    padding: 2px 10px;
    background-color: var(--el-color-danger);
    margin-top: 10px;
    color: var(--el-color-white);
    margin-left: 130px;
  }
}
</style>
