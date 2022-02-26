<i18n>
en:
  verbose: Include task params
  verbose_: 
  verbose_v: May include sensitive data in the trace, e.g. path to files, user name, password.
  verbose_vvv: includes more information in log, only used in development.
  v_: Normal
  v_v: Details
  v_vvv: More
  control_params: Control params
  operation: Operations
  offlineNodes: Offline nodes
  offlineNodesDesc: Exclude the following offline nodes.

  cert_renew: Renew Certificate
zh:
  verbose: 显示任务参数
  verbose_: 正常输出的日志，通常选用此选项
  verbose_v: 日志中会包含部分敏感信息，例如：文件路径、用户名密码等
  verbose_vvv: 日志中会包含最详细的信息，通常只在开发阶段使用
  v_: 正常
  v_v: 详细
  v_vvv: 更多
  control_params: 控制选项
  operation: 操作选项
  offlineNodes: 离线节点
  offlineNodesDesc: 排除以下不在线的节点：

  cert_renew: 更新 APIServer 证书
</i18n>

<template>
  <ExecuteTask :history="cluster.history" :loading="loading" :title="title" :startTask="execute" @refresh="$emit('refresh')" @visibleChange="onVisibleChange">
    <div style="width: 450px;">
      <div v-if="pingpong_loading" style="display: block;">
        <el-skeleton animated></el-skeleton>
      </div>
      <el-form v-else ref="form" :model="form" @submit.prevent.stop label-position="left" label-width="120px" class="app_form_mini">
        <el-form-item :label="$t('verbose')">
          <el-radio-group v-model="form.verbose">
            <el-radio-button label="">{{$t('v_')}}</el-radio-button>
            <el-radio-button label="v">{{$t('v_v')}}</el-radio-button>
            <el-radio-button label="vvv">{{$t('v_vvv')}}</el-radio-button>
          </el-radio-group>
          <div style="color: #aaa; font-size: 12px;">{{$t('verbose_' + form.verbose)}}</div>
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
        drain_node: {

        },
      },
      min_resource_package_version_rules: [
        { required: true, message: this.$t('newResourcePackageRequired') }
      ],
      kube_nodes_to_upgrade_rules: [
        { required: true, message: this.$t('upgrade_multi_nodes_desc') }
      ],
      pingpong: {},
      pingpong_loading: true,
      pods_on_node: undefined,
    }
  },
  computed: {
    title () {
      return this.$t('cert_renew')
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
        if (this.nodeName) {
          this.getPodsOnNode()
        }
      }
    },
    setAction () {

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
    getPodsOnNode () {
      this.pods_on_node = undefined
      this.kuboardSprayApi.get(`/clusters/${this.cluster.name}/state/pods_on_node/${this.nodeName}`).then(resp => {
        this.pods_on_node = resp.data.data.stdout
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
              fork: 1,
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
            this.kuboardSprayApi.post(`/clusters/${this.cluster.name}/renew_cert`, req).then(resp => {
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
