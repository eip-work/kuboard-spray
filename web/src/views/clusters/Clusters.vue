<i18n>
en:
  clusters: Cluster Management
  clusterList: Clusters List
  addCluster: Add Cluster Installation Plan
zh:
  clusters: 集群管理
  clusterList: 集群列表
  addCluster: 添加集群安装计划
</i18n>

<template>
  <div>
    <div class="app_block_title">
      {{$t('obj.cluster')}}
    </div>
    <div class="app_margin_bottom">
      <el-alert :closable="false" title="集群管理" type="default">
        <div style="line-height: 20px;">
          您可以：
          <li>制定新的集群安装计划</li>
          <li>维护使用 KuboardSpray 安装的集群</li>
        </div>
      </el-alert>
    </div>
    <el-card shadow="none">
      <div style="display: flex; flex-wrap: wrap;">
        <el-card v-for="(item, index) in clusters" :key="'cluster' + index" shadow="none" class="cluster" 
          @click="$router.push(`/clusters/${item}`)">
          <div class="noselect">
            {{item}}
          </div>
        </el-card>
        <el-button type="primary" size="large" icon="el-icon-plus">{{$t('addCluster')}}</el-button>
      </div>
    </el-card>
  </div>
</template>

<script>
import mixin from '../../mixins/mixin.js'

export default {
  mixins: [mixin],
  props: {
  },
  percentage () {
    return 100
  },
  breadcrumb () {
    return [
      { label: this.$t('clusterList') }
    ]
  },
  data() {
    return {
      clusters: []
    }
  },
  computed: {
  },
  components: { },
  mounted () {
    this.kuboardSprayApi.get('/clusters').then(resp => {
      this.clusters = resp.data.data
    }).catch(e => {
      this.$message.error('Error: ' + e)
    })
  },
  methods: {

  }
}
</script>

<style scoped lang="scss">
.cluster {
  margin-right: 10px;
  width: 200px;
  border-radius: 6px;
  cursor: pointer;
}
.cluster:hover {
  border-color: $--color-primary;
}
</style>
