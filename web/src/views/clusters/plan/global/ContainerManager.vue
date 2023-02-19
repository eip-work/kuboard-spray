<i18n>
en:
  label: ContainerManager
  description: Container Manager
  insecure_registries_placeholder: e.g. 192.168.30.56:5000, use http protocol to pull from the registry
  registry_mirrors_placeholder: Registry Mirrors

zh:
  label: 容器引擎
  description: 设置集群范围内所有节点使用的容器引擎相关参数
  insecure_registries_placeholder: 例如：192.168.30.56:5000，使用 http 协议从该仓库拉取镜像
  registry_mirrors_placeholder: 镜像仓库 mirror，例如：https://registry.docker-cn.com
  
</i18n>

<template>
  <ConfigSection v-if="cluster && cluster.resourcePackage" v-model:enabled="enabled" :label="$t('label')" :description="$t('description')" disabled anti-freeze>
    <FieldSelect :holder="vars" fieldName="container_manager"
      :prop="prop" required :loadOptions="loadContainerEngines" :disabled="cluster.resourcePackage === undefined">
      <template #view_append>
        <span style="float: right;" v-if="isClusterOnline && cluster.resourcePackage.data.supported_playbooks.sync_container_engine_params">
          <ContainerMangerSyncParamsTask v-if="editMode === 'view'"
            :cluster="cluster" @refresh="$emit('refresh')"></ContainerMangerSyncParamsTask>
          <span v-else-if="editMode === 'frozen'">保存后点击 「更新容器引擎参数」 按钮</span>
        </span>
      </template>
    </FieldSelect>
    <template v-if="vars.container_manager === 'containerd'">
      <FieldString :holder="vars" :prop="prop" fieldName="containerd_cfg_dir"></FieldString>
      <FieldString :holder="vars" :prop="prop" fieldName="containerd_storage_dir"></FieldString>
      <FieldString :holder="vars" :prop="prop" fieldName="containerd_state_dir"></FieldString>
      <FieldBool :holder="vars" :prop="prop" fieldName="containerd_use_systemd_cgroup" disabled></FieldBool>
      <ContainerManagerRegistry v-if="compareVersions(cluster.resourcePackage.data.kubernetes.kube_version, 'v1.26.0') >= 0" :cluster="cluster"></ContainerManagerRegistry>
    </template>
    <template v-else>
      <el-alert class="app_alert_mini app_margin_bottom" type="warning">KuboardSpray 资源包中并不包含 docker 安装文件，必须在 “OS 软件源” 中设置 docker-ce 软件源</el-alert>
      <FieldBool :holder="vars" :prop="prop" fieldName="docker_orphan_clean_up"></FieldBool>
      <FieldArray :holder="vars" :prop="prop" newItemTemplate="" fieldName="docker_insecure_registries" :itemRules="insecureRegistriesItemRules" anti-freeze>
        <template v-slot:editItem="scope">
          <el-input v-model.trim="vars.docker_insecure_registries[scope.index]" :placeholder="$t('insecure_registries_placeholder')"></el-input>
        </template>
      </FieldArray>
      <div style="margin: -15px 0 5px 120px;">
        <KuboardSprayLink href="https://kuboard-spray.cn/guide/options/insecure-registry.html#docker" target="_blank" style="font-size: 12px;">帮 助</KuboardSprayLink>
      </div>
      <FieldArray :holder="vars" :prop="prop" newItemTemplate="" fieldName="docker_registry_mirrors" :itemRules="insecureRegistriesItemRules" anti-freeze>
        <template v-slot:editItem="scope">
          <el-input v-model.trim="vars.docker_registry_mirrors[scope.index]" :placeholder="$t('registry_mirrors_placeholder')"></el-input>
        </template>
      </FieldArray>
    </template>
  </ConfigSection>
</template>

<script>
import ContainerManagerRegistry from './ContainerManagerRegistry.vue'
import ContainerMangerSyncParamsTask from './ContainerMangerSyncParamsTask.vue'
import compareVersions from 'compare-versions'

export default {
  props: {
    cluster: { type: Object, required: true },
  },
  data() {
    return {
      insecureRegistriesItemRules: [
        { required: true, type: 'string', message: '字段不能为空', trigger: 'blur' },
      ],
    }
  },
  inject: ['isClusterOnline', 'editMode'],
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
  components: { ContainerManagerRegistry, ContainerMangerSyncParamsTask },
  mounted () {
  },
  methods: {
    compareVersions,
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

<style scoped lang="scss">
</style>
