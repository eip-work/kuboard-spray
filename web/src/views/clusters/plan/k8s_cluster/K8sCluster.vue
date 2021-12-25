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
      <el-form-item :label="$t('container_manager')">
        <span class="app_text_mono">{{cluster.resourcePackage.container_engine.container_manager}}_{{cluster.resourcePackage.container_engine.version}}</span>
      </el-form-item>
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
  },
  components: { },
  mounted () {
  },
  methods: {

  }
}
</script>

<style scoped lang="scss">

</style>
