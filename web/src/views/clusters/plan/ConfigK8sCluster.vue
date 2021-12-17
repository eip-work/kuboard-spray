<i18n>
en:
  sshcommon: SSH Shared Params (apply to all the k8s nodes)
zh:
  sshcommon: SSH 共享参数（适用范围：所有 k8s 节点）
</i18n>


<template>
  <div>
    <SshParamsCluster :inventory="inventory" :clusterName="clusterName" :holder="inventory.all.children.k8s_cluster.vars" prop="all.children.k8s_cluster.vars" :description="$t('sshcommon')"></SshParamsCluster>
    <K8sCluster :resourcePackage="resourcePackage" :inventory="inventory"></K8sCluster>
    <NetworkPlugin :resourcePackage="resourcePackage" :inventory="inventory"></NetworkPlugin>
    <AddonMetricsServer :resourcePackage="resourcePackage" :inventory="inventory"></AddonMetricsServer>
  </div>
</template>

<script>
import SshParamsCluster from './common/SshParamsCluster.vue'
import AddonMetricsServer from './k8s_cluster/AddonMetricsServer.vue'
import NetworkPlugin from './k8s_cluster/NetworkPlugin.vue'
import K8sCluster from './k8s_cluster/K8sCluster.vue'

export default {
  props: {
    clusterName: { type: String, required: true },
    inventory: { type: Object, required: true },
    resourcePackage: { type: Object, required: false, default: undefined },
  },
  data() {
    return {

    }
  },
  computed: {
  },
  components: { SshParamsCluster, AddonMetricsServer, NetworkPlugin, K8sCluster },
  mounted () {
  },
  methods: {
  }
}
</script>

<style scoped lang="scss">

</style>
