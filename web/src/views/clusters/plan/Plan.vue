<i18n>
en:
  singleNode: Specific Node
  global_config: Global Config
  addons: Addons
  enabledBation: Enabled
  disabledBation: Disabled
  selectANode: Please select a node from the diagram to the left.
  bastionPrompt: Bastion params are placed at the bottom of KuboardSpray tab.
zh:
  singleNode: 单个节点
  global_config: 全局设置
  addons: 可选组件
  enabledBation: 使用堡垒机
  disabledBation: 不使用堡垒机
  selectANode: 请从左侧图中选择一个节点
  bastionPrompt: 堡垒机参数设置在 KuboardSpray 标签页的最后
</i18n>

<template>
  <div style="display:flex;">
    <div class="plan">
      <div class="left">
        <div style="padding: 5px; font-weight: bolder; font-size: 14px;">Kuboard Spray</div>
        <div>
          <Node class="localhost" name="localhost" :inventory="inventory" style="width: 100px;" hideDeleteButton
            :active="currentPropertiesTab === 'localhost'" @click="currentPropertiesTab = 'localhost'">
          </Node>
        </div>
        <div>
          <div v-if="bastionEnabled" class="verticalConnection"></div>
          <div v-else class="horizontalConnection"></div>
        </div>
        <div>
          <Node class="bastion" name="bastion" :inventory="inventory" style="width: 100px;" hideDeleteButton
            :active="currentPropertiesTab === 'bastion'" @click="showBastion">
            <div style="margin-top: 10px;">
              <el-tag v-if="bastionEnabled" type="danger" effect="dark" size="small" style="width: 100px; text-align: center;">{{$t('enabledBation')}}</el-tag>
              <el-tag v-else type="info" effect="light" size="small" style="width: 100px; text-align: center;">{{$t('disabledBation')}}</el-tag>
            </div>
          </Node>
          <div class="horizontalConnection" :style="bastionEnabled ? '' : 'border-color: white;'"></div>
        </div>
      </div>
      <div class="right">
        <div style="padding: 5px; font-weight: bolder; font-size: 14px; height: 28px; line-height: 28px;">
          <span style="margin-right: 20px;">Kubernetes Cluster</span>
          <AddNode :inventory="inventory" v-model:currentPropertiesTab="currentPropertiesTab"></AddNode>
        </div>
        <el-scrollbar height="calc(100vh - 283px)">
          <div class="masters">
            <Node v-for="(item, index) in inventory.all.children.target.children.k8s_cluster.children.kube_control_plane.hosts" :key="'control_plane' + index"
              @click="currentPropertiesTab = 'NODE_' + index" @delete_button="deleteNode(index)"
              :active="nodeRoles(index)[currentPropertiesTab] || currentPropertiesTab === 'global_config' || currentPropertiesTab === 'addons' || currentPropertiesTab === 'k8s_cluster' || 'NODE_' + index === currentPropertiesTab"
              :name="index" :inventory="inventory"></Node>
            <template v-for="(item, index) in inventory.all.children.target.children.etcd.hosts" :key="'etcd' + index">
              <Node v-if="isEtcdAndNotControlPlane(index)" :name="index" :inventory="inventory"
                @click="currentPropertiesTab = 'NODE_' + index" @delete_button="deleteNode(index)"
                :active="nodeRoles(index)[currentPropertiesTab] || currentPropertiesTab === 'global_config' || currentPropertiesTab === 'k8s_cluster' || 'NODE_' + index === currentPropertiesTab"></Node>
            </template>
          </div>
          <div class="workers">
            <template v-for="(item, index) in inventory.all.children.target.children.k8s_cluster.children.kube_node.hosts" :key="'node' + index">
              <Node v-if="isNode(index)" :name="index" :inventory="inventory"
                @click="currentPropertiesTab = 'NODE_' + index" @delete_button="deleteNode(index)"
                :active="nodeRoles(index)[currentPropertiesTab] || currentPropertiesTab === 'global_config' || currentPropertiesTab === 'addons' || currentPropertiesTab === 'k8s_cluster' || 'NODE_' + index === currentPropertiesTab"></Node>
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
              {{ $t('obj.localhost') }}
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
                <ConfigGlobal :cluster="cluster"></ConfigGlobal>
              </div>
            </el-scrollbar>
          </el-tab-pane>
          <el-tab-pane name="k8s_cluster">
            <template #label>
              Kubernetes
            </template>
            <el-scrollbar max-height="calc(100vh - 276px)">
              <div class="tab_content">
                <ConfigK8sCluster :cluster="cluster"></ConfigK8sCluster>
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
                <ConfigAddons :cluster="cluster"></ConfigAddons>
              </div>
            </el-scrollbar>
          </el-tab-pane>
          <el-tab-pane :name="currentPropertiesTab" v-if="currentPropertiesTab.indexOf('NODE_') === 0 && cluster.inventory.all.hosts[currentPropertiesTab.slice(5)]">
            <template #label>
              <div style="width: 100px; text-align: center;">{{ currentPropertiesTab.slice(5) }}</div>
            </template>
            <el-scrollbar max-height="calc(100vh - 276px)">
              <div class="tab_content">
                <ConfigNode :cluster="cluster" :nodeName="currentPropertiesTab.slice(5)"></ConfigNode>
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
import ConfigKuboardSpray from './kuboard_spray/ConfigKuboardSpray.vue'
import ConfigK8sCluster from './k8s_cluster/ConfigK8sCluster.vue'
import ConfigGlobal from './global/ConfigGlobal.vue'
import ConfigAddons from './addons/ConfigAddons.vue'
import ConfigNode from './node/ConfigNode.vue'
import ConfigEtcd from './etcd/ConfigEtcd.vue'
import AddNode from './common/AddNode.vue'

