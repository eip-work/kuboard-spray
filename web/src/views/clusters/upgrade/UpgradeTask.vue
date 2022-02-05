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
  offlineNodes: Offline nodes
  offlineNodesDesc: Exclude the following offline nodes.

  download_binaries: Distribute installation binaries
  download_binaries_desc: In the process of upgrading, we are going to change image version in some daemonsets, it would be the best if we distribute all required binaries and container images before we do the cactual upgrading.

  skip_downloads: Skip distributing installation binaries

  upgrade_cluster: "Upgrade cluster"
  upgrade_multi_nodes: Upgrade kube_nodes
  upgrade_multi_nodes_desc: Please select nodes to upgrade in this task run.
  upgrade_single_node: Upgrade {nodeName}
  upgrade_single_node_desc: Upgrade one node per task run.
  
  upgrade_all_nodes: Upgrade all nodes in one run
  upgrade_master_nodes: Upgrade kube_control_plane and etcd nodes first
zh:
  verbose: 显示任务参数
  verbose_: 正常输出的日志，通常选用此选项
  verbose_v: 日志中会包含部分敏感信息，例如：文件路径、用户名密码等
  verbose_vvv: 日志中会包含最详细的信息，通常只在开发阶段使用
  v_: 正常
  v_v: 详细
  v_vvv: 更多
  fork: 并发数量
  fork_more: 安装过程中可以同时操作的目标节点的最大数量。
  control_params: 控制选项
  operation: 操作选项
  offlineNodes: 离线节点
  offlineNodesDesc: 排除以下不在线的节点：

  download_binaries: 分发安装包
  download_binaries_desc: 升级过程中，将会改变一些 daemonset 的版本（例如：kube-proxy, calico-node 等），如果不事先分发这些镜像到所有各个节点，将导致升级过程中，部分节点可能因为缺少对应版本的镜像而出现问题。

  skip_downloads: 跳过加载安装包

  upgrade_cluster: 升级集群
  upgrade_multi_nodes: 升级多个工作节点
  upgrade_multi_nodes_desc: 请选择本次任务中要升级的节点
  upgrade_single_node: 升级 {nodeName}
  upgrade_single_node_desc: 一次只升级一个节点

  upgrade_all_nodes: 升级所有节点
  upgrade_all_nodes_desc: 一次性升级所有节点，如果当前有工作负载正在运行，将造成服务中断
  upgrade_master_nodes: 升级控制节点和 ETCD 节点
  upgrade_master_nodes_desc: 先逐个升级控制节点和 ETCD 节点，然后再逐个排空、升级工作节点，避免造成服务中断。（如果有3个控制节点，大约耗时 18 分钟）
</i18n>

