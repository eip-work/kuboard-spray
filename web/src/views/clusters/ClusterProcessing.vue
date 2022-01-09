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
  <ExecuteTask v-if="!loading" :history="cluster.history" :loading="loading" :title="title" :startTask="applyPlan" @refresh="$emit('refresh')" @visibleChange="onVisibleChange">
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
              <div v-if="pingpong_loading" style="display: block;">
                <el-skeleton animated></el-skeleton>
              </div>
              <div v-else style="margin-left: 25px; margin-top: 10px;" class="app_form_mini">
                <el-form-item prop="reset_nodes" required>
                  <el-radio-group v-model="reset_nodes">
                    <el-radio-button :label="true">重置节点（不能选择已停机的节点）</el-radio-button>
                    <el-radio-button :label="false">不重置节点</el-radio-button>
                  </el-radio-group>
                </el-form-item>
                <el-form-item prop="nodes_to_remove" required>
                  <el-checkbox-group v-model="form.nodes_to_remove">
                    <el-checkbox v-for="(item, key) in pendingRemoveNodes" :key="'h' + key" style="margin-top: 10px; margin-right: 10px;" :label="item.name"
                      :disabled="(pingpong[item.name] === undefined || pingpong[item.name].status !== 'SUCCESS') && form.reset_nodes">
                      <el-tooltip v-if="pingpong[item.name] && pingpong[item.name].status !== 'SUCCESS'" :content="pingpong[item.name].message" placement="top-start">
                        <span class="app_text_mono">{{item.name}}</span>
                      </el-tooltip>
                      <span v-else class="app_text_mono">{{item.name}}</span>
                    </el-checkbox>
                  </el-checkbox-group>
                </el-form-item>
              </div>
              <el-form-item label="排空等候时间" required prop="drain_out_time" style="margin-top: 20px;">
                <el-input-number v-model="form.drain_out_time" :min="1"></el-input-number> 分钟
              </el-form-item>
              <el-form-item label="排空重试次数" required prop="drain_retries">
                <el-input-number v-model="form.drain_retries" :min="1"></el-input-number> 次
              </el-form-item>
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
        reset_nodes: true,
        nodes_to_remove: [],
        drain_out_time: 10,
        drain_retries: 3,
        min_resource_package_version: '',
      },
      min_resource_package_version_rules: [
        { required: true, message: '当前资源包不支持此操作，请使用新的资源包' }
      ],
      pingpong: {},
      pingpong_loading: true,
    }
  },
  inject: ['isClusterInstalled', 'isClusterOnline', 'pendingRemoveNodes', 'pendingAddNodes'],
  computed: {
    reset_nodes: {
      get () { return this.form.reset_nodes },
      set (v) {
        let temp = []
        for (let node of this.pendingRemoveNodes) {
          console.log(node.name, this.pingpong[node.name].status)
          if ((this.pingpong[node.name].status === 'SUCCESS') === v) {
            temp.push(node.name)
          }
        }
        console.log(temp)
        this.form.nodes_to_remove = temp
        this.form.reset_nodes = v
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
        this.form.nodes_to_remove = []
        this.form.reset_nodes = true
        this.form.drain_out_time = 10
        this.form.drain_retries = 3
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
              req.reset_nodes = this.form.reset_nodes
              let temp = ''
              for (let node of this.form.nodes_to_remove) {
                temp += node + ','
              }
              temp = trimMark(temp)
              req.nodes_to_remove = temp
              req.drain_out_time = this.form.drain_out_time * 60 + 's'
              req.drain_retries = this.form.drain_retries + ''
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
</style>
