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
  <ConfigSection v-model:enabled="enabled" :label="$t('label')" :description="$t('obj.addon', {name: this.$t('description')})">
    <AddonFunction>{{$t('addon_function')}}</AddonFunction>
    <FieldString :holder="vars" fieldName="netcheck_namespace" :prop="prop"></FieldString>
    <FieldNumber :holder="vars" fieldName="netchecker_port" :prop="prop">
    </FieldNumber>
    <FieldNumber :holder="vars" fieldName="agent_report_interval" :prop="prop">
      <template #append>
        {{ $t('unit.second') }}
      </template>
    </FieldNumber>
  </ConfigSection>
</template>

<script>
import AddonFunction from './AddonFunction.vue'

export default {
  props: {
    cluster: { type: Object, required: true },
  },
  data() {
    return {

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
  components: { AddonFunction },
  mounted () {
  },
  methods: {

  }
}
</script>

<style scoped lang="scss">

</style>
