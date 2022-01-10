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
  control_params: Control params
  operation: Operations
  install_cluster: Install / Setup K8S Cluster
  add_node: "{count} nodes waiting to be added"
  remove_node: "{count} nodes waiting to be removed"
  aboutToRemoveNode: About to remove the following nodes
  aboutToAddNode: About to add the following nodes
  resetNodes: Reset nodes(non-applicable to offline nodes)
  resetNodesNo: Do not reset nodes
  allow_ungraceful_removal: Allow remove nodes failed to drain
  allow_ungraceful_removal_desc: 'kuboard-spray does "kubectl drain" before removing nodes, continue on failure?'
  drain_timeout: drain_timeout
  drain_retries: drain_retries
  drain_grace_period: drain_grace_period
  drain_retry_delay_seconds: drain_retry_delay_seconds
  nodes_to_remove_required: Please select nodes to remove.

  newResourcePackageRequired: This resource package does not support the operation, please choose a new one.
zh:
  verbose: 显示任务参数
  verbose_true: 日志中会包含部分敏感信息，例如：文件路径、用户名密码等
  verbose_false: 部分错误信息不能完整展示，使得出错时排查问题更困难
  vvv: 显示调试信息
  vvv_true: 日志中会包含最详细的信息
  vvv_false: 通常设置为 false
  fork: 并发数量
  fork_more: 安装过程中可以同时操作的目标节点的最大数量。ansible fork.
  control_params: 控制选项
  operation: 操作选项
  install_cluster: 安装 / 设置集群
  add_node: "{count} 个节点待添加"
  remove_node: "{count} 个节点待删除"
  aboutToRemoveNode: 将要删除如下节点
  aboutToAddNode: 将要添加如下节点
  resetNodes: 重置节点（不能选择已停机的节点）
  resetNodesNo: 不重置节点
  allow_ungraceful_removal: 允许删除未排空的节点
  allow_ungraceful_removal_desc: 删除节点时，kuboard-spray 会自动执行如下指令以排空节点，如果失败，是否继续删除节点？
  drain_timeout: 排空节点超时时间
  drain_retries: 排空节点重试次数
  drain_grace_period: 应用停止等候时间
  drain_retry_delay_seconds: 两次排空尝试间隔
  nodes_to_remove_required: 请选择要删除的节点

  newResourcePackageRequired: 当前资源包不支持此操作，请使用新的资源包
</i18n>