<template>
  <ExecuteTask :history="cluster.history" :loading="loading" :title="title" :startTask="execute" @refresh="$emit('refresh')" @visibleChange="onVisibleChange">
    <div style="width: 600px;">
      <div v-if="pingpong_loading" style="display: block;">
        <el-skeleton animated></el-skeleton>
      </div>
      <el-form v-else ref="form" :model="form" @submit.prevent.stop label-position="left" label-width="120px" class="app_form_mini">
        <el-form-item :label="$t('control_params')">
          <el-form-item :label="$t('verbose')">
            <el-radio-group v-model="form.verbose">
              <el-radio-button label="">{{$t('v_')}}</el-radio-button>
              <el-radio-button label="v">{{$t('v_v')}}</el-radio-button>
              <el-radio-button label="vvv">{{$t('v_vvv')}}</el-radio-button>
            </el-radio-group>
            <div style="color: #aaa; font-size: 12px;">{{$t('verbose_' + form.verbose)}}</div>
          </el-form-item>
          <el-form-item :label="$t('fork')" style="margin-top: 10px;">
            <el-input-number v-model="form.fork" :step="2" style="width: 166px;"></el-input-number>
            <div style="color: #aaa; font-size: 12px;">{{$t('fork_more')}}</div>
          </el-form-item>
        </el-form-item>
        <el-form-item :label="$t('operation')" prop="action" style="margin-top: 10px;">
          <el-radio-group v-model="form.action" class="app_margin_bottom">
            <template v-if="nodeName === undefined">
              <el-radio-button label="download_binaries" :disabled="!requireSeparateDownloadAction">
                {{ $t('download_binaries') }}
              </el-radio-button>
              <el-radio-button label="upgrade_all_nodes" :disabled="requireSeparateDownloadAction || !controlPlanePendingUpgrade">
                {{ $t('upgrade_all_nodes') }}
              </el-radio-button>
              <el-radio-button label="upgrade_master_nodes" :disabled="requireSeparateDownloadAction || !controlPlanePendingUpgrade">
                {{ $t('upgrade_master_nodes') }}
              </el-radio-button>
              <el-radio-button label="upgrade_multi_nodes" :disabled="requireSeparateDownloadAction || controlPlanePendingUpgrade">{{$t('upgrade_multi_nodes')}}</el-radio-button>
            </template>
            <el-radio-button v-else label="upgrade_single_node">{{$t('upgrade_single_node', {nodeName})}}</el-radio-button>
          </el-radio-group>
          <div class="form_description" :style="form.action === 'upgrade_all_nodes' ? 'color: var(--el-color-danger)':''">{{ $t(form.action + '_desc') }}</div>
          <div v-if="form.action === 'upgrade_multi_nodes'" :rules="kube_nodes_to_upgrade_rules" style="margin-bottom: 18px;">
            <el-form-item prop="kube_nodes_to_upgrade" required>
              <el-checkbox-group v-model="form.kube_nodes_to_upgrade" @change="form.skip_downloads = !requireDownloadForNodes(form.kube_nodes_to_upgrade)">
                <template v-for="(node, nodeName) in cluster.inventory.all.children.target.children.k8s_cluster.children.kube_node.hosts" :key="'nu' + nodeName">
                  <el-checkbox v-if="cluster.inventory.all.hosts[nodeName].kuboardspray_node_action === 'upgrade_node'" :label="nodeName">{{nodeName}}</el-checkbox>
                </template>
              </el-checkbox-group>
            </el-form-item>
          </div>
          <div v-if="form.action === 'upgrade_multi_nodes' || form.action === 'upgrade_single_node'">
            <el-switch v-model="form.skip_downloads" :disabled="requireDownloadForNodes(form.kube_nodes_to_upgrade)" :active-text="$t('skip_downloads')"></el-switch>
          </div>
        </el-form-item>

        <el-form-item :label="$t('offlineNodes')" class="app_margin_top" v-if="offlineNodes.length > 0">
          <div class="form_description">{{ $t('offlineNodesDesc') }}</div>
          <template v-for="node in offlineNodes" :key="'exclude' + node">
            <el-tooltip class="box-item" effect="dark" :content="pingpong[node].message" placement="top-end">
              <el-tag type="danger" effect="dark" style="margin: 0 10px 10px 0;">
                <span class="app_text_mono" style="font-size: 14px; margin-right: 10px;">{{ node }}</span>
                <el-icon :size="14" style="width: 14px; height: 14px; vertical-align: top;">
                  <el-icon-question-filled></el-icon-question-filled>
                </el-icon>
              </el-tag>
            </el-tooltip>
          </template>
        </el-form-item>

        <el-form-item v-if="cluster && cluster.resourcePackage && cluster.resourcePackage.data.supported_playbooks['upgrade_cluster'] === undefined" prop="min_resource_package_version" 
          style="margin-top: -10px;"
          :rules="min_resource_package_version_rules">
        </el-form-item>
      </el-form>
    </div>
  </ExecuteTask>
</template>

<script>
import ExecuteTask from '../../common/task/ExecuteTask.vue'

function trimMark(str) {
  if (str[str.length - 1] === ',') {
    return str.slice(0, str.length - 1)
  }
  return str
}

