<i18n>
en:
  singleNode: Specific Node
  global_config: Global Config
  addons: Addons
  enabledBation: Enabled
  disabledBation: Disabled
  selectANode: Please select a node from the diagram to the left.
  resourcePackage: Resource Package
zh:
  singleNode: 单个节点
  global_config: 全局设置
  addons: 可选组件
  enabledBation: 使用跳板机
  disabledBation: 不使用跳板机
  selectANode: 请从左侧图中选择一个节点
  resourcePackage: 资源包
</i18n>

<template>
  <div style="display:flex;">
    <div class="plan">
      <div class="left">
        <div style="padding: 5px; font-weight: bolder; font-size: 14px;">
          Kuboard Spray
        </div>
        <div>
          <Node class="localhost" name="localhost" :cluster="cluster" hideDeleteButton
            :active="currentPropertiesTab === 'localhost'" @click="currentPropertiesTab = 'localhost'">
          </Node>
        </div>
        <div>
          <div v-if="bastionEnabled" class="verticalConnection"></div>
          <div v-else class="horizontalConnection"></div>
        </div>
        <div>
          <Node class="bastion" name="bastion" :cluster="cluster" hideDeleteButton
            :active="currentPropertiesTab === 'bastion'" @click="showBastion">
            <div style="margin-top: 10px;">
              <el-tag v-if="bastionEnabled" type="danger" effect="dark" size="small" style="width: 100px; text-align: center;">{{$t('enabledBation')}}</el-tag>
              <el-tag v-else type="info" effect="light" size="small" style="width: 100px; text-align: center;">{{$t('disabledBation')}}</el-tag>
            </div>
          </Node>
          <div class="horizontalConnection" :style="bastionEnabled ? '' : 'border-color: white;'"></div>
        </div>
        <div style="line-height: 28px;">
          <KuboardSprayLink href="https://kuboard-spray.cn/guide/install-k8s.html" :size="12">安装 K8S 集群？</KuboardSprayLink>
          <KuboardSprayLink href="https://kuboard-spray.cn/guide/maintain/ha-mode.html" :size="12">实现高可用？</KuboardSprayLink>
          <KuboardSprayLink href="https://kuboard-spray.cn/guide/maintain/add-replace-node.html" :size="12">添加删除节点？</KuboardSprayLink>
          <KuboardSprayLink href="https://kuboard-spray.cn/guide/maintain/upgrade.html" :size="12">升级集群？</KuboardSprayLink>
        </div>
      </div>
      <div class="right">
        <div style="padding: 5px; font-weight: bolder; font-size: 14px; height: 28px; line-height: 28px; margin-bottom: 10px;">
          <span style="margin-right: 20px;">Kubernetes Cluster</span>
          <AddNode :inventory="inventory" v-model:currentPropertiesTab="currentPropertiesTab"></AddNode>
          <el-button v-if="mode == 'view'" type="primary" plain icon="el-icon-lightning" @click="ping" :loading="pingpong_loading">
            <span class="app_text_mono">PING</span>
          </el-button>
        </div>
        <el-scrollbar height="calc(100vh - 293px)">
          <div class="masters">
            <Node v-for="(item, index) in inventory.all.children.target.children.k8s_cluster.children.kube_control_plane.hosts" :key="'control_plane' + index" @deleted="currentPropertiesTab = 'node_nodes'"
              @click="currentPropertiesTab = 'NODE_' + index" :pingpong="pingpong" :pingpong_loading="pingpong_loading"
              :active="nodeRoles(index)[currentPropertiesTab] || currentPropertiesTab === 'global_config' || currentPropertiesTab === 'addons' || currentPropertiesTab === 'k8s_cluster' || 'NODE_' + index === currentPropertiesTab"
              :name="index" :cluster="cluster"></Node>
            <template v-for="(item, index) in inventory.all.children.target.children.etcd.hosts" :key="'etcd' + index">
              <Node v-if="isEtcdAndNotControlPlane(index)" :name="index" :cluster="cluster" @deleted="currentPropertiesTab = 'node_nodes'"
                @click="currentPropertiesTab = 'NODE_' + index" @delete_button="deleteNode(index)" :pingpong="pingpong" :pingpong_loading="pingpong_loading"
                :active="nodeRoles(index)[currentPropertiesTab] || currentPropertiesTab === 'global_config' || currentPropertiesTab === 'k8s_cluster' || 'NODE_' + index === currentPropertiesTab"></Node>
            </template>
          </div>
          <div class="workers">
            <template v-for="(item, index) in inventory.all.children.target.children.k8s_cluster.children.kube_node.hosts" :key="'node' + index">
              <Node v-if="isNode(index)" :name="index" :cluster="cluster" :pingpong="pingpong" :pingpong_loading="pingpong_loading"
                @click="currentPropertiesTab = 'NODE_' + index" @deleted="currentPropertiesTab = 'node_nodes'"
                :active="nodeRoles(index)[currentPropertiesTab] || currentPropertiesTab === 'global_config' || currentPropertiesTab === 'addons' || currentPropertiesTab === 'k8s_cluster' || 'NODE_' + index === currentPropertiesTab"></Node>
            </template>
          </div>
          <div class="workers">
            <template v-for="(item, index) in nodeGap.inventory.all.hosts" :key="'gap_node' + index">
              <Node :name="index" :cluster="nodeGap" :pingpong="pingpong" :pingpong_loading="pingpong_loading" @deleted="currentPropertiesTab = 'node_nodes'"
                @click="currentPropertiesTab = 'GAP_NODE_' + index"
                :active="'GAP_NODE_' + index === currentPropertiesTab"></Node>
            </template>
          </div>
        </el-scrollbar>
      </div>
    </div>
    <div class="properties">
      <el-form ref="form" label-width="120px" label-position="left" @submit.enter.prevent :model="inventory">
        <el-tabs type="card" v-model="currentPropertiesTab">
          <el-tab-pane name="localhost">
            <template #label>
              {{ $t('resourcePackage') }}
            </template>
            <el-scrollbar max-height="calc(100vh - 276px)" ref="configKuboardSprayScroll">
              <div class="tab_content">
                <ConfigKuboardSpray :cluster="cluster"></ConfigKuboardSpray>
              </div>
            </el-scrollbar>
          </el-tab-pane>
          <el-tab-pane name="global_config">
            <template #label>
              {{ $t('global_config') }}
            </template>
            <el-scrollbar max-height="calc(100vh - 276px)">
              <div class="tab_content">
                <ConfigGlobal :cluster="cluster" @refresh="$emit('refresh')"></ConfigGlobal>
              </div>
            </el-scrollbar>
          </el-tab-pane>
          <el-tab-pane name="k8s_cluster">
            <template #label>
              Kubernetes
            </template>
            <el-scrollbar max-height="calc(100vh - 276px)">
              <div class="tab_content">
                <ConfigK8sCluster :cluster="cluster" @switchTab="$emit('switchTab', $event)"></ConfigK8sCluster>
              </div>
            </el-scrollbar>
          </el-tab-pane>
          <el-tab-pane name="etcd">
            <template #label>
              ETCD
            </template>
            <el-scrollbar max-height="calc(100vh - 276px)">
              <div class="tab_content">
                <ConfigEtcd :cluster="cluster"></ConfigEtcd>
              </div>
            </el-scrollbar>
          </el-tab-pane>
          <el-tab-pane name="addons">
            <template #label>
              {{ $t('addons') }}
            </template>
            <el-scrollbar max-height="calc(100vh - 276px)">
              <div class="tab_content">
                <ConfigAddons :cluster="cluster" :currentTab="currentPropertiesTab" :pingpong="pingpong" @refresh="$emit('refresh')"></ConfigAddons>
              </div>
            </el-scrollbar>
          </el-tab-pane>
          <el-tab-pane :name="currentPropertiesTab" v-if="currentPropertiesTab.indexOf('NODE_') === 0 || currentPropertiesTab.indexOf('GAP_NODE_') === 0">
            <template #label>
              <div v-if="currentPropertiesTab.indexOf('NODE_') === 0" style="width: 100px; text-align: center;">{{ currentPropertiesTab.slice(5) }}</div>
              <div v-else style="width: 100px; text-align: center;">{{ currentPropertiesTab.slice(9) }}</div>
            </template>
            <el-scrollbar max-height="calc(100vh - 276px)">
              <div class="tab_content">
                <ConfigNode v-if="currentPropertiesTab.indexOf('NODE_') === 0" :cluster="cluster" :nodeName="currentPropertiesTab.slice(5)" 
                  :pingpong="pingpong" :pingpongLoading="pingpong_loading" @ping="ping"></ConfigNode>
                <CopyGapNodeToInventory v-else :cluster="nodeGap" :nodeName="currentPropertiesTab.slice(9)" 
                  :pingpong="pingpong" :pingpongLoading="pingpong_loading" @ping="ping"></CopyGapNodeToInventory>
              </div>
            </el-scrollbar>
          </el-tab-pane>
          <el-tab-pane v-else name="node_nodes">
            <template #label>
              <div style="width: 100px; text-align: center;">{{ $t('singleNode') }}</div>
            </template>
            <el-scrollbar max-height="calc(100vh - 276px)">
              <div class="tab_content">
                <el-alert type="warning" :closable="false" :title="$t('selectANode')">
                </el-alert>
              </div>
            </el-scrollbar>
          </el-tab-pane>
        </el-tabs>
      </el-form>
    </div>
  </div>
