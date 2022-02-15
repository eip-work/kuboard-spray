<i18n>
en:
  label: Common
  description: Most frequently referenced params to Kubernetes.
zh:
  label: 通用参数
  description: Kubernetes 集群的通用参数
</i18n>

<template>
  <ConfigSection v-model:enabled="enabled" :label="$t('label')" :description="$t('description')" disabled anti-freeze>
    <template v-if="cluster.resourcePackage !== undefined">
      <el-form-item label="kube_version">
        <span class="app_text_mono">{{cluster.resourcePackage.data.kubernetes.kube_version}}</span>
      </el-form-item>
    </template>
    <FieldString :holder="vars" :prop="prop" fieldName="cluster_name" required></FieldString>
    <FieldString :holder="vars" :prop="prop" fieldName="event_ttl_duration" required></FieldString>
    <FieldSelect v-if="cluster.resourcePackage && cluster.resourcePackage.data.kubernetes.candidate_admission_plugins" 
      :holder="vars" :prop="prop" fieldName="kube_apiserver_enable_admission_plugins" multiple
      :loadOptions="loadCandidateAdmissionPlugins">
    </FieldSelect>
    <FieldSelect v-if="cluster.resourcePackage && cluster.resourcePackage.data.kubernetes.default_enabled_admission_plugins" :holder="vars" :prop="prop" fieldName="kube_apiserver_disable_admission_plugins" multiple
      :loadOptions="loadDisableAdmissionPlugins"></FieldSelect>
    <FieldBool :holder="vars" :prop="prop" fieldName="kube_api_anonymous_auth" required>
      <template #edit_desc>
        <span v-if="vars.kube_api_anonymous_auth">{{$t('cis.conflict_warn')}}</span>
      </template>
    </FieldBool>
    <FieldRadio :holder="vars" :prop="prop" fieldName="kube_log_level" :options="[0, 1, 2, 3]" required></FieldRadio>
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
    }
  }
}
</script>

<style scoped lang="css">

</style>