export default {
  props: {
    cluster: { type: Object, required: true },
    mode: { type: String, required: false, default: 'view' },
  },
  data () {
    return {
      currentPropertiesTab: 'k8s_cluster',
    }
  },
  provide () {
    return {
      editMode: computed(() => {
        if (this.mode === 'view') {
          return 'view'
        }
        if (this.cluster && this.cluster.success_tasks.length > 0) {
          return 'frozen'
        }
        return this.mode
      })
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
  },
  components: { Node, ConfigKuboardSpray, ConfigK8sCluster, ConfigGlobal, ConfigAddons, ConfigNode, ConfigEtcd, AddNode },
  mounted () {
  },
  watch: {
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
      this.currentPropertiesTab = 'localhost'; 
      // setTimeout(() => {
      //   this.$refs.configKuboardSprayScroll.setScrollTop(10)
      // }, 200)
      setTimeout(() => {
        this.$refs.configKuboardSprayScroll.setScrollTop(3000)
      }, 400)
      this.$message.info(this.$t('bastionPrompt'))
    },
    deleteNode(nodeName) {
      delete this.inventory.all.hosts[nodeName]
      delete this.inventory.all.children.target.children.k8s_cluster.children.kube_control_plane.hosts[nodeName]
      delete this.inventory.all.children.target.children.k8s_cluster.children.kube_node.hosts[nodeName]
      delete this.inventory.all.children.target.children.etcd.hosts[nodeName]
    }
  }
}
</script>

<style lang="scss">
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
    // border: solid 1.5px $--color-primary-light-5;
    border-radius: 10px;
    margin-left: 5px;
    // padding: 10px;
    .horizontalConnection {
      flex-grow: 1;
      width: calc(100% - 30px);
      height: 20px;
      border-bottom-left-radius: 5px;
      border-bottom: dashed 2px $--color-primary;
      border-left: dashed 2px $--color-primary;
      margin-left: 60px;
      margin-bottom: 23px;
      margin-right: -40px;
    }
    .verticalConnection {
      margin-left: 60px;
      border-left: dashed 2px $--color-primary;
      height: 40px;
    }
  }
  .right {
    flex-grow: 1;
    min-height: 300px;
    border: solid 1.5px $--color-primary-light-5;
    border-radius: 10px;
    margin-left: 30px;
    padding: 10px;
    .masters {
      display: flex;
      flex-wrap: wrap;
      padding-bottom: 10px;
      margin-bottom: 10px;
      border-bottom: dotted 1.5px $--color-primary-light-5;
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
