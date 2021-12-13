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
    <ControlBar :title="name">
      <!-- <el-select style="margin-right: 10px;"></el-select> -->
      <el-button type="success">{{name}}</el-button>
      <el-button type="primary">History</el-button>
    </ControlBar>
    <el-tabs type="border-card">
      <el-tab-pane :label="$t('plan')">
        <Plan v-if="cluster" :cluster="cluster"></Plan>
      </el-tab-pane>
      <!-- <el-tab-pane :label="$t('nodes')">正在建设...</el-tab-pane> -->
      <el-tab-pane label="健康检查">检查集群的状态与集群规划内容的匹配情况（正在建设...）</el-tab-pane>
      <el-tab-pane label="备份">备份 etcd 内容（正在建设...）</el-tab-pane>
      <el-tab-pane label="CIS扫描">
        <li>https://github.com/aquasecurity/kube-bench</li>
      </el-tab-pane>
      <el-tab-pane label="升级包检测">检查是否有更新的 KuboardSpray 资源包（正在建设...）</el-tab-pane>
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
  watch: {
    'cluster.inventory.all.hosts.localhost.kuboardspray_resource_package': function() {
      this.loadResourcePackage()
    }
  },
  mounted () {
    this.kuboardSprayApi.get(`/clusters/${this.name}`).then(resp => {
      this.cluster = resp.data.data
      this.loadResourcePackage()
    }).catch(e => {
      console.log(e.response)
    })
  },
  methods: {
    loadResourcePackage () {
      this.cluster.resourcePackage = undefined
      let newValue = this.cluster.inventory.all.hosts.localhost.kuboardspray_resource_package
      if (newValue) {
        this.kuboardSprayApi.get(`/resources/${newValue}`).then(resp => {
          this.cluster.resourcePackage = resp.data.data.package
        }).catch(e => {
          console.log(e)
        })
      }
    }
  }
}
</script>

<style scoped lang="scss">

</style>
