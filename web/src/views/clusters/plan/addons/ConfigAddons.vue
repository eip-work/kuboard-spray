
<template>
  <div>
    <AddonMetricsServer :cluster="cluster"></AddonMetricsServer>
  </div>
</template>

<script>
import AddonMetricsServer from './AddonMetricsServer.vue'

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
  components: { AddonMetricsServer },
  mounted () {
  },
  methods: {
  }
}
</script>

<style scoped lang="scss">

</style>