<template>
  <ExecuteTask v-if="!loading" :history="cluster.history" :loading="loading" :title="title" :startTask="applyPlan" @refresh="$emit('refresh')" @visibleChange="onVisibleChange">
    <el-form ref="form" :model="form" @submit.prevent.stop label-position="left" label-width="120px" class="app_form_mini">
      <!-- <div style="height: 10px;"></div> -->
      <el-form-item :label="$t('control_params')">
        <el-form-item :label="$t('verbose')">
          <el-switch v-model="form.verbose"></el-switch>
          <span style="width: 350px; margin-left: 110px; color: #aaa; font-size: 12px;">{{$t('verbose_' + form.verbose)}}</span>
        </el-form-item>
        <el-form-item :label="$t('vvv')">
          <el-switch v-model="form.vvv"></el-switch>
          <span style="width: 350px; margin-left: 110px; color: #aaa; font-size: 12px;">{{$t('vvv_' + form.vvv)}}</span>
        </el-form-item>
        <el-form-item :label="$t('fork')">
          <el-input-number v-model="form.fork" :step="2"></el-input-number>
          <span style="width: 350px; margin-left: 20px; color: #aaa; font-size: 12px;">{{$t('fork_more')}}</span>
        </el-form-item>
      </el-form-item>
      <el-form-item v-if="pendingRemoveNodes.length > 0 || pendingAddNodes.length > 0" :label="$t('operation')" prop="action" required style="margin-top: 10px;">
        <el-radio-group v-model="form.action">
          <el-radio-button v-if="pendingRemoveNodes.length >0" label="remove_node">
            {{ $t('aboutToRemoveNode') }}
          </el-radio-button>
          <div style="margin-bottom: 15px;" v-if="action === 'remove_node'">
            <div v-if="pingpong_loading" style="display: block;">
              <el-skeleton animated></el-skeleton>
            </div>
            <div v-else style="margin-left: 0px; margin-top: 10px;" class="app_form_mini">
              <el-form-item prop="remove_node.reset_nodes" required>
                <el-radio-group v-model="reset_nodes">
                  <el-radio-button :label="true">{{ $t('resetNodes') }}</el-radio-button>
                  <el-radio-button :label="false">{{ $t('resetNodesNo') }}</el-radio-button>
                </el-radio-group>
              </el-form-item>
              <el-form-item prop="remove_node.nodes_to_remove" :rules="nodes_to_remove_rules">
                <el-checkbox-group v-model="form.remove_node.nodes_to_remove">
                  <el-checkbox v-for="(item, key) in pendingRemoveNodes" :key="'h' + key" style="margin-top: 10px; margin-right: 10px;" :label="item.name"
                    :disabled="(pingpong[item.name] === undefined || pingpong[item.name].status !== 'SUCCESS') && form.remove_node.reset_nodes">
                    <el-tooltip v-if="pingpong[item.name] && pingpong[item.name].status !== 'SUCCESS'" :content="pingpong[item.name].message" placement="top-start">
                      <span class="app_text_mono">{{item.name}}</span>
                    </el-tooltip>
                    <span v-else class="app_text_mono">{{item.name}}</span>
                  </el-checkbox>
                </el-checkbox-group>
              </el-form-item>
              <el-form-item label-width="150px" :label="$t('allow_ungraceful_removal')" style="margin-top: 10px;">
                <el-switch v-model="form.remove_node.allow_ungraceful_removal"></el-switch>
                <div style="line-height: 20px; display: inline-block; vertical-align: top; margin-left: 10px;">
                  <div class="form_description">{{$t('allow_ungraceful_removal_desc')}}</div>
                  <div class="form_description app_text_mono" style="color: #666; font-weight: bold;">
                    kubectl drain --force --ignore-daemonsets --grace-period {{ form.remove_node.drain_grace_period * 60}} 
                      --timeout {{ form.remove_node.drain_timeout * 60 }}s
                  </div>
                </div>
              </el-form-item>
              <el-form-item label-width="150px" :label="$t('drain_grace_period')" required prop="remove_node.drain_grace_period">
                <el-input-number v-model="form.remove_node.drain_grace_period" :min="1" :max="form.remove_node.drain_timeout - 1"></el-input-number> 分钟
                <span style="margin-left: 20px;" class="form_description">kubectl drain --grace-period</span>
              </el-form-item>
              <el-form-item label-width="150px" :label="$t('drain_timeout')" required prop="remove_node.drain_timeout">
                <el-input-number v-model="form.remove_node.drain_timeout" :min="form.remove_node.drain_grace_period + 1"></el-input-number> 分钟
                <span style="margin-left: 20px;" class="form_description">kubectl drain --timeout</span>
              </el-form-item>
              <el-form-item label-width="150px" :label="$t('drain_retries')" required prop="remove_node.drain_retries">
                <el-input-number v-model="form.remove_node.drain_retries" :min="1"></el-input-number> 次
              </el-form-item>
              <el-form-item label-width="150px" :label="$t('drain_retry_delay_seconds')" required prop="remove_node.drain_retry_delay_seconds">
                <el-input-number v-model="form.remove_node.drain_retry_delay_seconds" :min="5"></el-input-number> 秒
              </el-form-item>
            </div>
          </div>
          <div>
          <el-radio-button v-if="pendingAddNodes.length > 0" label="add_node">
            {{ $t('aboutToAddNode') }}
            <div style="margin-left: 25px;">
              <el-tag v-for="(item, key) in pendingAddNodes" :key="'h' + key" style="margin-top: 10px; margin-right: 10px;">
                <span class="app_text_mono">{{item.name}}</span>
              </el-tag>
            </div>
          </el-radio-button>
          </div>
        </el-radio-group>
        <el-form-item v-if="cluster && cluster.resourcePackage && cluster.resourcePackage.data.supported_playbooks[this.action] === undefined" prop="min_resource_package_version" 
          style="margin-top: -10px;"
          :rules="min_resource_package_version_rules">
        </el-form-item>
      </el-form-item>
    </el-form>
  </ExecuteTask>
</template>

<script>
import ExecuteTask from '../common/task/ExecuteTask.vue'

