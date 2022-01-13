
<template>
  <el-button icon="el-icon-document-copy"
    type="primary" plain
    :disabled="value === undefined || value === ''" 
    v-clipboard:copy="value"
    v-clipboard:success="success"
    v-clipboard:error="failed"
  >复制到剪贴板</el-button>
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
      this.$message.success('已成功复制到剪贴板！')
    },
    failed () {
      this.$message.error('未能复制到剪贴板！')
    }
  }
}
</script>

<style scoped lang="scss">

</style>