</template>

<script>
import { computed } from 'vue'
import Node from './Node.vue'
import ConfigKuboardSpray from './kuboardspray/ConfigKuboardSpray.vue'
import ConfigK8sCluster from './k8s_cluster/ConfigK8sCluster.vue'
import ConfigGlobal from './global/ConfigGlobal.vue'
import ConfigAddons from './addons/ConfigAddons.vue'
import ConfigNode from './node/ConfigNode.vue'
import CopyGapNodeToInventory from './node/CopyGapNodeToInventory.vue'
import ConfigEtcd from './etcd/ConfigEtcd.vue'
import AddNode from './common/AddNode.vue'

export default {
  props: {
    cluster: { type: Object, required: true },
    mode: { type: String, required: false, default: 'view' },
  },
  data () {
    return {
      currentPropertiesTab: 'global_config',
      pingpong: {},
      pingpong_loading: false,
    }
  },
  provide () {
    return {
      editMode: computed(() => {
        if (this.mode === 'view') {
          return 'view'
        }
        if (this.cluster && this.cluster.history.success_tasks.length > 0) {
          return 'frozen'
        }
        return this.mode
      }),
      currentTab: computed(() => {
        return this.currentPropertiesTab
      }),
    }
  },
  computed: {
    inventory: {
      get () {
        return this.cluster.inventory
      },
      set (v) {
        console.log(v)
      }
    },
    bastionEnabled() {
      return this.cluster.inventory.all.children.target.children.bastion !== undefined
    },
    nodeGap () {
      let temp = {
        name: this.cluster.name,
        type: 'gap',
        inventory: {
          all: {
            hosts: {},
            children: {
              target: {
                children: {
                  calico_rr: {
                    hosts: {}
                  },
                  etcd: {
                    hosts: {}
                  },
                  k8s_cluster: {
                    children: {
                      kube_control_plane: { hosts: {}},
                      kube_node: { hosts: {}},
                    }
                  }
                },
                vars: {},
              }
            }
          },
        }
      }
      if (this.cluster.state === undefined || this.cluster.state.nodes === undefined) {
        return temp
      }
      temp.inventory.all.children.target.vars = this.cluster.inventory.all.children.target.vars
      for (let nodeName in this.cluster.state.nodes) {
        let node = this.cluster.state.nodes[nodeName]
        if (this.cluster.inventory.all.hosts[nodeName] === undefined) {
          let ip = ''
          for (let i in node.status.addresses) {
            if (node.status.addresses[i].type === 'InternalIP') {
              ip = node.status.addresses[i].address
            }
          }
          temp.inventory.all.hosts[nodeName] = {
            ansible_host: ip
          }
        }
      }
      return temp
    }
  },
  components: { Node, ConfigKuboardSpray, ConfigK8sCluster, ConfigGlobal, ConfigAddons, ConfigNode, ConfigEtcd, AddNode, CopyGapNodeToInventory },
  mounted () {
    this.currentPropertiesTab = sessionStorage.getItem(this.cluster.name + '_plan_tab') || 'global_config'
  },
  watch: {
    currentPropertiesTab: function () {
      sessionStorage.setItem(this.cluster.name + '_plan_tab', this.currentPropertiesTab)
    },
    'inventory.all.hosts': function() {
      if (this.currentPropertiesTab.indexOf('NODE_') !== 0) {
        return
      }
      if (this.inventory.all.hosts[this.currentPropertiesTab.slice(5)] === undefined) {
        this.currentPropertiesTab = 'node_nodes'
      }
    }
  },
  methods: {
    validate (callback) {
      this.$refs.form.validate(callback)
    },
    nodeRoles (name) {
      let roles = {}
      for (let role in this.inventory.all.children.target.children.k8s_cluster.children) {
        for (let n in this.inventory.all.children.target.children.k8s_cluster.children[role].hosts) {
          if (n === name) {
            roles[role] = true
          }
        }
      }
      for (let n in this.inventory.all.children.target.children.etcd.hosts) {
        if (n === name) {
          roles.etcd = true
        }
      }
      return roles
    },
    isNode (name) {
      let roles = this.nodeRoles(name)
      let flag = false
      for (let key in roles) {
        if (key === 'kube_node') {
          flag = true
        } else {
          return false
        }
      }
      return flag
    },
    isEtcdAndNotControlPlane (name) {
      let roles = this.nodeRoles(name)
      return roles.etcd && !roles.kube_control_plane
    },
    showBastion () {
      this.currentPropertiesTab = 'global_config'; 
      setTimeout(() => {
        this.$refs.configKuboardSprayScroll.setScrollTop(0)
      }, 400)
    },
    ping () {
      this.pingpong = {}
      this.pingpong_loading = true
      let req = { nodes: 'target' }
      let count = 0
      for (let key in this.cluster.inventory.all.hosts) {
        if (key !== 'localhost' && key !== 'bastion') {
          count ++
        }
      }
      if (count == 0) {
        this.pingpong_loading = false
        return
      }
      this.kuboardSprayApi.post(`/clusters/${this.cluster.name}/state/ping`, req).then(resp => {
        this.pingpong = resp.data.data.items
        this.pingpong_loading = false
      }).catch(e => {
        if (e.response && e.response.data) {
          this.$message.error('不能测试节点是否在线: ' + e.response.data.message)
        } else {
          this.$message.error('不能测试节点是否在线: ' + e)
        }
        this.pingpong_loading = false
      })
    }
  }
}
</script>

