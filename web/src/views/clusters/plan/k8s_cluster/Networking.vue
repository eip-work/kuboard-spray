<i18n>
en:
  label: Network
  description: Network parameters to kubernetes.
zh:
  label: 网络参数
  description: Kubernetes 集群的网络参数
</i18n>

<template>
  <ConfigSection v-model:enabled="enabled" :label="$t('label')" :description="$t('description')" disabled anti-freeze>
    <template #more>
      待优化，需要弄清楚各个属性之间的数量关系
    </template>
    <FieldString :holder="vars" :prop="prop" fieldName="kube_service_addresses" required></FieldString>
    <FieldString :holder="vars" :prop="prop" fieldName="kube_pods_subnet" required></FieldString>
    <FieldNumber :holder="vars" :prop="prop" fieldName="kube_network_node_prefix" required></FieldNumber>
    <FieldNumber :holder="vars" :prop="prop" fieldName="kubelet_max_pods"></FieldNumber>
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
