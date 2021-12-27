<i18n>
en:
  label: Common
  description: Most frequently referenced params to Kubernetes.
  container_manager: container_manager
zh:
  label: 常用参数
  description: Kubernetes 集群的常用参数
  container_manager: 容器引擎
</i18n>

<template>
  <ConfigSection v-model:enabled="enabled" :label="$t('label')" :description="$t('description')" disabled>
    <template v-if="cluster.resourcePackage !== undefined">
      <el-form-item label="kube_version">
        <span class="app_text_mono">{{cluster.resourcePackage.kubernetes.kube_version}}</span>
      </el-form-item>
      <FieldSelect :holder="cluster.inventory.all.children.target.children.k8s_cluster.vars" fieldName="container_manager"
        prop="all.children.target.children.k8s_cluster.vars" required :loadOptions="loadContainerEngines">
        <template #display_value>
          {{ container_manager }}
        </template>
      </FieldSelect>
    </template>
    <FieldString :holder="cluster.inventory.all.children.target.children.k8s_cluster.vars" fieldName="cluster_name"></FieldString>
    <FieldString :holder="cluster.inventory.all.children.target.children.k8s_cluster.vars" fieldName="event_ttl_duration"></FieldString>
    <!-- <FieldBool :holder="cluster.inventory.all.children.target.children.k8s_cluster.vars" fieldName="auto_renew_certificates"></FieldBool> -->
  </ConfigSection>
</template>

<script>
export default {
  props: {
    cluster: { type: Object, required: true },
  },
  data() {
    return {

    }
  },
  computed: {
    enabled: {
      get () {return true},
      set () {},
    },
    container_manager () {
      if (this.cluster.resourcePackage) {
        let engines = this.cluster.resourcePackage.container_engine
        for (let i in engines) {
          let engine = engines[i]
          if (engine.container_manager === this.cluster.inventory.all.children.target.children.k8s_cluster.vars.container_manager) {
            return engine.container_manager + (engine.version ? '_' + engine.version : '')
          }
        }
      }
      return this.cluster.inventory.all.children.target.children.k8s_cluster.vars.container_manager
    }
  },
  components: { },
  mounted () {
  },
  methods: {
    async loadContainerEngines () {
      let result = []
      let engines = this.cluster.resourcePackage.container_engine
      for (let i in engines) {
        let engine = engines[i]
        result.push({
          label: engine.container_manager + (engine.version ? '_' + engine.version : ''),
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
