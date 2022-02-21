
<template>
  <div v-if="cluster.resourcePackage">
    <template v-for="(item, index) in cluster.resourcePackage.data.addon" :key="'addon' + index">
      <Component :is="item.name" :cluster="cluster" :pingpong="pingpong" @refresh="$emit('refresh')"></Component>
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
    pingpong: { prop: Object, required: false, default: () => { return {} }},
  },
  data() {
    return {

    }
  },
  computed: {
  },
  components: { metrics_server, netchecker, nodelocaldns },
  watch: {
    // currentTab (newValue) {
    //   if (newValue === 'addons') {
    //     console.log('activate')
    //   }
    // }
  },
  mounted () {
  },
  methods: {
  }
}
</script>

<style scoped lang="css">

</style>
