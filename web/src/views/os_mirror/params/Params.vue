<i18n>
en:
  params: kubespray params.
  params_desc: The following params will be used in kubespray if the os_mirror is selected.
zh:
  params: kubespray 参数
  params_desc: 如果被选中，以下参数将用于 kubespray。

</i18n>

<template>
  <div>
    <el-alert :title="$t('params')" :closable="false" class="app_margin_bottom">{{ $t('params_desc') }}</el-alert>
    <template v-for="(value, key) in os_mirror.status.params" :key="key">
      <FieldString :holder="os_mirror.status.params" :fieldName="key" label-width="240px"></FieldString>
    </template>
  </div>
</template>

<script>
export default {
  props: {
    os_mirror: { type: Object, required: true }
  },
  data() {
    return {

    }
  },
  computed: {
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
