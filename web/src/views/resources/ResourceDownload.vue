<i18n>
en:
  apply: Apply
  confirmToApply: Do the installation of Kubernetes Cluster, are you sure?
  processingTitle: Executing
  processingHints: There is a background task executing, just wait a moment.
  succeeded: Installed resource successfully before.
  viewLogs: View Task Logs
  taskFinished: Task is already completed
  viewPlan: I just want to view the Cluster Parameters.
zh:
  apply: 执 行
  confirmToApply: 将执行集群安装动作，请确认已完成集群规划！
  processingTitle: 任务执行中
  processingHints: 该集群有后台任务正在执行，请耐心等待。
  succeeded: 已经成功安装过 K8S 集群。
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
        <el-button type="primary" icon="el-icon-check" @click="viewTaskLogs(resource.current_pid)">{{$t('viewLogs')}}</el-button>
      </template>
    </el-dialog>
  </div>
  <el-popover v-else v-model:visible="showConfirm" placement="right-start" width="270" trigger="manual">
    <template #reference>
      <el-button type="danger" icon="el-icon-lightning" @click="showConfirm = !showConfirm">{{$t('apply')}}</el-button>
    </template>
    <el-form @submit.prevent.stop label-position="left" label-width="120px">
      <div style="height: 10px;"></div>
      <el-alert v-if="finished" type="success" effect="dark" style="margin-bottom: 10px;" :closable="false">
        <i class="el-icon-lightning" style="font-size: 16px; color: white; margin-right: 10px;"></i>
        <span class="confirmText" style="color: white;">{{$t('succeeded')}}</span>
      </el-alert>
      <el-alert type="error" style="margin-bottom: 10px;" :closable="false">
        <i class="el-icon-lightning" style="font-size: 16px; color: red; margin-right: 10px;"></i>
        <span class="confirmText">{{finished ? $t('reset') : $t('confirmToApply')}}</span>
      </el-alert>
      <div style="text-align: right; margin-top: 20px;">
        <el-button type="default" icon="el-icon-close" @click="showConfirm = false">{{$t('msg.cancel')}}</el-button>
        <el-button type="primary" icon="el-icon-lightning" @click="applyPlan">{{$t('msg.ok')}}</el-button>
      </div>
    </el-form>
  </el-popover>
</template>

<script>
import clone from 'clone'

export default {
  props: {
    resource: { type: Object, required: true },
    loading: { type: Boolean, required: false },
  },
  data() {
    return {
      forceHide: false,
      showConfirm: false,
    }
  },
  computed: {
    finished () {
      return false
    },
    dialogVisible: {
      get () {
        if (this.loading) {
          return false
        }
        return this.resource.processing && !this.forceHide
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
      this.forceHide = false
      this.showConfirm = false
      let request = clone(this.resource)
      request.downloadFrom = 'registry.cn-shanghai.aliyuncs.com/kuboard-spray/kuboard-spray-resource:spray-v2.17.1_k8s-v1.22.5_v1.0-amd64'
      this.kuboardSprayApi.post(`/resources/${this.resource.metadata.version}/download`, request).then(resp => {
        this.$emit('refresh')
        this.viewTaskLogs(resp.data.data.pid)
        this.$router.replace(`/settings/resources/${request.metadata.version}`)
      }).catch(e => {
        this.$message.error('' + e.response.data.message)
      })
    },
    viewTaskLogs (pid) {
      this.openUrlInBlank(`/#/tail/resource/${this.resource.metadata.version}/history/${pid}/execute.log`)
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
