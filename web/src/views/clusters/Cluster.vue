<i18n>
en:
  cluster: Cluster
  clusterList: Clusters List
  plan: Cluster Plan
  access: Access Cluster
  node: Nodes Maintainance
zh:
  clusterList: 集群列表
  cluster: 集群
  access: 访问集群
  plan: 集群规划
  nodes: 节点维护
</i18n>

<template>
  <div>
    <ControlBar :title="name">
      <span v-show="cluster && !cluster.history.processing" style="margin-right: 10px;">
        <template v-if="currentTab === 'plan'">
          <template v-if="mode === 'view' && (isClusterOnline || !isClusterInstalled)">
            <el-button type="primary" icon="el-icon-edit" @click="$router.replace(`/clusters/${name}?mode=edit`)">{{$t('msg.edit')}}</el-button>
          </template>
          <template v-if="mode === 'edit'">
            <el-popconfirm :confirm-button-text="$t('msg.ok')" :cancel-button-text="$t('msg.cancel')" placement="bottom-start"
              icon="el-icon-warning" icon-color="red" :title="$t('msg.confirmToCancel')" @confirm="cancelEdit">
              <template #reference>
                <el-button type="default" icon="el-icon-close">{{$t('msg.cancel')}}</el-button>
              </template>
            </el-popconfirm>
            <el-button type="primary" icon="el-icon-check" :disabled="noSaveRequired" @click="save">{{$t('msg.save')}}</el-button>
          </template>
        </template>
        <template v-if="currentTab === 'access' || currentTab === 'plan'">
          <ClusterProcessing v-if="mode === 'view' && cluster && !cluster.history.processing" :cluster="cluster" :name="name" @refresh="refresh" 
            :loading="loading"></ClusterProcessing>
        </template>
      </span>
      <template v-if="cluster && cluster.history.processing">
        <ClusterProcessing v-if="mode === 'view'" :cluster="cluster" :name="name" @refresh="refresh" :loading="loading"></ClusterProcessing>
      </template>
      <template v-if="cluster && cluster.state">
        <ClusterStateNodes :state="cluster.state"></ClusterStateNodes>
      </template>
    </ControlBar>
    <el-card shadow="none" v-if="loading">
      <el-skeleton animated :rows="10"></el-skeleton>
    </el-card>
    <el-tabs type="border-card" v-show="!loading" v-model="currentTab">
      <el-tab-pane :label="$t('plan')" name="plan">
        <Plan v-if="cluster" ref="plan" :cluster="cluster" :mode="mode"></Plan>
      </el-tab-pane>
      <el-tab-pane :label="$t('access')" name="access" :disabled="disableNonePlanTab">
        <Access v-if="cluster && cluster.history.success_tasks.length > 0" ref="access" :cluster="cluster" :loading="loading" @switch="currentTab = $event"></Access>
      </el-tab-pane>
      <el-tab-pane :disabled="disableNonePlanTab || !isClusterOnline" label="健康检查">
        <el-alert>检查集群当前的状况与集群安装计划的匹配情况，正在建设...</el-alert>
      </el-tab-pane>
      <el-tab-pane :disabled="disableNonePlanTab || !isClusterOnline" label="备份/恢复">
        <el-alert>对 etcd 及 control-plane 做备份及恢复，正在建设...</el-alert>      
      </el-tab-pane>
      <el-tab-pane :disabled="disableNonePlanTab || !isClusterOnline" label="CIS扫描">
        <el-alert>CIS 扫描，正在建设...</el-alert>
      </el-tab-pane>
      <el-tab-pane :disabled="disableNonePlanTab || !isClusterOnline" label="升级包检测">
        <el-alert>检测 kuboard-spray 升级资源包，正在建设...</el-alert>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script>
import mixin from '../../mixins/mixin.js'
import Plan from './plan/Plan.vue'
import yaml from 'js-yaml'
import ClusterProcessing from './ClusterProcessing.vue'
import Access from './access/Access.vue'
import { computed } from 'vue'
import ClusterStateNodes from './ClusterStateNodes.vue'
import clone from 'clone'

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
    isClusterInstalled () {
      return !(this.cluster && this.cluster.history.success_tasks.length == 0)
    },
    isClusterOnline () {
      if (this.cluster) {
        if (this.cluster.state && this.cluster.state.code === 200) {
          return true
        }
      }
      return false
    },
    disableNonePlanTab () {
      if (this.mode !== 'view') {
        return true
      }
      return !this.isClusterInstalled
    }
  },
  provide () {
    return {
      isClusterInstalled: computed(() => {
        return this.isClusterInstalled
      }),
      isClusterOnline: computed(() => {
        return this.isClusterOnline
      }),
      pendingRemoveNodes: computed(() => {
        return this.pendingNodes('remove_node')
      }),
      pendingAddNodes: computed(() => {
        return this.pendingNodes('add_node')
      }),
    }
  },
  components: { Plan, ClusterProcessing, Access, ClusterStateNodes },
  watch: {
    'cluster.inventory.all.hosts.localhost.kuboardspray_resource_package': function() {
      this.loadResourcePackage()
    }
  },
  mounted () {
    this.refresh()
  },
  methods: {
    pendingNodes (action) {
      let result = []
        if (this.isClusterInstalled) {
          for (let key in this.cluster.inventory.all.hosts) {
            let host = this.cluster.inventory.all.hosts[key]
            if (host.kuboardspray_node_action === action) {
              let h = clone(host)
              h.name = key
              if (this.cluster.inventory.all.children.target.children.k8s_cluster.children.kube_control_plane.hosts[key]) {
                h.isControlPlane = true
              }
              if (this.cluster.inventory.all.children.target.children.k8s_cluster.children.kube_node.hosts[key]) {
                h.isNode = true
              }
              if (this.cluster.inventory.all.children.target.children.etcd.hosts[key]) {
                h.isEtcd = true
              }
              result.push(h)
            }
          }
        }
        return result
    },
    cancelEdit () {
      this.$router.replace(`/clusters/${this.name}`)
      this.refresh()
    },
    async refresh() {
      this.loading = true
      await this.kuboardSprayApi.get(`/clusters/${this.name}`).then(resp => {
        this.cluster = resp.data.data
        // this.cluster.processing = true
        if (this.cluster.history.success_tasks.length > 0) {
          this.currentTab = 'access'
        } else {
          this.currentTab = 'plan'
        }
        this.originalInventoryYaml = yaml.dump(this.cluster.inventory)
        this.loadResourcePackage()
        setTimeout(() => {
          if (this.$refs.plan) {
            this.$refs.plan.ping()
          }
        }, 200)
        if (this.isClusterInstalled) {
          this.loadStateNodes()
        }
      }).catch(e => {
        console.log(e.response)
      })
      this.loading = false
    },
    loadStateNodes() {
      this.kuboardSprayApi.get(`/clusters/${this.name}/state/nodes`).then(resp => {
        if (resp.data.data.stdout_obj && resp.data.data.stdout_obj.items) {
          this.cluster.state = { nodes: {}, code: 200 }
          for (let item of resp.data.data.stdout_obj.items) {
            this.cluster.state.nodes[item.metadata.name] = item
          }
        }
        if (resp.data.data.return_code === '' && resp.data.data.stdout_obj.changed === false && resp.data.data.stdout_obj.msg) {
          this.cluster.state = {
            code: 500,
            msg: resp.data.data.stdout_obj.msg
          }
        }
      }).catch(e => {
        this.cluster.state = {
          code: 500,
          msg: e.response.data.message
        }
      })
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
