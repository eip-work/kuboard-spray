<i18n>
en:
    "adustFontSize": "Adust Font Size"
    "fontSize": "Font Size"
zh:
    "adustFontSize": "调整字体大小"
    "fontSize": "字体大小"
</i18n>

<template>
  <el-popover
    placement="bottom"
    :title="$t('adustFontSize')"
    width="200"
    trigger="click">
    <template #reference>
      <el-button type="info">{{$t('fontSize')}}</el-button>
    </template>
    <el-input-number v-model="fontSize" style="width: 100%" @change="changeFontSize" :min="12" :max="20"></el-input-number>
  </el-popover>
</template>

<script>
export default {
  props: {
    terminal: { type: Object, required: false, default: undefined },
    fitAddon: { type: Object, required: false, default: undefined }
  },
  data () {
    return {
      fontSize: parseInt(localStorage.getItem('terminal-font-size')) || 14,
    }
  },
  methods: {
    changeFontSize () {
      this.terminal.setOption('fontSize', this.fontSize)
      localStorage.setItem('terminal-font-size', this.fontSize)
      this.$nextTick(() => {
        this.fitAddon.fit()
      })
    },
  }
}
</script>

<style>

</style>
