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
    <el-card v-if="kube_network == 'calico'" shadow="none" :body-style="{padding: '10px 20px 10px 20px', 'background-color': 'var(--el-color-info-light-9)'}" class="app_margin_bottom">
      <component :is="kube_network" :cluster="cluster"></component>
    </el-card>
  </ConfigSection>
</template>

<script>
import calico from './network_plugin/Calico.vue'

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
    },
    kube_network () {
      return this.cluster.inventory.all.children.target.children.k8s_cluster.vars.kube_network_plugin
    }
  },
  components: { calico },
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
