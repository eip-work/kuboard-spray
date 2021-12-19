<i18n>
en:
  apply: Apply
  confirmToApply: Do the installation of Kubernetes Cluster, are you sure?
  processingTitle: Executing
  processingHints: There is a background task executing, just wait a moment.
  viewLogs: View Task Logs
  taskFinished: Task is already completed
  viewPlan: I just want to view the Cluster Parameters.
  verbose: Trace in details
  verbose_true: May include sensitive data in the trace, e.g. path to files, user name, password.
  verbose_false: Some information is hidden when there is a exception, which makes it more difficult to fix the issue.
zh:
  apply: 执 行
  confirmToApply: 将执行集群安装动作，请确认已完成集群规划！
  processingTitle: 任务执行中
  processingHints: 该集群有后台任务正在执行，请耐心等待。
  viewLogs: 查看任务日志
  taskFinished: 任务已结束
  viewPlan: 我想看看集群规划
  verbose: 显示详尽的日志信息
  verbose_true: 日志中会包含部分敏感信息，例如：文件路径、用户名密码等
  verbose_false: 部分错误信息不能完整展示，使得出错时排查问题更困难
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
        <el-button type="primary" icon="el-icon-check" @click="viewTaskLogs(cluster.current_pid)">{{$t('viewLogs')}}</el-button>
      </template>
    </el-dialog>
  </div>
  <el-popover v-else-if="!isInstalled" v-model:visible="showConfirm" placement="right-start" width="270" trigger="manual">
    <template #reference>
      <el-button type="danger" icon="el-icon-lightning" @click="showConfirm = !showConfirm">{{$t('apply')}}</el-button>
    </template>
    <el-form @submit.prevent.stop>
      <div style="height: 10px;"></div>
      <el-alert type="error" style="margin-bottom: 10px;" :closable="false">
        <i class="el-icon-lightning" style="font-size: 16px; color: red; margin-right: 10px;"></i>
        <span class="confirmText">{{$t('confirmToApply')}}</span>
      </el-alert>
      <el-form-item :label="$t('verbose')">
        <el-switch v-model="form.verbose"></el-switch>
        <div style="width: 240px; font-size: 12px;">{{$t('verbose_' + form.verbose)}}</div>
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
    cluster: { type: Object, required: true },
    name: { type: String, required: true },
    loading: { type: Boolean, required: false },
  },
  data() {
    return {
      forceHide: false,
      showConfirm: false,
      form: {
        verbose: false,
      }
    }
  },
  computed: {
    dialogVisible: {
      get () {
        if (this.loading) {
          return false
        }
        return this.cluster.processing && !this.forceHide
      },
      set () {}
    }
  },
  inject: ['isInstalled'],
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
      this.kuboardSprayApi.post(`/clusters/${this.name}/install`, this.form).then(resp => {
        this.$emit('refresh')
        this.viewTaskLogs(resp.data.data.pid)
      }).catch(e => {
        this.$message.error('' + e.response.data.message)
      })
    },
    viewTaskLogs (pid) {
      this.openUrlInBlank(`/#/clusters/${this.name}/history/${pid}/tail/execute.log`)
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
