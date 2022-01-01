<i18n>
en:
  apply: Apply
  processingTitle: Executing
  processingHints: There is a background task executing, just wait a moment.
  succeeded: Installed history successfully before.
  viewLogs: View Task Logs
  taskFinished: Task is already completed
  closeWindow: Force to close this dialog window.
zh:
  apply: 执 行
  processingTitle: 任务执行中
  processingHints: 当前有后台任务正在执行，请耐心等待。
  viewLogs: 查看任务日志
  taskFinished: 任务已结束
  closeWindow: 强制关闭此对话框

  reset: 再次尝试执行任务
  succeeded: 已经成功执行任务
  confirmToExecute: 执行任务
</i18n>

<template>
  <div v-if="dialogVisible" style="display: inline-block; width: 0px; height: 0px;">
    <el-dialog v-model="dialogVisible" :title="$t('processingTitle')" width="50%" :close-on-click-modal="false" :append-to-body="true" :show-close="false">
      <el-alert type="warning" :closable="false" :title="$t('processingTitle')" show-icon effect="dark" style="margin-top: 20px;">
        <span>{{$t('processingHints')}}</span>
      </el-alert>
      <template #footer>
        <el-button @click="forceHide = true" icon="el-icon-files">{{$t('closeWindow')}}</el-button>
        <el-button @click="$emit('refresh')" type="danger" icon="el-icon-finished">{{$t('taskFinished')}}</el-button>
        <el-button type="primary" icon="el-icon-check" @click="viewTaskLogs(history.current_pid)">{{$t('viewLogs')}}</el-button>
      </template>
    </el-dialog>
  </div>
  <template v-else-if="!loading">
    <el-popover v-if="!(finished && hideOnSuccess) && !history.processing" v-model:visible="showConfirm" placement="bottom" width="270" trigger="manual">
      <template #reference>
        <el-button type="warning" icon="el-icon-lightning" @click="showConfirm = !showConfirm">{{ label || $t('apply')}}</el-button>
      </template>
      <el-form @submit.prevent.stop label-position="left" label-width="120px">
        <div style="height: 10px;"></div>
        <template v-if="finished">
          <el-alert type="success" effect="dark" style="margin-bottom: 10px;" :closable="false">
            <i class="el-icon-lightning" style="font-size: 16px; color: white; margin-right: 10px;"></i>
            <span class="confirmText" style="color: white;">{{$t('succeeded')}} {{title ? `[ ${title} ]` : ''}}</span>
          </el-alert>
          <el-alert type="warning" style="margin-bottom: 10px;" :closable="false">
            <i class="el-icon-lightning" style="font-size: 16px; color: red; margin-right: 10px;"></i>
            <span class="confirmText">{{$t('reset')}} {{title ? `[ ${title} ]` : ''}}</span>
          </el-alert>
        </template>
        <template v-else>
          <el-alert type="error" style="margin-bottom: 10px;" :closable="false">
            <i class="el-icon-lightning" style="font-size: 16px; color: red; margin-right: 10px;"></i>
            <span class="confirmText">{{$t('confirmToExecute')}} {{title ? `[ ${title} ]` : ''}}</span>
          </el-alert>
        </template>

        <slot></slot>

        <div style="text-align: right; margin-top: 20px;">
          <el-button type="default" icon="el-icon-close" @click="showConfirm = false">{{$t('msg.cancel')}}</el-button>
          <el-button type="primary" icon="el-icon-lightning" @click="applyPlan">{{$t('msg.ok')}}</el-button>
        </div>
      </el-form>
    </el-popover>
  </template>
</template>

<script>
export default {
  props: {
    history: { type: Object, required: true },
    loading: { type: Boolean, required: false },
    startTask: { type: Function, required: true },
    label: { type: String, required: false, default: undefined },
    title: { type: String, required: false, default: undefined },
    hideOnSuccess: { type: Boolean, required: false, default: false },
  },
  data() {
    return {
      forceHide: false,
      showConfirm: false,
    }
  },
  computed: {
    finished () {
      if (this.history && this.history.success_tasks && this.history.success_tasks.length > 0) {
        return true
      }
      return false
    },
    dialogVisible: {
      get () {
        if (this.loading) {
          return false
        }
        return this.history.processing && !this.forceHide
      },
      set () {}
    }
  },
  watch: {
    'loading': function() {
      this.forceHide = false
      this.showConfirm = false
    },
  },
  components: { },
  emits: ['refresh'],
  mounted () {
  },
  methods: {
    applyPlan () {
      this.startTask.call(this.$parent).then(pid => {
        if (pid) {
          this.forceHide = false
          this.showConfirm = false
          this.viewTaskLogs(pid)
          this.$emit('refresh')
        }
      }).catch(e => {
        let msg = e
        if (e.response && e.response.data && e.response.data.message) {
          msg = e.response.data.message
        }
        this.$message.error('failed to start task: ' + msg)
      })
    },
    viewTaskLogs(pid) {
      this.openUrlInBlank(`/#/tail/${this.history.task_type}/${this.history.task_name}/history/${pid}/execute.log`)
    }
  }
}
</script>

<style scoped lang="scss">
.confirmText {
  font-size: 15px;
  // color: $--color-danger;
  font-weight: bold;
}
</style>
