<i18n>
en:
  verbose: Include task params
  verbose_true: May include sensitive data in the trace, e.g. path to files, user name, password.
  verbose_false: Some information is hidden when there is a exception, which makes it more difficult to fix the issue.
  vvv: vvv
  vvv_true: includes more information in log
  vvv_false: usually false
  fork: ansible fork
  fork_more: Max number of nodes can be operated in the installation.

  installCluster: Install Cluster
zh:
  verbose: 显示任务参数
  verbose_true: 日志中会包含部分敏感信息，例如：文件路径、用户名密码等
  verbose_false: 部分错误信息不能完整展示，使得出错时排查问题更困难
  vvv: 显示调试信息
  vvv_true: 日志中会包含最详细的信息
  vvv_false: 通常设置为 false
  fork: 并发数量
  fork_more: 安装过程中可以同时操作的目标节点的最大数量。ansible fork.

  installCluster: 安装集群
</i18n>

<template>
  <ExecuteTask :history="cluster.history" :loading="loading" :title="$t('installCluster')" :startTask="applyPlan" @refresh="$emit('refresh')">
    <el-form @submit.prevent.stop label-position="left" label-width="120px">
      <div style="height: 10px;"></div>
      <el-form-item :label="$t('verbose')">
        <el-switch v-model="form.verbose"></el-switch>
        <div style="width: 240px; font-size: 12px;">{{$t('verbose_' + form.verbose)}}</div>
      </el-form-item>
      <el-form-item :label="$t('vvv')">
        <el-switch v-model="form.vvv"></el-switch>
        <div style="width: 240px; font-size: 12px;">{{$t('vvv_' + form.vvv)}}</div>
      </el-form-item>
      <el-form-item :label="$t('fork')">
        <el-input-number v-model="form.fork" :step="2"></el-input-number>
        <div style="width: 240px; font-size: 12px;">{{$t('fork_more')}}</div>
      </el-form-item>
    </el-form>
  </ExecuteTask>
</template>

<script>
import ExecuteTask from '../common/task/ExecuteTask.vue'

export default {
  props: {
    cluster: { type: Object, required: true },
    name: { type: String, required: true },
    loading: { type: Boolean, required: false },
  },
  data() {
    return {
      form: {
        verbose: false,
        vvv: false,
        fork: 5,
      }
    }
  },
  watch: {
    'cluster.inventory.all.hosts': function (newValue) {
      if (newValue) {
        let count = 0
        for (let key in this.cluster.inventory.all.hosts) {
          if (key !== 'localhost' && key !== 'bastion') {
            count ++
          }
        }
        if (count > 50) {
          count = 50
        }
        this.form.fork = count
      }
    }
  },
  components: { ExecuteTask },
  emits: ['refresh'],
  mounted () {
  },
  methods: {
    async applyPlan () {
      let pid = undefined
      await this.kuboardSprayApi.post(`/clusters/${this.name}/install`, this.form).then(resp => {
        pid = resp.data.data.pid
      }).catch(e => {
        this.$message.error('' + e.response.data.message)
      })
      return pid
    },
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
