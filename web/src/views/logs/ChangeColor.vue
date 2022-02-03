<i18n>
en:
  changeColor: "Change Color"
  reason: "Change foreground color."
  confirmToReload: "This operation requires reloading the page, do you confirm to continue ?"
zh:
  changeColor: "修改前景色"
  reason: "修改前景色，以适应不同的光线环境"
  confirmToReload: "此操作需要将要重新加载页面，是否继续？"
</i18n>

<template>
  <el-popover
    placement="bottom"
    :title="$t('changeColor')"
    width="300"
    trigger="click">
    <template #reference>
      <el-button type="info">{{$t('changeColor')}}</el-button>
    </template>
    <span style="font-size: 13px; vertical-align: top; line-height: 28px; margin-right: 20px;">{{$t('reason')}}</span>
    <el-color-picker
      v-model="color"
      show-alpha @change="change"
      :predefine="predefineColors">
    </el-color-picker>
  </el-popover>
</template>

<script>
export default {
  props: {
  },
  data() {
    return {
      color: localStorage.getItem('terminalForegroundColor') || '#eee',
      predefineColors: [
        '#fff',
        '#eee',
        '#ddd',
        '#ccc',
        '#bbb',
      ]
    }
  },
  computed: {
  },
  components: { },
  mounted () {
  },
  methods: {
    change () {
      this.$confirm(this.$t('confirmToReload'), this.$t('message.prompt'), {
        type: 'warning'
      }).then(() => {
        localStorage.setItem('terminalForegroundColor', this.color)
        location.reload()
      }).catch(() => {
        this.$message({
          type: 'info',
          message: this.$t('message.canceled')
        });          
      });
    }
  }
}
</script>

<style scoped lang="css">

</style>
