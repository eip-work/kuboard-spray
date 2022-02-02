<i18n>
en:
  versionInResource: Version in installed Resource Package
  componentName: Component Name
  addon: Addons
  network_plugin: Cni plugins
  dependency: Dependencies
  etcd_cluster: etcd
  target_version: Version in target Resource Package

zh:
  versionInResource: 已安装资源包中的版本
  componentName: 组件名称
  addon: 可选组件
  network_plugin: 网络插件
  dependency: 依赖组件
  etcd_cluster: etcd
  target_version: 目标资源包中的版本
</i18n>

<template>
  <div>
    <el-table :data="tableData" style="width: 100%" row-key="name" :expand-row-keys="expanded" height="calc(90vh - 245px)">
      <el-table-column prop="name" :label="$t('componentName')" width="280">
        <template #header>
          <div class="compare_version_header">
            <div style="text-align: center;">{{ $t('componentName') }}</div>
          </div>
        </template>
        <template #default="scope">
          <div class="app_text_mono component_name nowrap">{{ $t(scope.row.name) }}</div>
        </template>
      </el-table-column>
      <el-table-column prop="version" :label="$t('versionInResource')" width="280">
        <template #header>
          <div class="compare_version_header">
            <div>
              <i class="el-icon-circle-check"></i>
              {{ $t('versionInResource') }}
            </div>
          </div>
          </template>
        <template #default="scope">
          <span class="app_text_mono component_name">{{ scope.row.version }}</span>
        </template>
      </el-table-column>
      <el-table-column>
        <template #header>
          <div class="compare_version_header">
            <i class="el-icon-cloudy"></i>
            <span style="margin-left: 10px;">{{$t('target_version')}}</span>
          </div>
        </template>
        <template #default="scope">
          <el-tag v-if="targetComponents[scope.row.name] && scope.row.version" :type="tagType(targetComponents[scope.row.name], scope.row.version)" :effect="tagType(targetComponents[scope.row.name], scope.row.version) === 'error' ? 'dark' : 'light'">
            {{ targetComponents[scope.row.name] }}
          </el-tag>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script>
import mixin from './compare_mixin.js'
import compareVersions from 'compare-versions'

export default {
  props: {
    cluster: { type: Object, required: true },
    target: { type: Object, required: true },
  },
  data() {
    return {
      version: undefined,
      expanded: ['kubernetes', 'container_engine', 'etcd_cluster', 'network_plugin', 'dependency', 'addon']
    }
  },
  computed: {
    targetComponents () {
      let result = {}
      let res = this.target.data
      let kube_version = res.kubernetes.kube_version
      let k8s = {
          'kube-apiserver': kube_version,
          kubectl: kube_version,
          kubelet: kube_version,
          kubeadm: kube_version,
      }
      result = Object.assign(result, k8s)

      for (let ce of res.container_engine) {
        if (ce.container_manager === 'docker') {
          result.docker = ce.params.docker_version
          result.containerd = ce.params.docker_containerd_version
        } else {
          result.containerd = ce.params.containerd_version
        }
      }

      result.etcd = res.etcd.etcd_version

      for (let np of res.network_plugin) {
        for (let param in np.params) {
          if (param.indexOf('_version') > 0) {
            result[param.slice(0, param.length - 8)] = np.params[param]
          }
        }
      }
      
      for (let dep of res.dependency) {
        result[dep.name] = dep.version
      }

      for (let add of res.addon) {
        for (let key in add.params) {
          if (key.indexOf('_version') > 0) {
            result[add.name] = add.params[key]
          }
        }
      }

      return result
    }
  },
  mixins: [mixin],
  components: { },
  mounted () {
  },
  methods: {
    tagType(targetVersion, installedVersion) {
      if (installedVersion && targetVersion) {
        let temp = compareVersions(targetVersion + '', installedVersion + '')
        if (temp > 0) {
          return 'success'
        }
        if (temp === 0) {
          return 'primary'
        }
        if (temp < 0) {
          return 'error'
        }
      }
      return undefined
    }
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
  text-align: left;
}
</style>