function trimMark(str) {
  if (str[str.length - 1] === ',') {
    return str.slice(0, str.length - 1)
  }
  return str
}

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
        remove_node: {
          allow_ungraceful_removal: false,
          reset_nodes: true,
          nodes_to_remove: [],
          drain_grace_period: 5,
          drain_timeout: 6,
          drain_retries: 3,
          drain_retry_delay_seconds: 10,
        },
        min_resource_package_version: '',
      },
      min_resource_package_version_rules: [
        { required: true, message: this.$t('newResourcePackageRequired') }
      ],
      nodes_to_remove_rules: [
        { required: true, message: this.$t('nodes_to_remove_required') }
      ],
      pingpong: {},
      pingpong_loading: true,
    }
  },
  inject: ['isClusterInstalled', 'isClusterOnline', 'pendingRemoveNodes', 'pendingAddNodes'],
  computed: {
    reset_nodes: {
      get () { return this.form.remove_node.reset_nodes },
      set (v) {
        let temp = []
        for (let node of this.pendingRemoveNodes) {
          console.log(node.name, this.pingpong[node.name].status)
          if ((this.pingpong[node.name].status === 'SUCCESS') === v) {
            temp.push(node.name)
          }
        }
        console.log(temp)
        this.form.remove_node.nodes_to_remove = temp
        this.form.remove_node.reset_nodes = v
      }
    },
    action () {
      if (this.pendingRemoveNodes.length > 0 || this.pendingAddNodes.length > 0) {
        return this.form.action
      } else {
        return 'install_cluster'
      }
    },
    title () {
      return this.$t(this.action, {count: this.pendingAddNodes.length || this.pendingRemoveNodes.length })
    },
    hosts () {
      if (this.cluster && this.cluster.inventory) {
        return this.cluster.inventory.all.hosts
      }
      return {}
    },
    nodesToRemove: {
      get () { return this.pendingRemoveNodes },
      set () {}
    }
  },
  watch: {
    'hosts': {
      deep: true,
      handler: function () {
        this.updateForm(true)
      }
    },
    isClusterOnline () {
      this.updateForm(true)
    },
    isClusterInstalled () {
      this.updateForm(true)
    }
  },
  components: { ExecuteTask },
  emits: ['refresh'],
  mounted () {
    this.updateForm(true)
  },
  methods: {
    onVisibleChange (flag) {
      if (flag) {
        this.updateForm()
        if (this.pendingAddNodes.length > 0) {
          this.testPingPong(this.pendingAddNodes)
        } else if (this.pendingRemoveNodes.length > 0) {
          this.testPingPong(this.pendingRemoveNodes)
        }
      }
    },
    testPingPong (nodes) {
      this.pingpong = {}
      this.pingpong_loading = true
      let temp = ''
      for (let node of nodes) {
        temp += node.name + ','
      }
      temp = trimMark(temp, ',')
      let req = { nodes: temp }
      this.kuboardSprayApi.post(`/clusters/${this.name}/state/ping`, req).then(resp => {
        this.pingpong = resp.data.data.items
        this.pingpong_loading = false
      }).catch(e => {
        if (e.response && e.response.data) {
          this.$message.error('不能测试节点是否在线: ' + e.response.data.message)
        } else {
          this.$message.error('不能测试节点是否在线: ' + e)
        }
      })
    },
    updateForm () {
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
        this.form.remove_node.nodes_to_remove = []
        this.form.remove_node.reset_nodes = true
        this.form.remove_node.allow_ungraceful_removal = false
        this.form.remove_node.drain_grace_period = 5
        this.form.remove_node.drain_timeout = 6
        this.form.remove_node.drain_retries = 3
        this.form.remove_node.drain_retry_delay_seconds = 10
      } else if (this.pendingAddNodes.length > 0) {
        this.form.action = 'add_node'
      }
    },
    async applyPlan () {
      return new Promise((resolve, reject) => {
        this.$refs.form.validate(flag => {
          if (flag) {
            let req = {
              verbose: this.form.verbose,
              vvv: this.form.vvv,
              fork: this.form.fork,
            }
            if (this.action === 'remove_node') {
              let temp = ''
              for (let node of this.form.remove_node.nodes_to_remove) {
                temp += node + ','
              }
              temp = trimMark(temp)
              req.nodes_to_remove = temp
              req.reset_nodes = this.form.remove_node.reset_nodes
              req.allow_ungraceful_removal = this.form.remove_node.allow_ungraceful_removal
              req.drain_grace_period = this.form.remove_node.drain_grace_period * 60 + ''
              req.drain_timeout = this.form.remove_node.drain_timeout * 60 + 's'
              req.drain_retries = this.form.remove_node.drain_retries + ''
              req.drain_retry_delay_seconds = this.form.remove_node.drain_retry_delay_seconds + ''
            }
            this.kuboardSprayApi.post(`/clusters/${this.name}/${this.action}`, req).then(resp => {
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
.form_description {
  font-size: 12px;
  // color: var(--el-text-color-placeholder);
  color: #aaa;

}
</style>
