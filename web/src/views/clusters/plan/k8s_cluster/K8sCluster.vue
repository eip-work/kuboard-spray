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
        <span class="app_text_mono">{{cluster.resourcePackage.kubernetes.kube_version}}</span>
      </el-form-item>
    </template>
    <FieldString :holder="vars" :prop="prop" fieldName="cluster_name" required></FieldString>
    <FieldString :holder="vars" :prop="prop" fieldName="event_ttl_duration" required></FieldString>
    <FieldBool :holder="vars" :prop="prop" fieldName="kube_api_anonymous_auth" required></FieldBool>
    <FieldRadio :holder="vars" :prop="prop" fieldName="kube_log_level" :options="[0, 1, 2, 3]" required></FieldRadio>
    <!-- <FieldBool :holder="vars" fieldName="auto_renew_certificates"></FieldBool> -->
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
  }
}
</script>

<style scoped lang="scss">

</style>
