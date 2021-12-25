<i18n>
en:
  copy: Copy to Clipboard
  copySuccess: Copied to Clipboard successfully!
  copyError: Failed to copy to Clipboard.
zh:
  copy: 复制到剪贴板
  copySuccess: 已成功复制到剪贴板！
  copyError: 未能复制到剪贴板！
</i18n>

<template>
  <el-button 
    icon="el-icon-copy-document"
    type="primary" plain
    :disabled="value === undefined || value === ''" 
    v-clipboard:copy="value"
    v-clipboard:success="success"
    v-clipboard:error="failed"
  >{{ $t('copy') }}</el-button>
</template>

<script>


export default {
  props: {
    value: { type: String, required: false, default: undefined }
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
    success () {
      this.$message.success(this.$t('copySuccess'))
    },
    failed () {
      this.$message.error(this.$t('copyError'))
    }
  }
}
</script>

<style scoped lang="scss">

</style>
