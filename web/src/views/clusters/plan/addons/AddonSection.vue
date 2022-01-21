<i18n>
en:
  is_installed_true: Not Installed
  is_installed_false: Installed
zh:
  is_installed_true: 已安装
  is_installed_false: 未安装
</i18n>

<template>
  <ConfigSection v-model:enabled="enabledRef" :label="label" :helpLink="`https://kuboard-spray.cn/guide/addons/${addonName}.html`">
    <template #header>
      {{ label }}
      <div v-if="cluster && cluster.state && cluster.state.addons && cluster.state.addons[addonName]" style="float: right; margin-right: 10px;">
        <el-tag v-if="cluster.state.addons[addonName].is_installed" type="success">
          <i class="el-icon-check" style="margin-right: 5px;"></i>
          {{$t('is_installed_true')}}
        </el-tag>
        <el-tag v-else type="warning">
          <i class="el-icon-close" style="margin-right: 5px;"></i>
          {{$t('is_installed_false')}}
        </el-tag>
      </div>
    </template>
    <template #more>
      <slot name="more"></slot>
    </template>
    <slot></slot>
  </ConfigSection>
</template>

<script>
export default {
  props: {
    enabled: { prop: Boolean, required: true },
    label: { prop: String, require: false, default: "" },
    cluster: { prop: Object, required: true },
    addonName: { prop: String, required: true },
  },
  data() {
    return {

    }
  },
  computed: {
    enabledRef: {
      get () { return this.enabled },
      set (v) { this.$emit('update:enabled', v) }
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
