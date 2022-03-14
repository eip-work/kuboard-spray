<i18n>
en:
  error: "Cannot access the following url via websocket protocol: "
  "commonProblems": "Refer to help on kuboard.cn"
  "confirm": "Do you confirm to retry ?"
zh:
  error: 不能通过 websocket 协议访问如下地址：
  "commonProblems": "查看问题原因及解决办法"
  "confirm": "确认重试？"
</i18n>

<template>
  <div style="text-align: left;">
    <el-dialog :title="$t('msg.prompt')" v-model="dialogVisible" width="60%" :before-close="handleClose">
      <p>{{ $t('error') }}</p>
      <p class="app_text_mono" style="font-size: 12px; margin-left: 20px; margin-bottom: 20px;">{{message}}</p>
      <KuboardSprayLink href="https://kuboard-spray.cn/guide/extra/websocket.html" :size="14" type="danger">{{$t('commonProblems')}}</KuboardSprayLink>
    </el-dialog>
  </div>
</template>

<script>
export default {
  data() {
    return {
      dialogVisible: false,
      message: undefined
    }
  },
  methods: {
    show(message) {
      this.message = message
      this.dialogVisible = true
    },
    handleClose() {
      this.$confirm(this.$t('confirm'))
        .then(() => {
          window.location.reload()
        })
        .catch(() => { });
    }
  }
}
</script>

<style>
</style>
