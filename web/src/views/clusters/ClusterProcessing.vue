<i18n>
en:
  apply: Apply
  confirmToApply: Do the installation of Kubernetes Cluster, are you sure?
  processingTitle: Executing
  processingHints: There is a background task executing, just wait a moment.
  viewLogs: View Task Logs
  taskFinished: Task is already completed
  viewPlan: I just want to view the Cluster Parameters.
zh:
  apply: 执 行
  confirmToApply: 将执行集群安装动作，请确认已完成集群规划！
  processingTitle: 任务执行中
  processingHints: 该集群有后台任务正在执行，请耐心等待。
  viewLogs: 查看任务日志
  taskFinished: 任务已结束
  viewPlan: 我想看看集群规划
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
  <el-popconfirm v-else-if="!isInstalled" :confirm-button-text="$t('msg.ok')" :cancel-button-text="$t('msg.cancel')" placement="bottom-start"
    icon="el-icon-lightning" icon-color="red" :title="$t('confirmToApply')" @confirm="applyPlan">
    <template #reference>
      <el-button type="danger" icon="el-icon-lightning">{{$t('apply')}}</el-button>
    </template>
  </el-popconfirm>
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
      forceHide: false
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
    }
  },
  components: { },
  emits: ['refresh'],
  mounted () {
  },
  methods: {
    applyPlan () {
      this.kuboardSprayApi.post(`/clusters/${this.name}/install`).then(resp => {
        this.$emit('refresh')
        this.viewTaskLogs(resp.data.data.pid)
      }).catch(e => {
        this.$message.error('' + e.response.data.message)
      })
    },
    viewTaskLogs (pid) {
      this.openUrlInBlank(`/#/clusters/${this.name}/history/${pid}/tail/command.log`)
    }
  }
}
</script>

<style scoped lang="scss">

</style>
