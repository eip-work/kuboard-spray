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
export default {
  props: {
    cluster: { type: Object, required: true },
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
  components: { },
  mounted () {
  },
  methods: {

  }
}
</script>

<style scoped lang="scss">

</style>
