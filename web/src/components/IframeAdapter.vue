<template>
  <div></div>
</template>

<script>
export default {
  props: {
    command: { type: String, required: false, default: undefined },
    handler: { type: Function, required: false, default: undefined }
  },
  mounted () {
    if (this.command) {
      window.addEventListener('message', this.onParentMessage, false)
    }
  },
  unmounted () {
    if (this.command) {
      window.removeEventListener('message', this.onParentMessage)
    }
  },
  methods: {
    onParentMessage(event) {
      if (event.data.kuboardCommandToIframe === 'kuboard.cn:' + window.KuboardMfe.KUBOARD_FRAME_ID + ':' + this.command) {
        if (this.handler) {
          this.handler(event.data.params)
        } else {
          console.log('您未定义 handler', event)
        }
      }
    },
  }
}
</script>

<style>

</style>
