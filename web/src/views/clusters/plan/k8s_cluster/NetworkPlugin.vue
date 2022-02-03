<i18n>
en:
  label: NetworkPlugin
zh:
  label: 网络插件
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
        let cnies = this.cluster.resourcePackage.data.network_plugin
        for (let i in cnies) {
          let cni = cnies[i]
          if (cni.name === this.cluster.inventory.all.children.target.children.k8s_cluster.vars.kube_network_plugin) {
            return cni.name + '_' + cni.params[cni.name + '_version']
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
        let cnies = resp.data.data.package.data.network_plugin
        for (let i in cnies) {
          let cni = cnies[i]
          result.push({
            label: cni.name + '_' + cni.params[cni.name + '_version'],
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

<style scoped lang="css">

</style>
