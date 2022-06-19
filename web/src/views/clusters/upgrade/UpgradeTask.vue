<i18n>
en:
  operation: Operations
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

  uncordon_node: Uncordon node {nodeName}
  uncordon_node_desc: Uncordon node {nodeName} to bring it back.
  drain_node: Drain node {nodeName}
  drain_node_desc: Drain node before upgrading it to avoid business interruption.
  drain_timeout: drain_timeout
  drain_retries: drain_retries
  drain_grace_period: drain_grace_period
  drain_retry_delay_seconds: drain_retry_delay_seconds
  nodes_to_remove_required: Please select nodes to remove.

  newResourcePackageRequired: This resource package does not support the operation, please choose a new one.
zh:
  operation: 操作选项
  download_binaries: 分发安装包
  download_binaries_desc: 升级过程中，将会改变一些 daemonset 的版本（例如：kube-proxy, calico-node 等），如果不事先分发这些镜像到所有各个节点，将导致升级过程中，部分节点可能因为缺少对应版本的镜像而出现问题。

  skip_downloads: 跳过加载安装包

  upgrade_cluster: 升级集群
  upgrade_multi_nodes: 升级多个工作节点
  upgrade_multi_nodes_desc: 请选择本次任务中要升级的节点
  upgrade_single_node: 升级 {nodeName}
  upgrade_single_node_desc: 一次只升级一个节点。以下是当前节点上的容器组列表，如果存在非 DaemonSet 类型的容器组，请考虑是否先排空该节点。

  upgrade_all_nodes: 升级所有节点
  upgrade_all_nodes_desc: 一次性升级所有节点，如果当前有工作负载正在运行，将造成服务中断
  upgrade_master_nodes: 升级控制节点和 ETCD 节点
  upgrade_master_nodes_desc: 先逐个升级控制节点和 ETCD 节点，然后再逐个排空、升级工作节点，避免造成服务中断。（如果有3个控制节点，大约耗时 18 分钟）

  uncordon_node: 恢复调度 {nodeName}
  uncordon_node_desc: 恢复调度 {nodeName}，以便新的容器组可以调度到该节点上。
  drain_node: 排空节点 {nodeName}
  drain_node_desc: 建议在升级节点前，排空节点上的容器组，以避免出现业务中断
  drain_timeout: 排空节点超时时间
  drain_retries: 排空节点重试次数
  drain_grace_period: 应用停止等候时间
  drain_retry_delay_seconds: 两次排空尝试间隔
  nodes_to_remove_required: 请选择要删除的节点

  newResourcePackageRequired: 当前资源包不支持此操作，请使用新的资源包
</i18n>

<template>
  <ClusterTask :history="cluster.history" :populateRequest="populateRequest" @onShow="onShow" :action="form.action" playbook_check="upgrade_cluster"
    :cluster="cluster" :title="title" @refresh="$emit('refresh')">
    <el-form ref="form" :model="form" style="width: 850px;" label-width="120px">
      <el-form-item :label="$t('operation')" prop="action" style="margin-top: 10px;">
        <el-radio-group v-model="form.action" class="app_margin_bottom">
          <template v-if="nodeName === undefined">
            <el-radio-button label="download_binaries">
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
          <template v-else>
            <el-radio-button label="drain_node" v-if="cluster.resourcePackage.data.supported_playbooks.drain_node"
              :disabled="cluster.state.nodes[nodeName].status.nodeInfo.kubeletVersion === cluster.resourcePackage.data.kubernetes.kube_version">
              {{$t('drain_node', {nodeName})}}
            </el-radio-button>
            <el-radio-button label="upgrade_single_node" :disabled="cluster.state.nodes[nodeName].status.nodeInfo.kubeletVersion === cluster.resourcePackage.data.kubernetes.kube_version">
              {{$t('upgrade_single_node', {nodeName})}}
            </el-radio-button>
            <el-radio-button label="uncordon_node" v-if="cluster.state.nodes[nodeName] && cluster.resourcePackage.data.supported_playbooks.uncordon_node"
              :disabled="!cluster.state.nodes[nodeName].spec.unschedulable">
              {{$t('uncordon_node', {nodeName})}}
            </el-radio-button>
          </template>
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
        <div v-if="form.action === 'drain_node'">
          <pre class="drain_cmd app_text_mono">kubectl cordon {{ nodeName }}
