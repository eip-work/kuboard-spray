<i18n>
en:
  label: nodelocaldns
  description: node local dns cache
  addon_function: Cache DNS result on node.
zh:
  label: 节点 DNS
  description: 节点 DNS 缓存
  addon_function: 在所有节点上缓存 DNS 的解析结果
</i18n>

<template>
  <AddonSection v-model:enabled="enabled" :label="$t('label')" :description="$t('obj.addon', {name: this.$t('description')})"
    :cluster="cluster" addonName="nodelocaldns" @refresh="$emit('refresh')">
    <template #more>{{$t('addon_function')}}</template>
    <FieldString :holder="vars" fieldName="nodelocaldns_ip" :prop="prop" :rules="nodelocaldns_ip_rules"></FieldString>
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
      nodelocaldns_ip_rules: [
        { 
          validator (rule, value, callback) {
            if (!/^(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])$/.test(value)) {
              return callback('requires an ip address')
            }
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
        return this.vars.enable_nodelocaldns
      },
      set (v) {
        this.vars.enable_nodelocaldns = v
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

<style scoped lang="css">

</style>
