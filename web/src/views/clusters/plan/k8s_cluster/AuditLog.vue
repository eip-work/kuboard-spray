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
    <FieldString :holder="vars" :prop="prop" fieldName="audit_log_path" :required="enabled"></FieldString>
    <FieldNumber :holder="vars" :prop="prop" fieldName="audit_log_maxage" :required="enabled"></FieldNumber>
    <FieldNumber :holder="vars" :prop="prop" fieldName="audit_log_maxbackups" :required="enabled"></FieldNumber>
    <FieldNumber :holder="vars" :prop="prop" fieldName="audit_log_maxsize" :required="enabled"></FieldNumber>
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
        if (v) {
          // Populate safe defaults if missing
          if (this.vars.audit_log_path === undefined) {
            this.vars.audit_log_path = '/var/log/audit/kube-apiserver-audit.log'
          }
          if (this.vars.audit_log_maxage === undefined) {
            this.vars.audit_log_maxage = 30
          }
          if (this.vars.audit_log_maxbackups === undefined) {
            this.vars.audit_log_maxbackups = 10
          }
          if (this.vars.audit_log_maxsize === undefined) {
            this.vars.audit_log_maxsize = 100
          }
        } else {
          // Clear fields when disabled to avoid passing empty flags downstream
          delete this.vars.audit_log_path
          delete this.vars.audit_log_maxage
          delete this.vars.audit_log_maxbackups
          delete this.vars.audit_log_maxsize
        }
      },
    },
  },
  components: { },
  mounted () {
    // Ensure defaults are present on initial render when audit is enabled
    if (this.enabled) {
      if (this.vars.audit_log_path === undefined) {
        this.vars.audit_log_path = '/var/log/audit/kube-apiserver-audit.log'
      }
      if (this.vars.audit_log_maxage === undefined) {
        this.vars.audit_log_maxage = 30
      }
      if (this.vars.audit_log_maxbackups === undefined) {
        this.vars.audit_log_maxbackups = 10
      }
      if (this.vars.audit_log_maxsize === undefined) {
        this.vars.audit_log_maxsize = 100
      }
    }
  },
  methods: {

  }
}
</script>

<style scoped lang="css">

</style>