<style lang="css">
.properties .el-tabs__item.is-active {
  padding-left: 0 20px !important;
  padding-right: 0 20px !important;
}
</style>

<style scoped lang="scss">
.plan {
  flex-grow: 2;
  display: flex;
  .left {
    width: 140px;
    min-height: 300px;
    border-radius: 10px;
    margin-left: 5px;
    .horizontalConnection {
      flex-grow: 1;
      width: calc(100% - 30px);
      height: 20px;
      border-bottom-left-radius: 5px;
      border-bottom: dashed 2px var(--el-color-primary);
      border-left: dashed 2px var(--el-color-primary);
      margin-left: 60px;
      margin-bottom: 23px;
      margin-right: -40px;
    }
    .verticalConnection {
      margin-left: 60px;
      border-left: dashed 2px var(--el-color-primary);
      height: 40px;
    }
  }
  .right {
    flex-grow: 1;
    min-height: 300px;
    border: solid 1.5px var(--el-color-primary-light-5);
    border-radius: 10px;
    margin-left: 30px;
    padding: 10px;
    .masters {
      display: flex;
      flex-wrap: wrap;
      padding-bottom: 10px;
      margin-bottom: 10px;
      border-bottom: dotted 1.5px var(--el-color-primary-light-5);
    }
    .workers {
      display: flex;
      flex-wrap: wrap;
    }
  }
}

.properties {
  flex-grow: 1;
  max-width: 45%;
  min-width: 45%;
  margin-left: 20px;
  .tab_content {
    padding: 10px;
    min-height: calc(100vh - 296px);
    background-color: #f1f4fa;
  }
}
</style>
