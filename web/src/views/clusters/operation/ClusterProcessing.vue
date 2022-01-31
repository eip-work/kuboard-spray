<i18n>
en:
  verbose: Include task params
  verbose_: 
  verbose_v: May include sensitive data in the trace, e.g. path to files, user name, password.
  verbose_vvv: includes more information in log, only used in development.
  v_: Normal
  v_v: Details
  v_vvv: More
  fork: ansible fork
  fork_more: Max number of nodes can be operated in the installation.
  control_params: Control params
  operation: Operations
  install_cluster: Install / Setup K8S Cluster
  add_node: "{count} nodes waiting to be added"
  remove_node: "{count} nodes waiting to be removed"
  aboutToRemoveNode: Remove nodes
  aboutToAddNode: Add nodes
  offlineNodes: Offline nodes
  offlineNodesDesc: Exclude the following offline nodes.
  excludeNodes: Exclude nodes
  excludeNodesDesc: Exclude the following nodes.

  sync_etcd_address: "Update apiserver's --etcd-servers param"
  sync_nginx_config: "Update apiserver list in loadbalancer"

  newResourcePackageRequired: This resource package does not support the operation, please choose a new one.
zh:
  verbose: 显示任务参数
  verbose_: 正常输出的日志，通常选用此选项
  verbose_v: 日志中会包含部分敏感信息，例如：文件路径、用户名密码等
  verbose_vvv: 日志中会包含最详细的信息，通常只在开发阶段使用
  v_: 正常
  v_v: 详细
  v_vvv: 更多
  fork: 并发数量
  fork_more: 安装过程中可以同时操作的目标节点的最大数量。ansible fork.
  control_params: 控制选项
  operation: 操作选项
  install_cluster: 安装 / 设置集群
  add_node: "{count} 个节点待添加"
  remove_node: "{count} 个节点待删除"
  aboutToRemoveNode: 删除节点
  aboutToAddNode: 添加节点
  offlineNodes: 离线节点
  offlineNodesDesc: 排除以下不在线的节点：
  excludeNodes: 排除节点
  excludeNodesDesc: 排除以下节点

  sync_etcd_address: "更新 apiserver 中 etcd 连接参数"
  sync_nginx_config: "更新负载均衡中 apiserver 列表"

  newResourcePackageRequired: 当前资源包不支持此操作，请使用新的资源包
</i18n>

<template>
  <ExecuteTask v-if="!loading" :history="cluster.history" :loading="loading" :title="title" :startTask="applyPlan" @refresh="$emit('refresh')" @visibleChange="onVisibleChange">
    <div v-if="pingpong_loading" style="display: block; min-width: 420px;">
      <el-skeleton animated></el-skeleton>
    </div>
    <el-form v-else ref="form" :model="form" @submit.prevent.stop label-position="left" label-width="120px" class="app_form_mini">
      <!-- <div style="height: 10px;"></div> -->
      <el-form-item :label="$t('control_params')">
        <el-form-item :label="$t('verbose')">
          <el-radio-group v-model="form.verbose">
            <el-radio-button label="">{{$t('v_')}}</el-radio-button>
            <el-radio-button label="v">{{$t('v_v')}}</el-radio-button>
            <el-radio-button label="vvv">{{$t('v_vvv')}}</el-radio-button>
          </el-radio-group>
          <span style="width: 350px; margin-left: 20px; color: #aaa; font-size: 12px;">{{$t('verbose_' + form.verbose)}}</span>
        </el-form-item>
        <el-form-item :label="$t('fork')" style="margin-top: 10px;">
          <el-input-number v-model="form.fork" :step="2" style="width: 166px;"></el-input-number>
          <span style="width: 350px; margin-left: 20px; color: #aaa; font-size: 12px;">{{$t('fork_more')}}</span>
        </el-form-item>
      </el-form-item>
      <el-form-item :label="$t('operation')" prop="action" style="margin-top: 10px;">
        <el-radio-group v-model="action" class="app_margin_bottom">
          <el-radio-button :disabled="action !== 'install_cluster'" label="install_cluster">
            {{ $t('install_cluster') }}
          </el-radio-button>
          <el-radio-button :disabled="pendingRemoveNodes.length == 0" label="remove_node">
            {{ $t('aboutToRemoveNode') }}
          </el-radio-button>
          <el-radio-button :disabled="pendingAddNodes.length == 0" label="add_node">
            {{ $t('aboutToAddNode') }}
          </el-radio-button>
          <el-radio-button :disabled="!cluster.inventory.all.hosts.localhost.kuboardspray_sync_nginx_config" label="sync_nginx_config">
            {{ $t('sync_nginx_config') }}
          </el-radio-button>
          <el-radio-button :disabled="!cluster.inventory.all.hosts.localhost.kuboardspray_sync_etcd_address" label="sync_etcd_address">
            {{ $t('sync_etcd_address') }}
          </el-radio-button>
        </el-radio-group>

        <Component :is="action" ref="action" v-if="cluster"
          :cluster="cluster" 
          :name="name" 
          :form="form"
          :hosts="hosts"
          :pingpong="pingpong"
          :pingpong_loading="pingpong_loading"
          :offlineNodes="offlineNodes">
        </Component>
      </el-form-item>
      <el-form-item :label="$t('offlineNodes')" class="app_margin_top" v-if="offlineNodes.length > 0">
        <div class="form_description">{{ $t('offlineNodesDesc') }}</div>
        <template v-for="node in offlineNodes" :key="'exclude' + node">
          <el-tooltip class="box-item" effect="dark" :content="pingpong[node].message" placement="top-end">
            <el-tag type="error" effect="dark" style="margin: 0 10px 10px 0;">
              <span class="app_text_mono" style="font-size: 14px; margin-right: 10px;">{{ node }}</span> <i class="el-icon-question"></i>
            </el-tag>
          </el-tooltip>
        </template>
      </el-form-item>
      <el-form-item v-if="cluster && cluster.resourcePackage && cluster.resourcePackage.data.supported_playbooks[this.action] === undefined" prop="min_resource_package_version" 
        style="margin-top: -10px;"
        :rules="min_resource_package_version_rules">
      </el-form-item>
    </el-form>
  </ExecuteTask>
