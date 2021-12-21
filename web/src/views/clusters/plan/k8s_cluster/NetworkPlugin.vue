<i18n>
en:
  label: NetworkPlugin
  // selectResourcePackageFirst: Please select resource package from KuboardSpray tab first.
zh:
  label: 网络插件
  // selectResourcePackageFirst: 请先在 KuboardSpray 标签页选择资源包
</i18n>


<template>
  <ConfigSection v-model:enabled="enabled" :label="$t('label')" disabled>
    <FieldSelect ref="kube_network_plugin" :disabled="!cluster.inventory.all.hosts.localhost.kuboardspray_resource_package"
      :holder="cluster.inventory.all.children.target.children.k8s_cluster.vars" fieldName="kube_network_plugin"
      :loadOptions="loadNetworkCandidates"></FieldSelect>
  </ConfigSection>
</template>

<script>
export default {
  props: {
    cluster: { type: Object, required: true },
    resourcePackage: { type: Object, required: false, default: undefined },
  },
  data() {
    return {

    }
  },
  computed: {
    enabled: {
      get () {
        return true
      },
      set () {}
    }
  },
  components: { },
  mounted () {
  },
  watch: {
    resourcePackage(newValue) {
      if (newValue !== undefined) {
        this.$refs.kube_network_plugin.load(true)
      }
    }
  },
  methods: {
    async loadNetworkCandidates () {
      let result = []
      if (this.resourcePackage !== undefined) {
        for (let item of this.resourcePackage.cni) {
          result.push({ value: item.name, label: `${item.name}_${item.version}`})
        }
      } else {
        // this.$message.error(this.$t('selectResourcePackageFirst'))
      }
      return result
    }
  }
}
</script>

<style scoped lang="scss">

</style>
