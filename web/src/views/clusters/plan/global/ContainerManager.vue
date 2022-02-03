<i18n>
en:
  label: ContainerManager
  description: Container Manager
  insecure_registries_placeholder: e.g. 172.19.16.11:5000, use http protocol to pull from the registry
  registry_mirrors_placeholder: Registry Mirrors
zh:
  label: 容器引擎
  description: 设置集群范围内所有节点使用的容器引擎相关参数
  insecure_registries_placeholder: 例如：172.19.16.11:5000，使用 http 协议从该仓库拉取镜像
  registry_mirrors_placeholder: 镜像仓库 mirror，例如：https://registry.docker-cn.com
</i18n>

<template>
  <ConfigSection v-if="cluster && cluster.resourcePackage" v-model:enabled="enabled" :label="$t('label')" :description="$t('description')" disabled anti-freeze>
    <FieldSelect :holder="vars" fieldName="container_manager"
      :prop="prop" required :loadOptions="loadContainerEngines" :disabled="cluster.resourcePackage === undefined">
    </FieldSelect>
    <template v-if="vars.container_manager === 'containerd'">
      <FieldBool :holder="vars" :prop="prop" fieldName="containerd_use_systemd_cgroup" disabled></FieldBool>
      <FieldArray :holder="vars" :prop="prop" newItemTemplate="" fieldName="containerd_insecure_registries" :itemRules="insecureRegistriesItemRules" anti-freeze>
        <template v-slot:editItem="scope">
          <el-input v-model.trim="vars.containerd_insecure_registries[scope.index]" :placeholder="$t('insecure_registries_placeholder')"></el-input>
        </template>
      </FieldArray>
    </template>
    <template v-else>
      <el-alert class="app_alert_mini app_margin_bottom" type="warning">KuboardSpray 资源包中并不包含 docker 安装文件，必须在 “OS 软件源” 中设置 docker-ce 软件源</el-alert>
      <FieldBool :holder="vars" :prop="prop" fieldName="docker_orphan_clean_up"></FieldBool>
      <FieldArray :holder="vars" :prop="prop" newItemTemplate="" fieldName="docker_insecure_registries" :itemRules="insecureRegistriesItemRules" anti-freeze>
        <template v-slot:editItem="scope">
          <el-input v-model.trim="vars.docker_insecure_registries[scope.index]" :placeholder="$t('insecure_registries_placeholder')"></el-input>
        </template>
      </FieldArray>
      <FieldArray :holder="vars" :prop="prop" newItemTemplate="" fieldName="docker_registry_mirrors" :itemRules="insecureRegistriesItemRules" anti-freeze>
        <template v-slot:editItem="scope">
          <el-input v-model.trim="vars.docker_registry_mirrors[scope.index]" :placeholder="$t('registry_mirrors_placeholder')"></el-input>
        </template>
      </FieldArray>
    </template>
  </ConfigSection>
</template>

<script>
export default {
  props: {
    cluster: { type: Object, required: true },
  },
  data() {
    return {
      insecureRegistriesItemRules: [
        { required: true, type: 'string', message: 'cannot be empty', trigger: 'blur' },
      ]
    }
  },
  computed: {
    enabled: {
      get () {return true},
      set () {},
    },
    vars: {
      get () { return this.cluster.inventory.all.children.target.vars },
      set () {},
    },
    prop() { return 'all.children.target.vars' }
  },
  components: { },
  mounted () {
  },
  methods: {
    async loadContainerEngines () {
      let result = []
      let engines = this.cluster.resourcePackage.data.container_engine
      for (let i in engines) {
        let engine = engines[i]
        let label = engine.container_manager
        if (engine.params.containerd_version) {
          label += '_' + engine.params.containerd_version
        } else if (engine.params.docker_version) {
          label += '_' + engine.params.docker_version
        }
        result.push({
          label,
          value: engine.container_manager,
        })
      }
      return result
    },
  }
}
</script>

<style scoped lang="css">

</style>
