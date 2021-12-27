
<template>
  <div>
    <NetChecker :cluster="cluster"></NetChecker>
    <MetricsServer :cluster="cluster"></MetricsServer>
  </div>
</template>

<script>
import MetricsServer from './MetricsServer.vue'
import NetChecker from './NetChecker.vue'

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
  components: { MetricsServer, NetChecker },
  mounted () {
  },
  methods: {
  }
}
</script>

<style scoped lang="scss">

</style>
