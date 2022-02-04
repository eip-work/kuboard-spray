<i18n>
en:
  continue: Continue to {text}
zh:
  continue: 确定（{text}）
</i18n>

<template>
  <el-popover v-model:visible="showConfirm" :placement="placement" :width="width" trigger="mannual">
    <template #reference>
      <el-button v-if="buttonText === ''" :style="buttonStyle" :type="type" :circle="circle" :icon="icon" @click="showConfirm = true" :plain="plain"></el-button>
      <el-button v-else :style="buttonStyle" :type="type" :circle="circle" :icon="icon" @click="showConfirm = true" :plain="plain">{{ buttonText || text }}</el-button>
    </template>
    <el-alert :type="alertType[type]" effect="light" :title="text" :icon="icon" :closable="false">
      <span>{{message}}</span>
    </el-alert>
    <slot></slot>
    <div style="text-align: right; margin-top: 20px;">
      <el-button @click="showConfirm = false" type="info" plain icon="el-icon-close">{{ $t('msg.cancel') }}</el-button>
      <el-button @click="confirm" icon="el-icon-check" :type="type">{{ $t('continue', {text}) }}</el-button>
    </div>
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
    circle: { type: Boolean, required: false, default: false },
    buttonText: { type: String, required: false, default: undefined },
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

<style scoped lang="css">

</style>
