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
  operation: Operations
  install_cluster: Install / Setup K8S Cluster
  add_node: "{count} nodes waiting to be added"
  remove_node: "{count} nodes waiting to be removed"
  aboutToRemoveNode: About to remove the following nodes
  aboutToAddNode: About to add the following nodes
zh:
  verbose: 显示任务参数
  verbose_true: 日志中会包含部分敏感信息，例如：文件路径、用户名密码等
  verbose_false: 部分错误信息不能完整展示，使得出错时排查问题更困难
  vvv: 显示调试信息
  vvv_true: 日志中会包含最详细的信息
  vvv_false: 通常设置为 false
  fork: 并发数量
  fork_more: 安装过程中可以同时操作的目标节点的最大数量。ansible fork.
  operation: 操作类型
  install_cluster: 安装 / 设置集群
  add_node: "{count} 个节点待添加"
  remove_node: "{count} 个节点待删除"
  aboutToRemoveNode: 将要删除如下节点
  aboutToAddNode: 将要添加如下节点
</i18n>

<template>
  <ExecuteTask :history="cluster.history" :loading="loading" :title="title" :startTask="applyPlan" @refresh="$emit('refresh')" @visibleChange="updateForm">
    <el-form ref="form" :model="form" @submit.prevent.stop label-position="left" label-width="120px">
      <div style="height: 10px;"></div>
      <el-form-item :label="$t('verbose')">
        <el-switch v-model="form.verbose"></el-switch>
        <div style="width: 350px; font-size: 12px;">{{$t('verbose_' + form.verbose)}}</div>
      </el-form-item>
      <el-form-item :label="$t('vvv')">
        <el-switch v-model="form.vvv"></el-switch>
        <div style="width: 350px; font-size: 12px;">{{$t('vvv_' + form.vvv)}}</div>
      </el-form-item>
      <el-form-item :label="$t('fork')">
        <el-input-number v-model="form.fork" :step="2"></el-input-number>
        <div style="width: 350px; font-size: 12px;">{{$t('fork_more')}}</div>
      </el-form-item>
      <el-form-item v-if="pendingRemoveNodes.length > 0 || pendingAddNodes.length > 0" :label="$t('operation')" prop="action" required>
        <el-radio-group v-model="form.action">
          <div style="margin-bottom: 15px;">
            <el-radio v-if="pendingRemoveNodes.length >0" label="remove_node">
              {{ $t('aboutToRemoveNode') }}
              <div style="margin-left: 25px;">
                <el-tag v-for="(item, key) in pendingRemoveNodes" :key="'h' + key" style="margin-top: 10px; margin-right: 10px;">
                  <span class="app_text_mono">{{item.name}}</span>
                </el-tag>
              </div>
            </el-radio>
          </div>
          <div>
            <el-radio v-if="pendingAddNodes.length > 0" label="add_node">
              {{ $t('aboutToAddNode') }}
              <div style="margin-left: 25px;">
                <el-tag v-for="(item, key) in pendingAddNodes" :key="'h' + key" style="margin-top: 10px; margin-right: 10px;">
                  <span class="app_text_mono">{{item.name}}</span>
                </el-tag>
              </div>
            </el-radio>
          </div>
        </el-radio-group>
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
        action: '',
      }
    }
  },
  inject: ['isClusterInstalled', 'isClusterOnline', 'pendingRemoveNodes', 'pendingAddNodes'],
  computed: {
    action () {
      if (this.pendingRemoveNodes.length > 0 || this.pendingAddNodes.length > 0) {
        return this.form.action
      } else {
        return 'install_cluster'
      }
    },
    title () {
      return this.$t(this.action, {count: this.pendingAddNodes.length || this.pendingRemoveNodes.length })
    }
  },
  watch: {
  },
  components: { ExecuteTask },
  emits: ['refresh'],
  mounted () {
  },
  methods: {
    updateForm (flag) {
      if (flag) {
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
        if (this.pendingRemoveNodes.length > 0) {
          this.form.action = 'remove_node'
        } else if (this.pendingAddNodes.length > 0) {
          this.form.action = 'add_node'
        }
      }
    },
    async applyPlan () {
      return new Promise((resolve, reject) => {
        this.$refs.form.validate(flag => {
          if (flag) {
            this.kuboardSprayApi.post(`/clusters/${this.name}/${this.action}`, this.form).then(resp => {
              let pid = resp.data.data.pid
              resolve(pid)
            }).catch(e => {
              // this.$message.error('' + e.response.data.message)
              reject(e.response.data.message)
            })
          } else {
            reject('请验证表单')
          }
        })
      })
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
