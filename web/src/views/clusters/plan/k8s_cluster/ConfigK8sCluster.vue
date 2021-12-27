<template>
  <div>
    <K8sCluster :cluster="cluster"></K8sCluster>
    <NetworkPlugin :cluster="cluster"></NetworkPlugin>
  </div>
</template>

<script>
import NetworkPlugin from './NetworkPlugin.vue'
import K8sCluster from './K8sCluster.vue'

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
  components: { NetworkPlugin, K8sCluster },
  mounted () {
  },
  methods: {
  }
}
</script>

<style scoped lang="scss">

</style>
