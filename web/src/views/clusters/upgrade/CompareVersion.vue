<i18n>
en:
  versionInResource: Version in Resource Package
  componentName: Component Name
  addon: Addons
  network_plugin: Cni plugins
  dependency: Dependencies
  etcd_cluster: etcd

zh:
  versionInResource: 资源包中的版本
  componentName: 组件名称
  addon: 可选组件
  network_plugin: 网络插件
  dependency: 依赖组件
  etcd_cluster: etcd
</i18n>


<template>
  <div>
    <el-table :data="tableData" style="width: 100%" row-key="name" :expand-row-keys="expanded" height="calc(100vh - 270px)">
      <el-table-column fixed prop="name" :label="$t('componentName')" width="210">
        <template #default="scope">
          <div class="app_text_mono component_name nowrap">{{ $t(scope.row.name) }}</div>
        </template>
      </el-table-column>
      <el-table-column fixed prop="version" :label="$t('versionInResource')" width="180">
        <template #default="scope">
          <span class="app_text_mono component_name">{{ scope.row.version }}</span>
        </template>
      </el-table-column>
      <template v-if="version">
        <el-table-column v-for="(nodeVersion, nodeName) in version" :key="'col' + nodeName">
          <template #header>
            {{ nodeName }}
            <el-button>升级</el-button>
          </template>
          <template #default="scope">
            <div class="versionStr">{{ nodeVersion[scope.row.name] ? nodeVersion[scope.row.name].stdout : '' }}</div>
          </template>
        </el-table-column>
      </template>
    </el-table>
  </div>
</template>

<script>
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
    tableData () {
      let result = []
      let inv = this.cluster.inventory
      let res = this.cluster.resourcePackage.data
      let kube_version = res.kubernetes.kube_version
      let k8s = {
        name: 'kubernetes',
        version: kube_version,
        children: [
          { name: 'kubeadm', version: kube_version },
          { name: 'kubectl', version: kube_version },
          { name: 'kubelet', version: kube_version },
        ]
      }
      result.push(k8s)

      let containerEngine = { name: 'container_engine', version: undefined, children: [] }
      for (let ce of res.container_engine) {
        if (ce.container_manager == inv.all.children.target.vars.container_manager) {
          if (ce.container_manager === 'docker') {
            containerEngine.children.push({ name: 'docker', version: ce.params.docker_version })
            containerEngine.children.push({ name: 'containerd', version: ce.params.docker_containerd_version })
          } else {
            containerEngine.children.push({ name: 'containerd', version: ce.params.containerd_version })
          }
        }
      }
      result.push(containerEngine)

      result.push({ name: 'etcd_cluster', version: undefined, children: [
        { name: 'etcd', version: res.etcd.etcd_version }
      ] })

      let networkPlugin = { name: 'network_plugin', children: [] }
      for (let np of res.network_plugin) {
        if (np.name === inv.all.children.target.children.k8s_cluster.vars.kube_network_plugin) {
          for (let param in np.params) {
            if (param.indexOf('_version') > 0) {
              networkPlugin.children.push({ name: param.slice(0, param.length - 8), version: np.params[param]})
            }
          }
        }
      }
      result.push(networkPlugin)
      
      let dependency = { name: 'dependency', children: [] }
      for (let dep of res.dependency) {
        dependency.children.push({ name: dep.name, version: dep.version })
      }
      result.push(dependency)

      let addon = { name: 'addon', children: [] }
      for (let add of res.addon) {
        for (let key in add.params) {
          if (key.indexOf('_version') > 0) {
            let v = { name: add.name }
            v.version = add.params[key]
            addon.children.push(v)
          }
        }
      }
      result.push(addon)

      return result
    },
  },
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
  width: 150px;
}
.versionStr {
  white-space: pre-line !important;
  vertical-align: top;
}
</style>
