<i18n>
en:
  singleNode: Specific Node
  global_config: Global Config
  enabledBation: Enabled
  disabledBation: Disabled
  selectANode: Please select a node from the diagram to the left.
  bastionPrompt: Bastion params are placed at the bottom of KuboardSpray tab.
zh:
  singleNode: 单个节点
  global_config: 全局设置
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
          <Node class="localhost" name="localhost" :inventory="inventory" style="width: 100px;"
            :active="currentPropertiesTab === 'localhost' || currentPropertiesTab === 'global_config'" @click="currentPropertiesTab = 'localhost'">
            <!-- <div class="role">{{$t('agent')}}</div> -->
          </Node>
        </div>
        <div>
          <div v-if="bastionEnabled" class="verticalConnection"></div>
          <div v-else class="horizontalConnection"></div>
        </div>
        <div>
          <Node class="bastion" name="bastion" :inventory="inventory" style="width: 100px;"
            :active="currentPropertiesTab === 'bastion'" @click="showBastion">
            <div style="margin-top: 10px;">
              <el-tag v-if="bastionEnabled" type="danger" effect="dark" size="small" style="width: 100px; text-align: center;">{{$t('enabledBation')}}</el-tag>
              <el-tag v-else type="info" effect="light" size="small" style="width: 100px; text-align: center;">{{$t('disabledBation')}}</el-tag>
            </div>
          </Node>
          <div class="horizontalConnection" :style="bastionEnabled ? '' : 'border-color: white;'"></div>
        </div>
        <!-- <el-button @click="$refs.form.validate()">校验</el-button>
        <el-radio-group v-model="mode">
          <el-radio-button label="create"></el-radio-button>
          <el-radio-button label="edit"></el-radio-button>
          <el-radio-button label="view"></el-radio-button>
        </el-radio-group> -->
      </div>
      <div class="right">
        <el-scrollbar height="calc(100vh - 240px)">
          <div class="masters">
            <Node v-for="(item, index) in inventory.all.children.k8s_cluster.children.kube_control_plane.hosts" :key="'control_plane' + index"
              @click="currentPropertiesTab = 'NODE_' + index"
              :active="nodeRoles(index)[currentPropertiesTab] || currentPropertiesTab === 'global_config' || currentPropertiesTab === 'k8s_cluster' || 'NODE_' + index === currentPropertiesTab"
              :name="index" :inventory="inventory"></Node>
            <template v-for="(item, index) in inventory.all.children.k8s_cluster.children.etcd.hosts" :key="'etcd' + index">
              <Node v-if="isEtcdAndNotControlPlane(index)" :name="index" :inventory="inventory"
                @click="currentPropertiesTab = 'NODE_' + index"
                :active="nodeRoles(index)[currentPropertiesTab] || currentPropertiesTab === 'global_config' || currentPropertiesTab === 'k8s_cluster' || 'NODE_' + index === currentPropertiesTab"></Node>
            </template>
          </div>
          <div class="workers">
            <template v-for="(item, index) in inventory.all.children.k8s_cluster.children.kube_node.hosts" :key="'node' + index">
              <Node v-if="isNode(index)" :name="index" :inventory="inventory"
                @click="currentPropertiesTab = 'NODE_' + index"
                :active="nodeRoles(index)[currentPropertiesTab] || currentPropertiesTab === 'global_config' || currentPropertiesTab === 'k8s_cluster' || 'NODE_' + index === currentPropertiesTab"></Node>
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
            <el-scrollbar max-height="calc(100vh - 274px)" ref="configKuboardSprayScroll">
              <div class="tab_content">
                <ConfigKuboardSpray v-model:resourcePackage="resourcePackage" :inventory="inventory" :cluster="cluster"></ConfigKuboardSpray>
              </div>
            </el-scrollbar>
          </el-tab-pane>
          <!-- <el-tab-pane name="global_config">
            <template #label>
              {{ $t('global_config') }}
            </template>
            <el-scrollbar max-height="calc(100vh - 274px)">
              <div class="tab_content">
                {{ currentPropertiesTab }}
              </div>
            </el-scrollbar>
          </el-tab-pane> -->
          <el-tab-pane name="k8s_cluster">
            <template #label>
              {{ $t('node.k8s_cluster') }}
            </template>
            <el-scrollbar max-height="calc(100vh - 274px)">
              <div class="tab_content">
                <ConfigK8sCluster :resourcePackage="resourcePackage" :inventory="inventory" :clusterName="cluster.name"></ConfigK8sCluster>
              </div>
            </el-scrollbar>
          </el-tab-pane>
          <!-- <el-tab-pane name="kube_control_plane">
            <template #label>
              {{ $t('node.kube_control_plane') }}
            </template>
            <el-scrollbar max-height="calc(100vh - 274px)">
              <div class="tab_content">
                {{ currentPropertiesTab }}
              </div>
            </el-scrollbar>
          </el-tab-pane> -->
          <el-tab-pane name="etcd">
            <template #label>
              {{ $t('node.etcd') }}
            </template>
            <el-scrollbar max-height="calc(100vh - 274px)">
              <div class="tab_content">
                <ConfigEtcd :resourcePackage="resourcePackage" :inventory="inventory" :clusterName="cluster.name"></ConfigEtcd>
              </div>
            </el-scrollbar>
          </el-tab-pane>
          <!-- <el-tab-pane name="kube_node">
            <template #label>
              {{ $t('node.kube_node') }}
            </template>
            <el-scrollbar max-height="calc(100vh - 274px)">
              <div class="tab_content">
                {{ currentPropertiesTab }}
              </div>
            </el-scrollbar>
          </el-tab-pane> -->
          <el-tab-pane :name="currentPropertiesTab" v-if="currentPropertiesTab.startsWith('NODE_')">
            <template #label>
              <div style="width: 100px; text-align: center;">{{ currentPropertiesTab.slice(5) }}</div>
            </template>
            <el-scrollbar max-height="calc(100vh - 274px)">
              <div class="tab_content">
                <ConfigNode :resourcePackage="resourcePackage" :inventory="inventory" :clusterName="cluster.name" :nodeName="currentPropertiesTab.slice(5)"></ConfigNode>
              </div>
            </el-scrollbar>
          </el-tab-pane>
          <el-tab-pane v-else name="node_nodes">
            <template #label>
              <div style="width: 100px; text-align: center;">{{ $t('singleNode') }}</div>
            </template>
            <el-scrollbar max-height="calc(100vh - 274px)">
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
import ConfigKuboardSpray from './ConfigKuboardSpray.vue'
import ConfigK8sCluster from './ConfigK8sCluster.vue'
import ConfigNode from './ConfigNode.vue'
import ConfigEtcd from './ConfigEtcd.vue'

export default {
  props: {
    cluster: { type: Object, required: true }
  },
  data () {
    return {
      currentPropertiesTab: 'k8s_cluster',
      mode: 'edit',
      resourcePackage: undefined,
    }
  },
  provide () {
    return {
      editMode: computed(() => this.mode)
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
      return this.cluster.inventory.all.children.bastion !== undefined
    },
  },
  components: { Node, ConfigKuboardSpray, ConfigK8sCluster, ConfigNode, ConfigEtcd },
  mounted () {
  },
  methods: {
    nodeRoles (name) {
      let roles = {}
      for (let role in this.inventory.all.children.k8s_cluster.children) {
        for (let n in this.inventory.all.children.k8s_cluster.children[role].hosts) {
          if (n === name) {
            roles[role] = true
          }
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
    min-height: calc(100vh - 294px);
    background-color: #f1f4fa;
  }
}
</style>
