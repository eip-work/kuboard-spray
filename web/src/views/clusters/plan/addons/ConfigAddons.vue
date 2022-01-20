
<template>
  <div v-if="cluster.resourcePackage">
    <el-button @click="checkStatus">检查状态</el-button>
    <!-- <NetChecker :cluster="cluster"></NetChecker>
    <MetricsServer :cluster="cluster"></MetricsServer> -->
    <template v-for="(item, index) in cluster.resourcePackage.data.addon" :key="'addon' + index">
      <Component :is="item.name" :cluster="cluster"></Component>
    </template>
  </div>
</template>

<script>
import metrics_server from './addons/metrics_server.vue'
import netchecker from './addons/netchecker.vue'
import nodelocaldns from './addons/nodelocaldns.vue'

export default {
  props: {
    cluster: { type: Object, required: true },
    currentTab: { type: String, required: false, default: '' },
  },
  data() {
    return {

    }
  },
  computed: {
  },
  components: { metrics_server, netchecker, nodelocaldns },
  watch: {
    currentTab (newValue) {
      if (newValue === 'addons') {
        console.log('activate')
      }
    }
  },
  mounted () {
  },
  methods: {
    checkStatus () {
      this.kuboardSprayApi.get(`/clusters/${this.cluster.name}/state/addons`).then(resp => {
        console.log(resp.data)
      }).catch(e => {
        console.log(e)
      })
    },
  }
}
</script>

<style scoped lang="scss">

</style>
