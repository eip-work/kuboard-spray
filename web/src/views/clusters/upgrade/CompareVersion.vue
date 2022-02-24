<i18n>
en:
  versionInResource: Version in Resource Package
  componentName: Component Name
  addon: Addons
  network_plugin: Cni plugins
  dependency: Dependencies
  etcd_cluster: etcd
  loading: Loading versions installed on nodes
  command_for_version: Command
  cordoned: Cordoned
zh:
  versionInResource: 资源包中的版本
  componentName: 组件名称
  addon: 可选组件
  network_plugin: 网络插件
  dependency: 依赖组件
  etcd_cluster: etcd
  loading: 查询节点上已安装的版本
  command_for_version: 查询命令
  cordoned: 已暂停调度
</i18n>


<template>
  <div>
    <el-table :data="tableData" style="width: 100%" row-key="name" :expand-row-keys="expanded" height="calc(100vh - 270px)">
      <el-table-column prop="name" :label="$t('componentName')" width="180" fixed>
        <template #header>
          <div class="compare_version_header">
            <div style="text-align: center;">{{ $t('componentName') }}</div>
          </div>
        </template>
        <template #default="scope">
          <div class="app_text_mono component_name nowrap">{{ $t(scope.row.name) }}</div>
        </template>
      </el-table-column>
      <el-table-column prop="version" :label="$t('versionInResource')" width="120" fixed>
        <template #header>
          <div class="compare_version_header">
            <div>{{ $t('versionInResource') }}</div>
          </div>
          </template>
        <template #default="scope">
          <span class="app_text_mono component_name">{{ scope.row.version }}</span>
        </template>
      </el-table-column>
      <template v-if="version">
        <el-table-column min-width="120px" fixed>
          <template #header>
            <div class="compare_version_header">{{ $t('command_for_version') }}</div>
          </template>
          <template #default="scope">
            <el-tooltip trigger="click" v-if="firstNode && firstNode[scope.row.name] && firstNode[scope.row.name].cmd" placement="top" width="420">
              <el-button icon="el-icon-finished">{{$t('command_for_version')}}</el-button>
              <template #content>
                <pre style="margin: 0 10px;">{{ firstNode[scope.row.name].cmd }}</pre>
              </template>
            </el-tooltip>
          </template>
        </el-table-column>
        <el-table-column v-for="(nodeVersion, nodeName) in version" :key="'col' + nodeName" min-width="120px">
          <template #header>
            <div class="compare_version_header">
              <UpgradeTask v-if="showUpgradeButton(nodeName)" :cluster="cluster" :nodeName="nodeName" :controlPlanePendingUpgrade="false" @refresh="$emit('refresh')"></UpgradeTask>
              <div v-else class="nowrap">{{ nodeName }}</div>
            </div>
          </template>
          <template #default="scope">
            <template v-if="scope.row.name === 'kubernetes'">
              <template v-if="cluster.state.nodes[nodeName]">
                <el-tag v-if="cluster.state.nodes[nodeName].spec.unschedulable" type="danger" effect="dark">{{ $t('cordoned') }}</el-tag>
              </template>
            </template>
            <template v-else-if="nodeVersion[scope.row.name]">
              <el-tag v-if="nodeVersion[scope.row.name].skipped" type="info">skipped</el-tag>
              <el-tooltip v-else-if="nodeVersion[scope.row.name].unreachable" trigger="hover" placement="top" width="420">
                <el-tag type="danger" effect="dark" style="cursor: pointer">unreachable</el-tag>
                <template #content>
                  <pre style="margin: 0 10px; width: 450px;">{{ nodeVersion[scope.row.name].msg }}</pre>
                </template>
              </el-tooltip>
              <el-tag v-else-if="nodeVersion[scope.row.name].stdout === scope.row.version" type="success">{{ nodeVersion[scope.row.name].stdout }}</el-tag>
              <el-tag v-else-if="nodeVersion[scope.row.name].stdout" 
                :type="scope.row.name === 'docker' && nodeVersion[scope.row.name].stdout.indexOf(scope.row.version) === 0 ? 'success' : 'danger'">
                {{ nodeVersion[scope.row.name].stdout }}
              </el-tag>
              <el-tooltip v-else-if="nodeVersion[scope.row.name].stderr" trigger="hover" placement="top" width="420">
                <el-tag type="danger" effect="dark" style="cursor: pointer">error</el-tag>
                <template #content>
                  <pre style="margin: 0 10px; width: 450px;">{{ nodeVersion[scope.row.name].stderr }}</pre>
                </template>
              </el-tooltip>
              <el-tag v-else type="info">empty</el-tag>
            </template>
          </template>
        </el-table-column>
      </template>
      <el-table-column v-if="version === undefined" min-width="180px">
        <template #header>
          <div class="compare_version_header">
            <el-icon style="vertical-align: middle;" class="is-loading">
              <el-icon-loading></el-icon-loading>
            </el-icon>
            <span style="margin-left: 10px;">{{$t('loading')}}</span>
          </div>
        </template>
        <template #default>
          <el-icon class="is-loading">
            <el-icon-loading></el-icon-loading>
          </el-icon>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script>
import mixin from './compare_mixin.js'
import UpgradeTask from './UpgradeTask.vue'

export default {
  props: {
    cluster: { type: Object, required: true },
    version: { type: Object, required: false, default: undefined },
    controlPlanePendingUpgrade: { type: Boolean, required: false, default: true },
  },
  data() {
    return {
      expanded: ['kubernetes', 'container_engine', 'etcd_cluster', 'network_plugin', 'dependency', 'addon']
    }
  },
  computed: {
    firstNode() {
      if (this.version) {
        for (let key in this.version) {
          return this.version[key]
        }
      }
      return undefined
    }
  },
  mixins: [mixin],
  components: { UpgradeTask },
  mounted () {
  },
  methods: {
    showUpgradeButton(nodeName) {
      if (this.cluster.history.processing) {
        return false
      }
      if (this.controlPlanePendingUpgrade) {
        return false
      }
      let host = this.cluster.inventory.all.hosts[nodeName]
      if (host.kuboardspray_node_action === 'upgrade_node') {
        return true
      }
      if (this.cluster.state.nodes[nodeName]) {
        if (this.cluster.state.nodes[nodeName].spec.unschedulable) {
          return true
        }
        return this.cluster.state.nodes[nodeName].status.nodeInfo.kubeletVersion !== this.cluster.resourcePackage.data.kubernetes.kube_version
      }
      return false
    }
  }
}
</script>

<style scoped lang="css">
.component_name {
  font-weight: bolder;
  display: inline-block;
  vertical-align: top;
  width: 120px;
}
.versionStr {
  vertical-align: top;
}
.compare_version_header {
  text-align: left;
}
</style>
