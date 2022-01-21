<i18n>
en:
  label: net_checker
  description: net_checker
  addon_function: Deploy some agents in K8S to verify if overlay network functions well.
zh:
  label: 网络检查
  description: net_checker
  addon_function: 在 K8S 中部署几个测试的 Agent，用于验证容器之间的网络是否互通。
</i18n>

<template>
  <AddonSection v-model:enabled="enabled" :label="$t('label')" :description="$t('obj.addon', {name: this.$t('description')})" anti-freeze help-link="https://kuboard-spray.cn/guide/validate/netcheck.html"
    :cluster="cluster" addonName="netchecker">
    <template #more>{{$t('addon_function')}}</template>
    <FieldString :holder="vars" fieldName="netcheck_namespace" :prop="prop" :rules="namespaceRules"></FieldString>
    <FieldNumber :holder="vars" fieldName="netchecker_port" :prop="prop">
    </FieldNumber>
    <FieldNumber :holder="vars" fieldName="agent_report_interval" :prop="prop">
      <template #append>
        {{ $t('unit.second') }}
      </template>
    </FieldNumber>
  </AddonSection>
</template>

<script>
import AddonSection from '../AddonSection.vue'

export default {
  props: {
    cluster: { type: Object, required: true },
  },
  data() {
    return {
      namespaceRules: [
        { 
          validator (rule, value, callback) {
            let reg = new RegExp("^[a-z0-9]([-a-z0-9]*[a-z0-9])?$")
            if (!reg.test(value)) {
              return callback('[a-z0-9]([-a-z0-9]*[a-z0-9])?')
            }
            //a lowercase RFC 1123 label must consist of lower case alphanumeric characters or '-', and must start and end with an alphanumeric character (e.g. 'my-name',  or '123-abc', regex used for validation is '[a-z0-9]([-a-z0-9]*[a-z0-9])?')
            return callback()
          },
          trigger: 'blur'
        }
      ]
    }
  },
  computed: {
    vars: {
      get () { return this.cluster.inventory.all.children.target.children.k8s_cluster.vars },
      set () {},
    },
    prop () {
      return 'all.children.target.children.k8s_cluster.vars'
    },
    enabled: {
      get () {
        return this.vars.deploy_netchecker
      },
      set (v) {
        this.vars.deploy_netchecker = v
      }
    }
  },
  components: { AddonSection },
  mounted () {
  },
  methods: {

  }
}
</script>

<style scoped lang="scss">

</style>
