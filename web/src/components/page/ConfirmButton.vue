<i18n>
en:
  continue: Continue to {text}
zh:
  continue: 确定（{text}）
</i18n>

<template>
  <el-popover v-model:visible="showConfirm" :placement="placement" :width="width" trigger="click">
    <template #reference>
      <el-button :style="buttonStyle" :type="type" :icon="icon" @click="showConfirm = true" :plain="plain">{{ text }}</el-button>
    </template>
    <slot>
      <el-alert :type="alertType[type]" effect="dark" :title="text" :icon="icon" :closable="false">
        <span>{{message}}</span>
      </el-alert>
      <div style="text-align: right; margin-top: 20px;">
        <el-button @click="showConfirm = false" type="info" plain icon="el-icon-close">{{ $t('msg.cancel') }}</el-button>
        <el-button @click="confirm" icon="el-icon-check" :type="type">{{ $t('continue', {text}) }}</el-button>
      </div>
    </slot>
  </el-popover>
</template>

<script>
export default {
  props: {
    placement: { type: String, required: false, default: 'bottom-end' },
    width: { type: String, required: false, default: '420px' },
    type: { type: String, required: false, default: 'primary' },
    icon: { type: String, required: false, default: '' },
    text: { type: String, required: false, default: 'text' },
    message: { type: String, required: false, default: 'message' },
    buttonStyle: { type: String, required: false, default: '' },
    plain: { type: Boolean, required: false, default: false },
  },
  data() {
    return {
      showConfirm: false
    }
  },
  emits: ['confirm'],
  computed: {
    alertType () {
      return {
        primary: 'warning',
        success: 'success',
        info: 'info',
        danger: 'error',
        warning: 'warning',
        default: 'default',
      }
    }
  },
  components: { },
  mounted () {
  },
  methods: {
    confirm () {
      this.showConfirm = false
      this.$emit('confirm')
    }
  }
}
</script>

<style scoped lang="scss">

</style>