</template>

<script>
import ExecuteTask from '../../common/task/ExecuteTask.vue'
import clone from 'clone'

import add_node from './cluster/add_node.vue'
import install_cluster from './cluster/install_cluster.vue'
import remove_node from './cluster/remove_node.vue'
import sync_etcd_address from './cluster/sync_etcd_address.vue'
import sync_nginx_config from './cluster/sync_nginx_config.vue'

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
        verbose: '',
        fork: 5,
        action: '',
        min_resource_package_version: '',
        nodes_to_exclude: []
      },
      min_resource_package_version_rules: [
        { required: true, message: this.$t('newResourcePackageRequired') }
      ],
      pingpong: {},
      pingpong_loading: true,
    }
  },
  inject: ['isClusterInstalled', 'isClusterOnline', 'pendingRemoveNodes', 'pendingAddNodes'],
  computed: {
    offlineNodes () {
      let result = []
      if (this.action === 'remove_node' && this.form.remove_node && this.form.remove_node.reset_nodes) {
        return result
      }
      for (let key in this.cluster.inventory.all.hosts) {
        if (this.pingpong[key] && this.pingpong[key].status !== 'SUCCESS') {
          result.push(key)
        }
      }
      return result
    },
    action: {
      get () {
        return this.form.action
      },
      set (v) { this.form.action = v }
    },
    title () {
      return this.$t(this.action, { count: this.pendingAddNodes.length || this.pendingRemoveNodes.length })
    },
    hosts () {
      if (this.cluster && this.cluster.inventory) {
        let temp = clone(this.cluster.inventory.all.hosts)
        delete temp.localhost
        delete temp.bastion
        return temp
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
    },
  },
  components: { 
    ExecuteTask,
    add_node, install_cluster, remove_node, sync_etcd_address, sync_nginx_config,
  },
  emits: ['refresh'],
  mounted () {
    this.updateForm(true)
  },
  methods: {
    onVisibleChange (flag) {
      if (flag) {
        this.updateForm()
        this.testPingPong('target')
      }
    },
    testPingPong (nodes) {
      this.pingpong = {}
      this.pingpong_loading = true
      let req = { nodes: nodes }
      this.kuboardSprayApi.post(`/clusters/${this.name}/state/ping`, req).then(resp => {
        this.pingpong = resp.data.data.items
        this.pingpong_loading = false
      }).catch(e => {
        if (e.response && e.response.data) {
          this.$message.error('不能测试节点是否在线: ' + e.response.data.message)
        } else {
          this.$message.error('不能测试节点是否在线: ' + e)
        }
        this.pingpong_loading = false
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
      if (this.cluster.inventory.all.hosts.localhost.kuboardspray_sync_etcd_address) {
        this.form.action = 'sync_etcd_address'
      } else if (this.cluster.inventory.all.hosts.localhost.kuboardspray_sync_nginx_config) {
        this.form.action = 'sync_nginx_config'
      } else if (this.etcdInDanger) {
          this.form.action = 'add_node'
      } else if (this.pendingRemoveNodes.length > 0) {
        this.form.action = 'remove_node'
      } else if (this.pendingAddNodes.length > 0) {
        this.form.action = 'add_node'
      } else {
        this.form.action = 'install_cluster'
      }
    },
    async applyPlan () {
      if (this.pingpong_loading) {
        this.$message.error('Wait ..')
        return
      }
      return new Promise((resolve, reject) => {
        this.$refs.form.validate(flag => {
          if (flag) {
            let req = {
              verbose: this.form.verbose,
              fork: this.form.fork,
            }
            let discovered_interpreter_python = {}
            for (let node in this.pingpong) {
              let msg = this.pingpong[node].message
              if (msg) {
                let pong = JSON.parse(msg)
                if (pong && pong.ansible_facts && pong.ansible_facts.discovered_interpreter_python) {
                  discovered_interpreter_python[node] = pong.ansible_facts.discovered_interpreter_python
                }
              }
            }
            req.discovered_interpreter_python = discovered_interpreter_python
            { // 排除节点
              let temp = ''
              let excludes = {}
              for (let node of this.offlineNodes) {
                excludes[node] = true
              }
              for (let i in this.$refs.action.excludeNodes) {
                excludes[this.$refs.action.excludeNodes[i]] = true
              }
              for (let node in excludes) {
                temp += '!' + node + ','
              }
              req.nodes_to_exclude = trimMark(temp)
            }
            req = Object.assign(req, this.$refs.action.populateRequest())
            console.log(req)
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
  max-width: 700px;
}
</style>
