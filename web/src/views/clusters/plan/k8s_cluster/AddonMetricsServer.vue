<i18n>
en:
  hi: kuboard
zh:
  hi: kuboard
</i18n>

<template>
  <ConfigSection v-model:enabled="enabledMetricsServer" label="Metrics" :description="$t('obj.addon', {name: 'metrics_server'})">
    Hello
  </ConfigSection>
</template>

<script>
import ConfigSection from '../ConfigSection.vue'

export default {
  props: {
    inventory: { type: Object, required: true },
    resourcePackage: { type: Object, required: false, default: undefined },
  },
  data() {
    return {

    }
  },
  computed: {
    enabledMetricsServer: {
      get () {
        return false
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

  }
}
</script>

<style scoped lang="scss">

</style>
