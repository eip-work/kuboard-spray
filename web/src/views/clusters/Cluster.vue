<i18n>
en:
  cluster: Cluster
  clusterList: Clusters List
  plan: Cluster Plan
  node: Nodes Maintainance
zh:
  clusterList: 集群列表
  cluster: 集群
  plan: 集群规划
  nodes: 节点维护
</i18n>

<template>
  <div>
    <div class="app_block_title">{{$t('cluster')}} {{name}}</div>
    <el-tabs type="border-card">
      <el-tab-pane :label="$t('plan')">
        <Plan v-if="cluster" :cluster="cluster"></Plan>
      </el-tab-pane>
      <el-tab-pane :label="$t('nodes')">正在建设...</el-tab-pane>
    </el-tabs>
  </div>
</template>

<script>
import mixin from '../../mixins/mixin.js'
import Plan from './plan/Plan.vue'

export default {
  mixins: [mixin],
  props: {
    name: { type: String, required: true },
  },
  percentage () {
    return this.loading ? 10 : 100
  },
  breadcrumb () {
    return [
      { label: this.$t('clusterList'), to: '/clusters' },
      { label: this.name }
    ]
  },
  data () {
    return {
      loading: false,
      cluster: undefined,
    }
  },
  computed: {
  },
  components: { Plan },
  mounted () {
    this.kuboardSprayApi.get(`/clusters/${this.name}`).then(resp => {
      this.cluster = resp.data.data
    }).catch(e => {
      console.log(e.response)
    })
  },
  methods: {

  }
}
</script>

<style scoped lang="scss">

</style>
