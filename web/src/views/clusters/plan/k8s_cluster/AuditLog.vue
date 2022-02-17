<i18n>
en:
  label: Audit Log
  description: Kubernetes audit log
zh:
  label: 审计日志
  description: Kubernetes 审计日志
</i18n>

<template>
  <ConfigSection v-model:enabled="enabled" :label="$t('label')" :description="$t('description')">
    <template #more>
      将 Kubernetes 审计信息存储到指定的路径。禁用后将触发 CIS Scan 错误
    </template>
    <FieldString :holder="vars" :prop="prop" fieldName="audit_log_path"></FieldString>
    <FieldNumber :holder="vars" :prop="prop" fieldName="audit_log_maxage"></FieldNumber>
    <FieldNumber :holder="vars" :prop="prop" fieldName="audit_log_maxbackups"></FieldNumber>
    <FieldNumber :holder="vars" :prop="prop" fieldName="audit_log_maxsize"></FieldNumber>
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
    prop () {
      return 'all.children.target.children.k8s_cluster.vars'
    },
    vars: {
      get () { return this.cluster.inventory.all.children.target.children.k8s_cluster.vars },
      set () {}
    },
    enabled: {
      get () {return this.vars.kubernetes_audit },
      set (v) {
        this.vars.kubernetes_audit = v
      },
    },
  },
  components: { },
  mounted () {
  },
  methods: {

  }
}
</script>

<style scoped lang="css">

</style>