export default {
  props: {
    cluster: { type: Object, required: true },
    nodeName: { type: String, required: false, default: undefined },
    controlPlanePendingUpgrade: { type: Boolean, required: false, default: true },
  },
  data() {
    return {
      loading: false,
      form: {
        verbose: '',
        fork: 5,
        action: 'upgrade_master_nodes',
        min_resource_package_version: '',
        nodes_to_exclude: [],
        kube_nodes_to_upgrade: [],
        skip_downloads: false,
      },
      min_resource_package_version_rules: [
        { required: true, message: this.$t('newResourcePackageRequired') }
      ],
      kube_nodes_to_upgrade_rules: [
        { required: true, message: this.$t('upgrade_multi_nodes_desc') }
      ],
      pingpong: {},
      pingpong_loading: true,
    }
  },
  computed: {
    title () {
      if (this.nodeName) {
        return this.$t('upgrade_single_node', { nodeName: this.nodeName })
      }
      return this.$t('upgrade_cluster') + ' : ' + this.$t(this.form.action)
    },
    requireSeparateDownloadAction () {
      for (let hostName in this.cluster.inventory.all.children.target.children.k8s_cluster.children.kube_control_plane.hosts) {
        let host = this.cluster.inventory.all.hosts[hostName]
        if (host.kuboardspray_require_download) {
          return true
        }
      }
      return false
    },
    offlineNodes () {
      let result = []
      for (let key in this.cluster.inventory.all.hosts) {
        if (this.pingpong[key] && this.pingpong[key].ping !== 'pong') {
          result.push(key)
        }
      }
      return result
    },
  },
  components: { ExecuteTask },
  mounted () {
    this.setAction()
  },
  methods: {
    requireDownloadForNodes(nodes) {
      for (let node of nodes) {
        if (this.cluster.inventory.all.hosts[node].kuboardspray_require_download) {
          return true
        }
      }
      return false
    },
    onVisibleChange(flag) {
      if (flag) {
        this.setAction()
        this.form.kube_nodes_to_upgrade = []
        this.testPingPong('target')
      }
    },
    setAction () {
      if (!this.controlPlanePendingUpgrade) {
        if (this.nodeName === undefined) {
          this.form.action = 'upgrade_multi_nodes'
          this.form.skip_downloads = false
        } else {
          this.form.action = 'upgrade_single_node'
          this.form.skip_downloads = !this.requireDownloadForNodes([this.nodeName])
          console.log(this.nodeName, this.requireDownloadForNodes([this.nodeName]), this.form.skip_downloads)
        }
      }
      if (this.requireSeparateDownloadAction) {
        this.form.action = 'download_binaries'
      }
    },
    testPingPong (nodes) {
      this.pingpong = {}
      this.pingpong_loading = true
      let req = { nodes: nodes }
      this.kuboardSprayApi.post(`/clusters/${this.cluster.name}/state/ping`, req).then(resp => {
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
    async execute () {
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
            { // 排除节点
              let temp = ''
              let excludes = {}
              for (let node of this.offlineNodes) {
                excludes[node] = true
              }
              for (let node in excludes) {
                temp += '!' + node + ','
              }
              req.nodes_to_exclude = trimMark(temp)
            }
            if (this.form.action === 'upgrade_all_nodes') {
              req.nodes = 'target'
              req.skip_downloads = true
            } else if (this.form.action === 'upgrade_master_nodes') {
              req.nodes = 'kube_control_plane,etcd'
              req.skip_downloads = true
            } else if (this.form.action === 'upgrade_multi_nodes') {
              let temp = ''
              for (let item of this.form.kube_nodes_to_upgrade) {
                temp += item + ','
              }
              temp = trimMark(temp)
              req.nodes = temp
              req.skip_downloads = this.form.skip_downloads
            } else if (this.form.action === 'upgrade_single_node') {
              req.nodes = this.nodeName
              req.skip_downloads = this.form.skip_downloads
            }
            this.kuboardSprayApi.post(`/clusters/${this.cluster.name}/${this.form.action}`, req).then(resp => {
              let pid = resp.data.data.pid
              resolve(pid)
            }).catch(e => {
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

<style scoped lang="css">
.form_description {
  font-size: 12px;
  color: #aaa;
  max-width: 700px;
}
</style>
