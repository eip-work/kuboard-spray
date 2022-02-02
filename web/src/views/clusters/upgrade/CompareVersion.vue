<i18n>
en:
  versionInResource: Version in Resource Package
  componentName: Component Name
  addon: Addons
  network_plugin: Cni plugins
  dependency: Dependencies
  etcd_cluster: etcd
  loading: Loading versions installed on nodes

zh:
  versionInResource: 资源包中的版本
  componentName: 组件名称
  addon: 可选组件
  network_plugin: 网络插件
  dependency: 依赖组件
  etcd_cluster: etcd
  loading: 查询节点上已安装的版本
</i18n>


<template>
  <div>
    <el-table :data="tableData" style="width: 100%" row-key="name" :expand-row-keys="expanded" height="calc(100vh - 270px)">
      <el-table-column prop="name" :label="$t('componentName')" width="180" fixed>
        <template #header>
          <div class="compare_version_header">
            <div style="line-height: 51px; text-align: center;">{{ $t('componentName') }}</div>
          </div>
        </template>
        <template #default="scope">
          <div class="app_text_mono component_name nowrap">{{ $t(scope.row.name) }}</div>
        </template>
      </el-table-column>
      <el-table-column prop="version" :label="$t('versionInResource')" width="120" fixed>
        <template #header>
          <div class="compare_version_header">
            <div style="line-height: 51px;">{{ $t('versionInResource') }}</div>
          </div>
          </template>
        <template #default="scope">
          <span class="app_text_mono component_name">{{ scope.row.version }}</span>
        </template>
      </el-table-column>
      <template v-if="version">
        <el-table-column v-for="(nodeVersion, nodeName) in version" :key="'col' + nodeName" min-width="120px">
          <template #header>
            <div class="compare_version_header">
              <div class="nowrap">{{ nodeName }}</div>
              <el-button>升级</el-button>
            </div>
          </template>
          <template #default="scope">
            <template v-if="nodeVersion[scope.row.name]">
              <el-tag v-if="nodeVersion[scope.row.name].stdout === scope.row.version" type="success">{{ nodeVersion[scope.row.name].stdout }}</el-tag>
              <el-tag v-else-if="nodeVersion[scope.row.name].stdout" type="error">{{ nodeVersion[scope.row.name].stdout }}</el-tag>
            </template>
          </template>
        </el-table-column>
      </template>
      <el-table-column v-if="version === undefined">
        <template #header>
          <div class="compare_version_header" style="line-height: 51px;">
            <i class="el-icon-loading"></i>
            <span style="margin-left: 10px;">{{$t('loading')}}</span>
          </div>
        </template>
        <template #default>
          <i class="el-icon-loading"></i>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script>
import mixin from './compare_mixin.js'

export default {
  props: {
    cluster: { type: Object, required: true },
    version: { type: Object, required: false, default: undefined },
  },
  data() {
    return {
      expanded: ['kubernetes', 'container_engine', 'etcd_cluster', 'network_plugin', 'dependency', 'addon']
    }
  },
  computed: {
  },
  mixins: [mixin],
  components: { },
  mounted () {
  },
  methods: {

  }
}
</script>

<style scoped lang="scss">
.component_name {
  font-weight: bolder;
  display: inline-block;
  vertical-align: top;
  width: 120px;
}
.versionStr {
  // white-space: pre-line !important;
  vertical-align: top;
}
.compare_version_header {
  height: 51px;
}
</style>
