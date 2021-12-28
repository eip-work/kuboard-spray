<i18n>
en:
  label: NetworkPlugin
  // selectResourcePackageFirst: Please select resource package from KuboardSpray tab first.
zh:
  label: 网络插件
  // selectResourcePackageFirst: 请先在 KuboardSpray 标签页选择资源包
</i18n>


<template>
  <ConfigSection v-model:enabled="enabled" :label="$t('label')" disabled anti-freeze>
    <FieldSelect :holder="cluster.inventory.all.children.target.children.k8s_cluster.vars" fieldName="kube_network_plugin"
      prop="all.children.target.children.k8s_cluster.vars"
      required :loadOptions="loadKubeNetworkPlugin"></FieldSelect>
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
      get () {
        return true
      },
      set () {}
    },
    kube_network_plugin () {
      if (this.cluster.resourcePackage) {
        let cnies = this.cluster.resourcePackage.cni
        for (let i in cnies) {
          let cni = cnies[i]
          if (cni.name === this.cluster.inventory.all.children.target.children.k8s_cluster.vars.kube_network_plugin) {
            return cni.name + (cni.version ? '_' + cni.version : '')
          }
        }
      }
      return this.cluster.inventory.all.children.target.children.k8s_cluster.vars.kube_network_plugin
    }
  },
  components: { },
  mounted () {
  },
  methods: {
    async loadKubeNetworkPlugin () {
      let result = []
      await this.kuboardSprayApi.get(`/resources/${this.cluster.inventory.all.hosts.localhost.kuboardspray_resource_package}`).then(resp => {
        let cnies = resp.data.data.package.cni
        for (let i in cnies) {
          let cni = cnies[i]
          result.push({
            label: cni.name + (cni.version ? '_' + cni.version : ''),
            value: cni.name,
          })
        }
      }).catch(e => {
        console.log(e)
      })
      return result
    },
  }
}
</script>

<style scoped lang="scss">

</style>
