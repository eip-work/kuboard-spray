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
    <el-form-item :label="$t('field.kube_network_plugin')">
      <span class="app_text_mono">{{ kube_network_plugin }}</span>
    </el-form-item>
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
  }
}
</script>

<style scoped lang="scss">

</style>
