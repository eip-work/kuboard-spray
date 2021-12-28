<i18n>
en:
  label: Metrics
  description: metrics_server
  addon_function: Regularly takes metris info from K8S components, a must for features like Horizontal Auto Scaler.
zh:
  label: Metrics
  description: metrics_server
  addon_function: 定时采集 K8S 各组件的性能数据，是容器组水平伸缩等功能的基础依赖。
</i18n>

<template>
  <ConfigSection v-model:enabled="enabled" :label="$t('label')" :description="$t('obj.addon', {name: this.$t('description')})" anti-freeze>
    <AddonFunction>{{$t('addon_function')}}</AddonFunction>
    <FieldString :holder="vars" :prop="prop" fieldName="metrics_server_metric_resolution" :rules="resolutionRules"></FieldString>
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
      resolutionRules: [
        {
          validator: (rule, value, callback) => {
            if (value) {
              if (/^[0-9]+s/.test(value)) {
                return callback()
              } else {
                return callback('必须是数字，并加上 s 作为结尾')
              }
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
        return this.vars.metrics_server_enabled
      },
      set (v) {
        this.vars.metrics_server_enabled = v
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