kubectl drain --force --ignore-daemonsets --grace-period {{ form.drain_node.drain_grace_period * 60}} --timeout {{ form.drain_node.drain_timeout * 60 }}s</pre>
          <el-form-item label-width="150px" :label="$t('drain_grace_period')" required prop="drain_node.drain_grace_period">
            <el-input-number v-model="form.drain_node.drain_grace_period" :min="1" :max="form.drain_node.drain_timeout - 1"></el-input-number> 分钟
            <span style="margin-left: 20px;" class="form_description">kubectl drain --grace-period</span>
          </el-form-item>
          <el-form-item label-width="150px" :label="$t('drain_timeout')" required prop="drain_node.drain_timeout">
            <el-input-number v-model="form.drain_node.drain_timeout" :min="form.drain_node.drain_grace_period + 1"></el-input-number> 分钟
            <span style="margin-left: 20px;" class="form_description">kubectl drain --timeout</span>
          </el-form-item>
          <el-form-item label-width="150px" :label="$t('drain_retries')" required prop="drain_node.drain_retries">
            <el-input-number v-model="form.drain_node.drain_retries" :min="1"></el-input-number> 次
          </el-form-item>
          <el-form-item label-width="150px" :label="$t('drain_retry_delay_seconds')" required prop="drain_node.drain_retry_delay_seconds">
            <el-input-number v-model="form.drain_node.drain_retry_delay_seconds" :min="5" :step="5"></el-input-number> 秒
          </el-form-item>
        </div>
        <el-scrollbar v-if="form.action === 'upgrade_single_node'" :max-height="240" style="max-height: 240px; margin-bottom: 10px;">
          <pre class="pods_on_node">{{pods_on_node}}</pre>
        </el-scrollbar>
        <div v-if="form.action === 'upgrade_multi_nodes' || form.action === 'upgrade_single_node'">
          <el-switch v-model="form.skip_downloads" :disabled="requireDownloadForNodes(form.kube_nodes_to_upgrade)" :active-text="$t('skip_downloads')"></el-switch>
        </div>
      </el-form-item>
    </el-form>
  </ClusterTask>
</template>

<script>
import ClusterTask from '../../common/task/ClusterTask.vue'

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
        action: 'upgrade_master_nodes',
        kube_nodes_to_upgrade: [],
        skip_downloads: false,
        drain_node: {
        },
      },
      kube_nodes_to_upgrade_rules: [
        { required: true, message: this.$t('upgrade_multi_nodes_desc') }
      ],
      pods_on_node: undefined,
    }
  },
  computed: {
    title () {
      if (this.nodeName) {
        return this.$t(this.form.action, { nodeName: this.nodeName })
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
  },
  components: { ClusterTask },
  emits: ['refresh'],
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
    onShow() {
      this.setAction()
      this.form.kube_nodes_to_upgrade = []
      if (this.nodeName) {
        this.getPodsOnNode()
      }
    },
    setAction () {
      if (!this.controlPlanePendingUpgrade) {
        if (this.nodeName === undefined) {
          this.form.action = 'upgrade_multi_nodes'
          this.form.skip_downloads = false
        } else {
          this.form.action = 'drain_node'
          this.form.skip_downloads = !this.requireDownloadForNodes([this.nodeName])
          if (this.cluster.state.nodes[this.nodeName].spec.unschedulable) {
            this.form.action = 'upgrade_single_node'
          }
          if (this.cluster.state.nodes[this.nodeName].status.nodeInfo.kubeletVersion === this.cluster.resourcePackage.data.kubernetes.kube_version) {
            this.form.action = 'uncordon_node'
          }
          this.form.drain_node = {}
          this.form.drain_node.drain_grace_period = 5
          this.form.drain_node.drain_timeout = 6
          this.form.drain_node.drain_retries = 2
          this.form.drain_node.drain_retry_delay_seconds = 10
        }
      }
      if (this.requireSeparateDownloadAction) {
        this.form.action = 'download_binaries'
      }
    },
    getPodsOnNode () {
      this.pods_on_node = undefined
      this.kuboardSprayApi.get(`/clusters/${this.cluster.name}/state/pods_on_node/${this.nodeName}`).then(resp => {
        this.pods_on_node = resp.data.data.stdout
      })
    },
    populateRequest () {
      console.log('populateRequest')
      let _this = this
      return new Promise((resolve, reject) => {
        _this.$refs.form.validate(flag => {
          if (!flag) {
            return reject ()
          }
          let req = {}
          if (_this.form.action === 'upgrade_all_nodes') {
            req.nodes = 'target'
            req.skip_downloads = true
          } else if (_this.form.action === 'upgrade_master_nodes') {
            req.nodes = 'kube_control_plane,etcd'
            req.skip_downloads = true
          } else if (_this.form.action === 'upgrade_multi_nodes') {
            let temp = ''
            for (let item of _this.form.kube_nodes_to_upgrade) {
              temp += item + ','
            }
            temp = trimMark(temp)
            req.nodes = temp
            req.skip_downloads = _this.form.skip_downloads
          } else if (_this.form.action === 'upgrade_single_node') {
            req.nodes = _this.nodeName
            req.skip_downloads = _this.form.skip_downloads
          } else if (_this.form.action === 'drain_node') {
            req.nodes = _this.nodeName
            req.drain_grace_period = _this.form.drain_node.drain_grace_period + ''
            req.drain_timeout = _this.form.drain_node.drain_timeout + 's'
            req.drain_retries = _this.form.drain_node.drain_retries + ''
            req.drain_retry_delay_seconds = _this.form.drain_node.drain_retry_delay_seconds + ''
          } else if (_this.form.action === 'uncordon_node') {
            req.nodes = _this.nodeName
          }
          resolve(req)
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
.pods_on_node {
  background-color: var(--el-text-color-primary);
  color: var(--el-color-white);
  padding: 20px;
  margin: 0;
  min-height: 200px;
}
.drain_cmd {
  background-color: var(--el-text-color-primary);
  color: var(--el-color-white);
  padding: 10px 20px;
  font-weight: bold;
}
</style>
