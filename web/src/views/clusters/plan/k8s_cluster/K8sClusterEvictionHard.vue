<i18n>
en:
  label: eviction_hard
  description: eviction_hard params
  nolimit: No Limit
zh:
  label: eviction_hard
  description: eviction_hard 参数
  kube_api_anonymous_auth_and_insecure_port: 不能同时禁用匿名用户和非安全端口
  nolimit: 不限制
</i18n>

<template>
  <ConfigSection v-model:enabled="enabled" :label="$t('label')" :description="(fieldName == 'eviction_hard' ? '工作节点的 ' : '控制节点的 ') + $t('description')"
    helpLink="https://kubernetes.io/docs/concepts/scheduling-eviction/node-pressure-eviction/"
    anti-freeze labelWidth="120px">
    <FieldString :holder="vars[fieldName]['eviction.hard']" :prop="prop" placeholder="默认值: 100Mi" fieldName="memory.available"></FieldString>
    <FieldString :holder="vars[fieldName]['eviction.hard']" :prop="prop" placeholder="默认值: 10%" fieldName="nodefs.available"></FieldString>
    <FieldString :holder="vars[fieldName]['eviction.hard']" :prop="prop" placeholder="默认值: 5%" fieldName="nodefs.inodesFree"></FieldString>
    <FieldString :holder="vars[fieldName]['eviction.hard']" :prop="prop" placeholder="默认值: 15%" fieldName="imagefs.available"></FieldString>
  </ConfigSection>
</template>
    
<script>
export default {
  props: {
    cluster: { type: Object, required: true },
    fieldName: { type: String, required: true },
  },
  data() {
    return {
      
    }
  },
  computed: {
    prop () {
      return 'all.children.target.children.k8s_cluster.vars.' + this.fieldName + '.' + 'eviction.hard'
    },
    vars: {
      get () { return this.cluster.inventory.all.children.target.children.k8s_cluster.vars },
      set () {}
    },
    enabled: {
      get () {
        if (this.vars[this.fieldName]) {return true}
        return false
      },
      set (v) {
        if (v) {
          this.vars[this.fieldName] = {}
          this.vars[this.fieldName]['eviction.hard'] = {}
        } else {
          this.vars[this.fieldName] = undefined
        }
      },
    },
  },
  components: { },
  mounted () {
  },
  methods: {
    async loadCandidateAdmissionPlugins () {
      if (this.cluster.resourcePackage === undefined) {
        return []
      }
      let temp = this.cluster.resourcePackage.data.kubernetes.candidate_admission_plugins
      if (temp === undefined) {
        return []
      }
      temp = temp.split(',')
      let result = []
      for (let item of temp) {
        result.push({ label: item, value: item})
      }
      return result
    },
    async loadDisableAdmissionPlugins() {
      if (this.cluster.resourcePackage === undefined) {
        return []
      }
      let temp = this.cluster.resourcePackage.data.kubernetes.default_enabled_admission_plugins
      if (temp === undefined) {
        return []
      }
      temp = temp.split(',')
      let result = []
      for (let item of temp) {
        result.push({ label: item, value: item})
      }
      return result
    },
    async loadTlsCipherOptions () {
      return [
        { label: "TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA", value: "TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA" },
        { label: "TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256", value: "TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256" },
        { label: "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256", value: "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256" },
        { label: "TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA", value: "TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA" },
        { label: "TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384", value: "TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384" },
        { label: "TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305", value: "TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305" },
        { label: "TLS_ECDHE_ECDSA_WITH_RC4_128_SHA", value: "TLS_ECDHE_ECDSA_WITH_RC4_128_SHA" },
        { label: "TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA", value: "TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA" },
        { label: "TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA", value: "TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA" },
        { label: "TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256", value: "TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256" },
        { label: "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256", value: "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256" },
        { label: "TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA", value: "TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA" },
        { label: "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384", value: "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384" },
        { label: "TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305", value: "TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305" },
        { label: "TLS_ECDHE_RSA_WITH_RC4_128_SHA", value: "TLS_ECDHE_RSA_WITH_RC4_128_SHA" },
        { label: "TLS_RSA_WITH_3DES_EDE_CBC_SHA", value: "TLS_RSA_WITH_3DES_EDE_CBC_SHA" },
        { label: "TLS_RSA_WITH_AES_128_CBC_SHA", value: "TLS_RSA_WITH_AES_128_CBC_SHA" },
        { label: "TLS_RSA_WITH_AES_128_CBC_SHA256", value: "TLS_RSA_WITH_AES_128_CBC_SHA256" },
        { label: "TLS_RSA_WITH_AES_128_GCM_SHA256", value: "TLS_RSA_WITH_AES_128_GCM_SHA256" },
        { label: "TLS_RSA_WITH_AES_256_CBC_SHA", value: "TLS_RSA_WITH_AES_256_CBC_SHA" },
        { label: "TLS_RSA_WITH_AES_256_GCM_SHA384", value: "TLS_RSA_WITH_AES_256_GCM_SHA384" },
        { label: "TLS_RSA_WITH_RC4_128_SHA", value: "TLS_RSA_WITH_RC4_128_SHA" },
      ]
    }
  }
}
</script>

<style scoped lang="css">

</style>
