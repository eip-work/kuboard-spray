<i18n>
en:
  ETCD: ETCD Params (Apply to all the etcd nodes)
zh:
  ETCD: ETCD 参数（适用范围：所有 etcd 节点）
</i18n>

<template>
  <div>
    <ConfigSection v-model:enabled="enabledEtcd" :description="$t('ETCD')" label="ETCD" disabled>
      <FieldSelect :holder="cluster.inventory.all.children.etcd.vars" prop="all.children.etcd.vars" required
        fieldName="etcd_deployment_type" :loadOptions="loadEtcdType"></FieldSelect>
    </ConfigSection>
  </div>
</template>

<script>
import ConfigSection from './ConfigSection.vue'

export default {
  props: {
    cluster: { type: Object, required: true },
  },
  data() {
    return {

    }
  },
  computed: {
    enabledEtcd: {
      get () {
        return true
      },
      set (v) {
        console.log(v)
      }
    }
  },
  components: { ConfigSection },
  mounted () {
  },
  methods: {
    async loadEtcdType () {
      return [
        { label: '以二进制形式部署', value: 'host' }, 
        { label: '以 Docker 形式部署', value: 'docker' }
      ]
    }
  }
}
</script>

<style scoped lang="scss">

</style>
