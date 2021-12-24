<i18n>
en:
  apply: Apply
  confirmToApply: Do the installation of OS Mirror, are you sure?
  succeeded: Successfully installed OS Mirror service before.
  reset: Continue to set new parameters to the OS Mirror service.
  processingTitle: Executing
  processingHints: There is a background task executing, just wait a moment.
  viewLogs: View Task Logs
  taskFinished: Task is already completed
  viewPlan: I just want to view the Parameters.
  verbose: Trace in details
  verbose_true: May include sensitive data in the trace, e.g. path to files, user name, password.
  verbose_false: Some information is hidden when there is a exception, which makes it more difficult to fix the issue.
  vvv: vvv
  vvv_true: includes more information in log
  vvv_false: usually false
zh:
  apply: 执 行
  confirmToApply: 将执行镜像服务的安装，请确认已完成参数配置！
  succeeded: 已经成功安装过镜像服务。
  reset: 此操作将重新设置镜像服务的参数！
  processingTitle: 任务执行中
  processingHints: 当前有后台任务正在执行，请耐心等待。
  viewLogs: 查看任务日志
  taskFinished: 任务已结束
  viewPlan: 我想看看参数
  verbose: 显示任务参数
  verbose_true: 日志中会包含部分敏感信息，例如：文件路径、用户名密码等
  verbose_false: 部分错误信息不能完整展示，使得出错时排查问题更困难
  vvv: 显示调试信息
  vvv_true: 日志中会包含最详细的信息
  vvv_false: 通常设置为 false
</i18n>

<template>
  <div v-if="dialogVisible" style="display: inline-block; width: 0px; height: 0px;">
    <el-dialog v-model="dialogVisible" :title="$t('processingTitle')" width="50%" :close-on-click-modal="false" :append-to-body="true" :show-close="false">
      <el-alert type="warning" :closable="false" :title="$t('processingTitle')" show-icon effect="dark" style="margin-top: 20px;">
        <span>{{$t('processingHints')}}</span>
      </el-alert>
      <template #footer>
        <el-button @click="forceHide = true" icon="el-icon-files">{{$t('viewPlan')}}</el-button>
        <el-button @click="$emit('refresh')" type="danger" icon="el-icon-finished">{{$t('taskFinished')}}</el-button>
        <el-button type="primary" icon="el-icon-check" @click="viewTaskLogs(os_mirror.current_pid)">{{$t('viewLogs')}}</el-button>
      </template>
    </el-dialog>
  </div>
  <el-popover v-else v-model:visible="showConfirm" placement="right-start" width="270" trigger="manual">
    <template #reference>
      <el-button type="danger" icon="el-icon-lightning" @click="showConfirm = !showConfirm">{{$t('apply')}}</el-button>
    </template>
    <el-form @submit.prevent.stop label-position="left" label-width="120px">
      <div style="height: 10px;"></div>
      <el-alert v-if="isInstalled" type="success" effect="dark" style="margin-bottom: 10px;" :closable="false">
        <i class="el-icon-lightning" style="font-size: 16px; color: white; margin-right: 10px;"></i>
        <span class="confirmText" style="color: white;">{{$t('succeeded')}}</span>
      </el-alert>
      <el-alert type="error" style="margin-bottom: 10px;" :closable="false">
        <i class="el-icon-lightning" style="font-size: 16px; color: red; margin-right: 10px;"></i>
        <span class="confirmText">{{isInstalled ? $t('reset') : $t('confirmToApply')}}</span>
      </el-alert>
      <el-form-item :label="$t('verbose')">
        <el-switch v-model="form.verbose"></el-switch>
        <div style="width: 240px; font-size: 12px;">{{$t('verbose_' + form.verbose)}}</div>
      </el-form-item>
      <el-form-item :label="$t('vvv')">
        <el-switch v-model="form.vvv"></el-switch>
        <div style="width: 240px; font-size: 12px;">{{$t('vvv_' + form.vvv)}}</div>
      </el-form-item>
      <div style="text-align: right; margin-top: 20px;">
        <el-button type="default" icon="el-icon-close" @click="showConfirm = false">{{$t('msg.cancel')}}</el-button>
        <el-button type="primary" icon="el-icon-lightning" @click="applyPlan">{{$t('msg.ok')}}</el-button>
      </div>
    </el-form>
  </el-popover>
</template>

<script>
export default {
  props: {
    os_mirror: { type: Object, required: true },
    name: { type: String, required: true },
    loading: { type: Boolean, required: false },
  },
  data() {
    return {
      forceHide: false,
      showConfirm: false,
      form: {
        verbose: false,
        vvv: false,
      }
    }
  },
  inject: ['isInstalled'],
  computed: {
    dialogVisible: {
      get () {
        if (this.loading) {
          return false
        }
        return this.os_mirror.processing && !this.forceHide
      },
      set () {}
    }
  },
  watch: {
    'loading': function() {
      this.forceHide = false
      this.showConfirm = false
    }
  },
  components: { },
  emits: ['refresh'],
  mounted () {
  },
  methods: {
    applyPlan () {
      this.forceHide = false
      this.showConfirm = false
      this.kuboardSprayApi.post(`/mirrors/${this.name}/install`, this.form).then(resp => {
        this.$emit('refresh')
        this.viewTaskLogs(resp.data.data.pid)
      }).catch(e => {
        this.$message.error('' + e.response.data.message)
      })
    },
    viewTaskLogs (pid) {
      this.openUrlInBlank(`/#/tail/mirror/${this.name}/history/${pid}/execute.log`)
    }
  }
}
</script>

<style scoped lang="scss">
.confirmText {
  font-size: 15px;
  color: $--color-danger;
  font-weight: bold;
}
</style>
