<i18n>
en:
  apply: Apply
  processingTitle: Executing
  processingHints: There is a background task executing, just wait a moment.
  succeeded: Installed history successfully before.
  viewLogs: View Task Logs
  taskFinished: Task is already completed
  closeWindow: Force to close this dialog window.
  taskInCurrent: Task in Running
  viewLastSuccessLog: Last success log
  viewLastLog: Last log

  install_cluster: Install / Setup K8S Cluster
  add_node: Add node to cluster
  remove_node: Remove node from cluster

zh:
  apply: 执 行
  processingTitle: 任务执行中
  processingHints: 当前有后台任务正在执行，请耐心等待。
  viewLogs: 查看任务日志
  taskFinished: 任务已结束
  closeWindow: 强制关闭此对话框

  reset: 将要执行任务
  succeeded: 已经成功完成
  confirmToExecute: 执行任务
  taskInCurrent: 当前有任务正在执行
  viewLastSuccessLog: 最后成功日志
  viewLastLog: 最后日志

  install_cluster: 安装 / 设置集群
  add_node: 添加节点到集群
  remove_node: 删除集群节点
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
    <el-popover v-if="!(finished && hideOnSuccess) && !history.processing" v-model:visible="showConfirm" :placement="placement" width="420" trigger="manual" :disabled="disabled">
      <template #reference>
        <el-button type="warning" icon="el-icon-lightning" @click="showConfirm = !showConfirm" @click.prevent.stop :disabled="disabled">{{ title || $t('apply')}}</el-button>
      </template>
      <el-form @submit.prevent.stop label-position="left" label-width="120px">
        <!-- <div style="height: 10px;"></div> -->
        <template v-if="finished">
          <!-- <el-alert type="success" effect="dark" style="margin-bottom: 10px;" :closable="false">
            <i class="el-icon-lightning" style="font-size: 15px; color: white; margin-right: 10px;"></i>
            <span class="confirmText" style="color: white;">{{$t('succeeded')}} [ {{ $t(lastSuccess.type) }} ]</span>
          </el-alert> -->
          <el-alert type="warning" effect="dark" style="margin-bottom: 10px;" :closable="false">
            <i class="el-icon-lightning" style="font-size: 15px; color: white; margin-right: 10px;"></i>
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
          <el-button v-if="history.success_tasks.length > 0" @click="viewTaskLogs(lastSucessPid)"
            type="success" plain icon="el-icon-files" style="float: left;">{{ $t('viewLastSuccessLog') }}</el-button>
          <el-button v-if="history.current_pid && history.current_pid !== lastSucessPid" @click="viewTaskLogs(history.current_pid)" 
            type="danger" plain icon="el-icon-document" style="float: left;">{{$t('viewLastLog')}}</el-button>
          <el-button type="default" icon="el-icon-close" @click="showConfirm = false">{{$t('msg.cancel')}}</el-button>
          <el-button type="primary" icon="el-icon-lightning" @click="applyPlan" :loading="executing">{{$t('msg.ok')}}</el-button>
        </div>
      </el-form>
    </el-popover>
    <el-button v-if="history.processing" type="danger" round @click="forceHide = false">{{ $t('taskInCurrent') }}</el-button>
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
    placement: { type: String, required: false, default: 'bottom-start' },
    disabled: { type: Boolean, required: false, default: false },
  },
  data() {
    return {
      forceHide: false,
      showConfirmData: false,
      executing: false,
    }
  },
  computed: {
    lastSuccess () {
      if (this.history.success_tasks.length > 0) {
        let temp = { successTime: '2000', type: 'install_cluster', pid: '' }
        for (let item of this.history.success_tasks) {
          if (temp.pid < item.pid) {
            temp = item
          }
        }
        return temp
      }
      return {}
    },
    lastSucessPid () {
      let task = undefined
      for (let t of this.history.success_tasks) {
        task = task || t
        if (task.pid < t.pid) {
          task = t
        }
      }
      return task ? task.pid : ''
    },
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
    },
    showConfirm: {
      get () { return this.showConfirmData },
      set (v) {
        this.$emit('visibleChange', v)
        this.showConfirmData = v
      }
    }
  },
  watch: {
    'loading': function() {
      this.forceHide = false
      this.showConfirm = false
    },
  },
  components: { },
  emits: ['refresh', 'visibleChange'],
  mounted () {
  },
  methods: {
    async applyPlan () {
      this.executing = true
      await this.startTask.call(this.$parent).then(pid => {
        if (pid) {
          this.forceHide = false
          this.showConfirm = false
          this.viewTaskLogs(pid)
          this.$emit('refresh')
        }
      }).catch(e => {
        let msg = e
        if (e && e.response && e.response.data && e.response.data.message) {
          msg = e.response.data.message
        }
        this.$message.error('failed to start task: ' + msg)
      })
      this.executing = false
    },
    viewTaskLogs(pid) {
      this.openUrlInBlank(`/#/tail/${this.history.task_type}/${this.history.task_name}/history/${pid}/execute.log`)
    }
  }
}
</script>

<style scoped lang="css">
.confirmText {
  font-size: 15px;
  font-weight: bold;
}
</style>
