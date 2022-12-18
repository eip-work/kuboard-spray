<template>
  <div>
    <K8sCluster :cluster="cluster"></K8sCluster>
    <RotateCertificate :cluster="cluster" @switchTab="$emit('switchTab', $event)"></RotateCertificate>
    <Networking :cluster="cluster"></Networking>
    <NetworkPlugin :cluster="cluster"></NetworkPlugin>
    <NetworkDNS :cluster="cluster"></NetworkDNS>
    <AuditLog :cluster="cluster"></AuditLog>
  </div>
</template>

<script>
import NetworkPlugin from './NetworkPlugin.vue'
import K8sCluster from './K8sCluster.vue'
import Networking from './Networking.vue'
import NetworkDNS from './NetworkDNS.vue'
import AuditLog from './AuditLog.vue'
import RotateCertificate from './RotateCertificate.vue'

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
  components: { NetworkPlugin, K8sCluster, Networking, NetworkDNS, AuditLog, RotateCertificate },
  mounted () {
  },
  methods: {
  }
}
</script>

<style scoped lang="css">

</style>
