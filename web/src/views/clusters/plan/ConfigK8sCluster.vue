<i18n>
en:
  sshcommon: SSH Shared Params (apply to all the k8s nodes)
zh:
  sshcommon: SSH 共享参数（适用范围：所有 k8s 节点）
</i18n>


<template>
  <div>
    <SshParamsCluster :cluster="cluster" :holder="cluster.inventory.all.children.target.children.k8s_cluster.vars" prop="all.children.target.children.k8s_cluster.vars" :description="$t('sshcommon')"></SshParamsCluster>
    <OsMirror :cluster="cluster"></OsMirror>
    <K8sCluster :cluster="cluster"></K8sCluster>
    <NetworkPlugin :cluster="cluster"></NetworkPlugin>
    <AddonMetricsServer :cluster="cluster"></AddonMetricsServer>
  </div>
</template>

<script>
import SshParamsCluster from './common/SshParamsCluster.vue'
import AddonMetricsServer from './k8s_cluster/AddonMetricsServer.vue'
import NetworkPlugin from './k8s_cluster/NetworkPlugin.vue'
import K8sCluster from './k8s_cluster/K8sCluster.vue'
import OsMirror from './k8s_cluster/OsMirror.vue'

export default {
  props: {
    cluster: { type: Object, required: true },
  },
  data() {
    return {

    }
  },
  computed: {
  },
  components: { SshParamsCluster, AddonMetricsServer, NetworkPlugin, K8sCluster, OsMirror },
  mounted () {
  },
  methods: {
  }
}
</script>

<style scoped lang="scss">

</style>
