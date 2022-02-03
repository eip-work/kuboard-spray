<template>
  <div>
    <K8sCluster :cluster="cluster"></K8sCluster>
    <Networking :cluster="cluster"></Networking>
    <NetworkPlugin :cluster="cluster"></NetworkPlugin>
  </div>
</template>

<script>
import NetworkPlugin from './NetworkPlugin.vue'
import K8sCluster from './K8sCluster.vue'
import Networking from './Networking.vue'

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
  components: { NetworkPlugin, K8sCluster, Networking },
  mounted () {
  },
  methods: {
  }
}
</script>

<style scoped lang="css">

</style>
