<i18n>
en:
  label: Common
  description: Most frequently referenced params to Kubernetes.
  nolimit: No Limit
zh:
  label: 通用参数
  description: Kubernetes 集群的通用参数
  kube_api_anonymous_auth_and_insecure_port: 不能同时禁用匿名用户和非安全端口
  nolimit: 不限制
</i18n>

<template>
  <ConfigSection v-model:enabled="enabled" :label="$t('label')" :description="$t('description')" disabled anti-freeze>
    <template v-if="cluster.resourcePackage !== undefined">
      <el-form-item label="kube_version">
        <span class="app_text_mono">{{cluster.resourcePackage.data.kubernetes.kube_version}}</span>
      </el-form-item>
    </template>
    <FieldString :holder="vars" :prop="prop" fieldName="cluster_name" required :rules="clusterNameRules"></FieldString>
    <FieldString :holder="vars" :prop="prop" fieldName="event_ttl_duration" required></FieldString>
    <FieldSelect v-if="cluster.resourcePackage && cluster.resourcePackage.data.kubernetes.candidate_admission_plugins" 
      :holder="vars" :prop="prop" fieldName="kube_apiserver_enable_admission_plugins" multiple anti-freeze
      :loadOptions="loadCandidateAdmissionPlugins">
      <template #view>
        <el-tag v-for="(item, index) in vars.kube_apiserver_enable_admission_plugins" :key="'adm' + index" style="margin: 5px;">
          {{ item }}
        </el-tag>
      </template>
    </FieldSelect>
    <FieldSelect v-if="cluster.resourcePackage && cluster.resourcePackage.data.kubernetes.default_enabled_admission_plugins" 
      :holder="vars" :prop="prop" fieldName="kube_apiserver_disable_admission_plugins" multiple anti-freeze
      :loadOptions="loadDisableAdmissionPlugins">
      <template #view>
        <el-tag v-for="(item, index) in vars.kube_apiserver_disable_admission_plugins" :key="'adm' + index" style="margin: 5px;">
          {{ item }}
        </el-tag>
      </template>
    </FieldSelect>
    <FieldRadio :holder="vars" :prop="prop" fieldName="kube_log_level" :options="[0, 1, 2, 3]" required></FieldRadio>
    <FieldNumber :holder="vars" :prop="prop" fieldName="kubelet_event_record_qps">
      <template #append>
        <span v-if="vars.kubelet_event_record_qps === 0">{{$t('nolimit')}}</span>
        <span v-else>个</span>
      </template>
      <template #edit_desc>
        <span v-if="vars.kubelet_event_record_qps < 10 && vars.kubelet_event_record_qps !== 0" style="color: var(--el-color-danger);">
          请设置为 0 或者大于等于 10 的数字
        </span>
      </template>
    </FieldNumber>
    <FieldSelect :holder="vars" :prop="prop" fieldName="tls_cipher_suites" :loadOptions="loadTlsCipherOptions" multiple>
      <template #view>
        <el-tag v-for="cipher in vars.tls_cipher_suites" :key="cipher" type="info">{{ cipher }}</el-tag>
      </template>
    </FieldSelect>
    <!-- <FieldBool :holder="vars" fieldName="auto_renew_certificates"></FieldBool> -->

    <!-- kubelet_rotate_server_certificates -->
  </ConfigSection>
</template>

<script>
export default {
  props: {
    cluster: { type: Object, required: true },
  },
  data() {
    return {
      insecurePortRules: [
        {
          validator: (rule, value, callback) => {
            // if (!this.vars.kube_api_anonymous_auth && (value === 0 || value === undefined)) {
            //   return callback(this.$t('kube_api_anonymous_auth_and_insecure_port'))
            // }
            return callback()
          },
          trigger: 'blur',
        }
      ],
      clusterNameRules: [
        {
          validator: (rule, value, callback) => {
            let reg = /^[a-z0-9]([-a-z0-9]*[a-z0-9])?(.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$/
            if (!reg.test(value)) {
              return callback("^[a-z0-9]([-a-z0-9]*[a-z0-9])?(.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$")
            }
            return callback()
          }
        }
      ]
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
      get () {return true},
      set () {},
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
