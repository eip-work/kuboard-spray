<i18n>
en:
  clusters: Cluster Management
  clusterList: Clusters List
  addCluster: Add Cluster Installation Plan
  confirmToDelete: Do you confirm to delete the cluster info in KuboardSpray? Cluster itself will not be impacted.
  rename: 重命名
zh:
  clusters: 集群管理
  clusterList: 集群列表
  addCluster: 添加集群安装计划
  confirmToDelete: 是否删除此集群在 KuboardSpray 的信息？（该集群本身将继续运行不受影响。）
  rename: 重命名
</i18n>

<template>
  <div>
    <div class="app_block_title">
      {{$t('obj.cluster')}}
    </div>
    <div class="app_margin_bottom">
      <el-alert :closable="false" title="集群管理" class="app_white_alert">
        <div style="line-height: 20px;">
          您可以：
          <li>制定新的集群安装计划，并安装 kubernetes 集群</li>
          <li>维护使用 KuboardSpray 安装的集群</li>
        </div>
      </el-alert>
    </div>
    <el-card shadow="none" style="min-height: 234px;">
      <el-skeleton v-if="loading" :rows="5" animated />
      <div v-else style="display: flex; flex-wrap: wrap;">
        <div v-for="(item, index) in clusters" :key="'cluster' + index" class="cluster">
          <div class="deleteButton">
            <ClustersRename :clusterName="item" @success="refresh"></ClustersRename>
            <el-popconfirm :confirm-button-text="$t('msg.ok')" :cancel-button-text="$t('msg.cancel')" icon="el-icon-warning" icon-color="red"
              placement="top-start" :title="$t('confirmToDelete')" @confirm="deleteCluster(item)">
              <template #reference>
                <el-button icon="el-icon-delete" circle type="danger"></el-button>
              </template>
            </el-popconfirm>
          </div>
          <el-card shadow="none"
            @click="$router.push(`/clusters/${item}`)">
            <div class="noselect">
              {{item}}
            </div>
          </el-card>
        </div>
        <el-button class="cluster" style="margin-bottom: 10px;" type="primary" size="large" icon="el-icon-plus" @click="$refs.create.show()">{{$t('addCluster')}}</el-button>
      </div>
    </el-card>
    <CreateCluster ref="create"></CreateCluster>
  </div>
</template>

<script>
import mixin from '../../mixins/mixin.js'
import CreateCluster from './create/CreateCluster.vue'
import ClustersRename from './create/ClustersRename.vue'

export default {
  mixins: [mixin],
  props: {
  },
  percentage () {
    return this.loading ? 10 : 100
  },
  breadcrumb () {
    return [
      { label: this.$t('clusterList') }
    ]
  },
  refresh () {
    this.refresh()
  },
  data() {
    return {
      clusters: [],
      loading: false,
    }
  },
  computed: {
  },
  components: { CreateCluster, ClustersRename },
  mounted () {
    this.refresh()
  },
  methods: {
    async refresh () {
      this.loading = true
      await this.kuboardSprayApi.get('/clusters').then(resp => {
        this.clusters = resp.data.data
      }).catch(e => {
        this.$message.error('Error: ' + e)
      })
      this.loading = false
    },
    deleteCluster(cluster) {
      this.kuboardSprayApi.delete('/clusters/' + cluster).then(() => {
        this.refresh()
        this.$message.success(this.$t('msg.delete_succeeded'))
      }).catch(e => {
        this.$message.error(this.$t('msg.delete_failed', {msg: e.response.data.message }))
      })
    }
  }
}
</script>

<style scoped lang="scss">
.cluster {
  margin-right: 10px;
  margin-bottom: 10px;
  height: 60px;
  width: 200px;
  border-radius: 6px;
  cursor: pointer;
  .deleteButton {
    height: 0px;
    overflow: hidden;
    position: relative;
    top: -5px;
    left: 5px;
    text-align: right;
  }
}
.cluster:hover {
  border-color: var(--el-color-primary);
}
.cluster:hover .deleteButton {
  overflow: visible;
}
</style>
