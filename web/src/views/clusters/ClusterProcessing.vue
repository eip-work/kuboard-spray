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
  aboutToRemoveNode: Remove nodes
  aboutToAddNode: Add nodes
  resetNodes: Reset nodes(non-applicable to offline nodes)
  resetNodesNo: Do not reset nodes
  allow_ungraceful_removal: Allow remove nodes failed to drain
  allow_ungraceful_removal_desc: 'kuboard-spray does "kubectl drain" before removing nodes, continue on failure?'
  drain_timeout: drain_timeout
  drain_retries: drain_retries
  drain_grace_period: drain_grace_period
  drain_retry_delay_seconds: drain_retry_delay_seconds
  nodes_to_remove_required: Please select nodes to remove.
  nodesToExclude: Exclude nodes.
  nodesToExcludeDesc: Exclude the following offline nodes.
  nodesToIncludeDesc: Include the following nodes.
  allNodesMustBeOnline: All nodes must be reachable.
  removeOneEtcdNodeOnce: If you want to remove etcd node, can only remove them one by one.
  requiresAtLeastOneOnlineNode: All nodes pending to add are offline.
  requiresAllControlNodeOnline: "{node} is offline. Remove it, or bring it back."
  requireAtLeastOneControlPlane: Requires at least one control_plane.
  requiresOddEtcdCount: Etcd count should be an odd number.
  requiresAllEtcdNodeOnline: "All Etcd nodes must be online, currently, we cannot reach node {node}"
  requiresKubeNodeCount: Requires at lease one kube_node.
  sync_etcd_address: "Update apiserver's --etcd-servers param"
  sync_etcd_address_desc: "Member list of etcd cluster changed, this task is going to update --etcd-servers param in /etc/kubernetes/manifests/kube-apiserver.yaml on all the remaining kube_control_plane nodes, to match the latest etcd member list."
  sync_nginx_config: "Update apiserver list in loadbalancer"
  sync_nginx_config_desc: "Control_plane list in k8s cluster changed, this task is going to update 'upstream kube_apiserver' block in /etc/nginx/nginx.conf on all the kube_node nodes."

  add_nodes_desc: You are going to add the following nodes into the cluster.

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
  aboutToRemoveNode: 删除节点
  aboutToAddNode: 添加节点
  resetNodes: 重置节点（不能选择已停机的节点）
  resetNodesNo: 不重置节点
  allow_ungraceful_removal: 允许删除未排空的节点
  allow_ungraceful_removal_desc: 删除节点时，kuboard-spray 会自动执行如下指令以排空节点，如果失败，是否继续删除节点？
  drain_timeout: 排空节点超时时间
  drain_retries: 排空节点重试次数
  drain_grace_period: 应用停止等候时间
  drain_retry_delay_seconds: 两次排空尝试间隔
  nodes_to_remove_required: 请选择要删除的节点
  nodesToExclude: 排除节点
  nodesToExcludeDesc: 排除以下不在线的节点：
  nodesToIncludeDesc: 包含以下节点：
  allNodesMustBeOnline: 所有节点必须在线，请删除不在线节点或者使其在线
  removeOneEtcdNodeOnce: 如果要删除 etcd 节点，一次只能选择一个节点删除。
  requiresAtLeastOneOnlineNode: 至少需要一个在线节点
  requiresAllControlNodeOnline: "安装集群或者添加控制节点时，所有控制节点必须在线。请启动 {node} 或者将其移除。"
  requireAtLeastOneControlPlane: 至少需要一个控制节点
  requiresOddEtcdCount: ETCD 节点的数量必须为奇数
  requiresAllEtcdNodeOnline: "所有 ETCD 节点必须在线，当前 {node} 不在线"
  requiresKubeNodeCount: 至少要有一个在线的工作节点
  sync_etcd_address: "更新 apiserver 中 etcd 连接参数"
  sync_etcd_address_desc: "ETCD 集群的成员列表已经发生变化，此操作将更新剩余控制节点上 /etc/kubernetes/manifests/kube-apiserver.yaml 文件中 --etcd-servers 的参数，以符合匹配最新的 etcd 成员列表。"
  sync_nginx_config: "更新负载均衡中 apiserver 列表"
  sync_nginx_config_desc: "集群中控制节点的列表发生变化，此操作将更新所有工作节点上 /etc/nginx/nginx.conf 文件中 upstream kube_apiserver 的列表"

  add_nodes_desc: 将要添加以下节点


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
      <el-form-item :label="$t('operation')" prop="action" style="margin-top: 10px;" :rules="nodes_to_exclude_rules">
        <el-radio-group v-model="action">
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

        <!-- 安装集群的参数 -->
        <div v-if="action === 'install_cluster'">
          <div class="form_description">{{ $t('nodesToIncludeDesc') }}</div>
          <template v-for="(item, key) in hosts" :key="'h' + key">
            <el-tag v-if="pingpong[key] && pingpong[key].status === 'SUCCESS'" style="margin-bottom: 10px; margin-right: 10px;" effect="dark">
              <span class="app_text_mono" style="font-size: 14px;">{{key}}</span>
            </el-tag>
          </template>
        </div>

        <!-- 添加节点的参数 -->
        <div v-if="action === 'add_node'">
          <div class="form_description">
            {{ $t('add_nodes_desc') }}
          </div>
          <template v-for="(item, key) in pendingAddNodes" :key="'h' + key">
            <el-tag v-if="pingpong[item.name] && pingpong[item.name].status === 'SUCCESS'" style="margin-top: 10px; margin-right: 10px;" :disabled="pingpong[item.name].status !== 'SUCCESS'" :label="item.name">
              <span class="app_text_mono" style="font-size: 14px;">{{item.name}}</span>
            </el-tag>
          </template>
        </div>

        <!-- 删除节点的参数 -->
        <div style="margin-bottom: 15px;" v-if="action === 'remove_node'">
          <div style="margin-left: 0px; margin-top: 10px;" class="app_form_mini">
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
              <el-input-number v-model="form.remove_node.drain_retry_delay_seconds" :min="5" :step="5"></el-input-number> 秒
            </el-form-item>
          </div>
        </div>

        <!-- 同步配置的参数 -->
        <div v-if="action === 'sync_nginx_config'" style="margin-top: 10px;">
          <div class="form_description" style="margin-bottom: 10px;">
            {{ $t('sync_nginx_config_desc') }}
          </div>
        </div>

        <div v-if="action === 'sync_etcd_address'" style="margin-top: 10px;">
          <div class="form_description" style="margin-bottom: 10px;">
            {{ $t('sync_etcd_address_desc') }}
          </div>
        </div>
      </el-form-item>
      <el-form-item :label="$t('nodesToExclude')" class="app_margin_top" v-if="nodesToExclude.length > 0">
        <div class="form_description">{{ $t('nodesToExcludeDesc') }}</div>
        <template v-for="node in nodesToExclude" :key="'exclude' + node">
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
import ExecuteTask from '../common/task/ExecuteTask.vue'
import clone from 'clone'

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
        nodes_to_exclude: []
      },
      min_resource_package_version_rules: [
        { required: true, message: this.$t('newResourcePackageRequired') }
      ],
      nodes_to_remove_rules: [
        { required: true, message: this.$t('nodes_to_remove_required') }
      ],
      nodes_to_exclude_rules: [
        {
          validator: (rule, value, callback) => {
            if (this.action === 'install_cluster' || this.action === 'add_node') {
              // if (this.nodesToExclude.length > 0) {
              //   return callback(this.$t('allNodesMustBeOnline'))
              // }
              
              // 至少有一个被添加的节点在线
              let temp = ''
              for (let node in this.cluster.inventory.all.hosts) {
                if (this.pingpong[node] && this.pingpong[node].status === 'SUCCESS') {
                  if (this.action === 'add_node') {
                    if (this.cluster.inventory.all.hosts[node].kuboardspray_node_action === 'add_node') {
                      temp += node + ','
                    }
                  } else if (this.action === 'install_cluster') {
                    temp += node + ','
                  }
                }
              }
              if (temp === '') {
                return callback(this.$t('requiresAtLeastOneOnlineNode'))
              }

              // 至少有一个工作节点
              let countKubeNodeCount = 0
              for (let node in this.cluster.inventory.all.children.target.children.k8s_cluster.children.kube_node.hosts) {
                if (this.pingpong[node] && this.pingpong[node].status === 'SUCCESS') {
                  countKubeNodeCount ++
                }
              }
              if (countKubeNodeCount === 0) {
                return callback(this.$t('requiresKubeNodeCount'))
              }


              // 所有的控制节点必须在线
              let controlPlaneCount = 0
              for (let controlPlane in this.cluster.inventory.all.children.target.children.k8s_cluster.children.kube_control_plane.hosts) {
                if (this.pingpong[controlPlane] === undefined || this.pingpong[controlPlane].status !== 'SUCCESS') {
                  return callback(this.$t('requiresAllControlNodeOnline', { node: controlPlane }))
                }
                controlPlaneCount ++
              }
              if (controlPlaneCount === 0) {
                return callback(this.$t('requireAtLeastOneControlPlane'))
              }
              // 所有的 etcd 节点必须在线
              let etcdCount = 0
              for (let etcd in this.cluster.inventory.all.children.target.children.etcd.hosts) {
                if (this.pingpong[etcd] === undefined || this.pingpong[etcd].status !== 'SUCCESS') {
                  return callback(this.$t('requiresAllEtcdNodeOnline, { node: etcd }'))
                }
                etcdCount ++
              }
              if (etcdCount % 2 == 0) {
                return callback(this.$t('requiresOddEtcdCount'))
              }
            } else if (this.action === 'remove_node') {
              for (let node of this.form.remove_node.nodes_to_remove) {
                if (this.cluster.inventory.all.children.target.children.etcd.hosts[node] !== undefined && this.form.remove_node.nodes_to_remove.length > 1) {
                  return callback(this.$t('removeOneEtcdNodeOnce'))
                }
              }
            }
            return callback()
          },
          trigger: 'change',
        }
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
          if ((this.pingpong[node.name].status === 'SUCCESS') === v) {
            temp.push(node.name)
          }
        }
        this.form.remove_node.nodes_to_remove = temp
        this.form.remove_node.reset_nodes = v
      }
    },
    nodesToExclude () {
      let result = []
      if (this.action === 'remove_node' && !this.reset_nodes) {
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
        this.form.remove_node.nodes_to_remove = []
        this.form.remove_node.reset_nodes = true
        this.form.remove_node.allow_ungraceful_removal = false
        this.form.remove_node.drain_grace_period = 5
        this.form.remove_node.drain_timeout = 6
        this.form.remove_node.drain_retries = 3
        this.form.remove_node.drain_retry_delay_seconds = 10
      } else if (this.pendingAddNodes.length > 0) {
        this.form.action = 'add_node'
      } else {
        this.form.action = 'install_cluster'
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
            if (this.nodesToExclude.length > 0) {
              let temp = ''
              for (let node of this.nodesToExclude) {
                temp += '!' + node + ','
              }
              req.nodes_to_exclude = trimMark(temp)
            }
            if (this.action === 'install_cluster') {
              let temp = ''
              for (let node in this.hosts) {
                if (this.nodesToExclude.indexOf(node) < 0) {
                  temp += node + ','
                }
              }
              req.nodes = trimMark(temp)
            } else if (this.action === 'add_node') {
              let temp = ''
              for (let node of this.pendingAddNodes) {
                if (this.pingpong[node.name] && this.pingpong[node.name].status === 'SUCCESS') {
                  temp += node.name + ','
                }
              }
              req.nodes_to_add = trimMark(temp)
            } else if (this.action === 'remove_node') {
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
            } else if (this.action === 'sync_nginx_config') {
              console.log('sync_nginx_config')
            } else if (this.action === 'sync_etcd_address') {
              console.log('sync_etcd_address')
            }
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

}
</style>
