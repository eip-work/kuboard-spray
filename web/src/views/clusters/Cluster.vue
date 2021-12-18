<i18n>
en:
  cluster: Cluster
  clusterList: Clusters List
  plan: Cluster Plan
  node: Nodes Maintainance
  confirmToCancel: Modifications will be lost to proceed, do you confirm ?
zh:
  clusterList: 集群列表
  cluster: 集群
  plan: 集群规划
  nodes: 节点维护
  confirmToCancel: 将丢失已修改内容，确认取消编辑？
</i18n>

<template>
  <div>
    <ControlBar :title="name">
      <div v-show="cluster && !cluster.processing">
      <!-- <el-select style="margin-right: 10px;"></el-select> -->
        <template v-if="currentTab === 'plan'">
          <template v-if="mode === 'view'">
            <el-button type="primary" icon="el-icon-edit" @click="$router.replace(`/clusters/${name}?mode=edit`)">{{$t('msg.edit')}}</el-button>
            <ClusterProcessing v-if="cluster" :cluster="cluster" :name="name" @refresh="refresh" :loading="loading"></ClusterProcessing>
          </template>
          <template v-if="mode === 'edit'">
            <el-popconfirm :confirm-button-text="$t('msg.ok')" :cancel-button-text="$t('msg.cancel')" placement="bottom-start"
              icon="el-icon-warning" icon-color="red" :title="$t('confirmToCancel')" @confirm="cancelEdit">
              <template #reference>
                <el-button type="default" icon="el-icon-close">{{$t('msg.cancel')}}</el-button>
              </template>
            </el-popconfirm>
            <el-button type="primary" icon="el-icon-check" :disabled="noSaveRequired" @click="save">{{$t('msg.save')}}</el-button>
          </template>
          <template v-if="mode === 'create'">
            <el-button type="primary" icon="el-icon-check">{{$t('msg.save')}}</el-button>
          </template>
        </template>
      </div>
    </ControlBar>
    <el-card shadow="none" v-if="loading">
      <el-skeleton animated :rows="10"></el-skeleton>
    </el-card>
    <el-tabs type="border-card" v-show="!loading" v-model="currentTab">
      <el-tab-pane :label="$t('plan')" name="plan">
        <Plan v-if="cluster" ref="plan" :cluster="cluster" :mode="mode"></Plan>
      </el-tab-pane>
      <el-tab-pane disabled label="健康检查">检查集群的状态与集群规划内容的匹配情况（正在建设...）</el-tab-pane>
      <el-tab-pane disabled label="备份">备份 etcd 内容（正在建设...）</el-tab-pane>
      <el-tab-pane disabled label="CIS扫描">
        <li>https://github.com/aquasecurity/kube-bench</li>
      </el-tab-pane>
      <el-tab-pane disabled label="升级包检测">检查是否有更新的 KuboardSpray 资源包（正在建设...）</el-tab-pane>
    </el-tabs>
  </div>
</template>

<script>
import mixin from '../../mixins/mixin.js'
import Plan from './plan/Plan.vue'
import yaml from 'js-yaml'
import ClusterProcessing from './ClusterProcessing.vue'
import { computed } from 'vue'

export default {
  mixins: [mixin],
  props: {
    name: { type: String, required: true },
    mode: { type: String, required: false, default: 'view' },
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
  refresh () {
    this.refresh()
  },
  data () {
    return {
      loading: false,
      cluster: undefined,
      originalInventoryYaml: '',
      currentTab: 'plan'
    }
  },
  computed: {
    noSaveRequired () {
      if (this.cluster === undefined) {
        return true
      } 
      return this.originalInventoryYaml == yaml.dump(this.cluster.inventory)
    },
  },
  provide () {
    return {
      isInstalled: computed(() => {
        if (this.cluster === undefined) {
          return false
        }
        if (this.cluster.success_tasks.length > 0) {
          return true
        }
        return false
      })
    }
  },
  components: { Plan, ClusterProcessing },
  watch: {
    'cluster.inventory.all.hosts.localhost.kuboardspray_resource_package': function() {
      this.loadResourcePackage()
    }
  },
  mounted () {
    this.refresh()
  },
  methods: {
    cancelEdit () {
      this.$router.replace(`/clusters/${this.name}`)
      this.refresh()
    },
    async refresh() {
      this.loading = true
      await this.kuboardSprayApi.get(`/clusters/${this.name}`).then(resp => {
        this.cluster = resp.data.data
        // this.cluster.processing = true
        this.originalInventoryYaml = yaml.dump(this.cluster.inventory)
        this.loadResourcePackage()
      }).catch(e => {
        console.log(e.response)
      })
      setTimeout(() => {
        this.loading = false
      }, 200)
    },
    async loadResourcePackage () {
      this.cluster.resourcePackage = undefined
      let newValue = this.cluster.inventory.all.hosts.localhost.kuboardspray_resource_package
      if (newValue) {
        await this.kuboardSprayApi.get(`/resources/${newValue}`).then(resp => {
          this.cluster.resourcePackage = resp.data.data.package
        }).catch(e => {
          console.log(e)
        })
      }
    },
    save () {
      this.$refs.plan.validate(flag => {
        if (flag) {
          this.kuboardSprayApi.put(`/clusters/${this.name}`, this.cluster.inventory).then(() => {
            this.$message.success(this.$t('msg.save_succeeded'))
            this.refresh()
            this.$router.replace(`/clusters/${this.name}`)
          }).catch(e => {
            this.$message.error(this.$t('msg.save_failed', e.response.data.message))
          })
        }
      })
    },
  }
}
</script>

<style scoped lang="scss">

</style>
