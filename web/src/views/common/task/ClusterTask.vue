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

</i18n>

<template>
  <ExecuteTask :history="cluster.history" :loading="loading" :title="title" :startTask="execute" @refresh="$emit('refresh')"
    :disabled="disabled" :type="type" @visibleChange="onVisibleChange">
    <div :style="`width: ${width}px;`">
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
        
        <div style="font-size: 12px;">
          <slot v-bind:form="form"></slot>
        </div>

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

        <el-form-item v-if="action && cluster && cluster.resourcePackage && cluster.resourcePackage.data.supported_playbooks[playbook_check || action] === undefined" prop="min_resource_package_version" 
          style="margin-top: -10px;"
          :rules="min_resource_package_version_rules">
        </el-form-item>
      </el-form>
    </div>
  </ExecuteTask>
</template>

<script>
import ExecuteTask from './ExecuteTask.vue'

function trimMark(str) {
  if (str[str.length - 1] === ',') {
    return str.slice(0, str.length - 1)
  }
  return str
}

export default {
  props: {
    cluster: { type: Object, required: true },
    action: { type: String, required: false, default: undefined },
    playbook_check: { type: String, required: false, default: undefined },
    title: { type: String, required: true },
    populateRequest: { type: Function, required: true },
    width: { type: Number, required: false, default: 850 },
    type: { type: String, required: false, default: 'warning' },
    disabled: { type: Boolean, required: false, default: false },
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
      },
      min_resource_package_version_rules: [
        { required: true, message: this.$t('newResourcePackageRequired') }
      ],
      pingpong: {},
      pingpong_loading: true,
      pods_on_node: undefined,
    }
  },
  computed: {
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
  },
  emits: ['onShow', 'refresh'],
  methods: {
    onVisibleChange(flag) {
      if (flag) {
        let count = 0
        for (let key in this.cluster.inventory.all.hosts) {
          if (key !== 'localhost')
          count ++
        }
        this.form.fork = count
        this.testPingPong('target')
        this.$emit('onShow')
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
      let _this = this
      return new Promise((resolve, reject) => {
        _this.$refs.form.validate(flag => {
          if (flag) {
            let req = {
              verbose: _this.form.verbose,
              fork: _this.form.fork,
            }
            { // 排除节点
              let temp = ''
              let excludes = {}
              for (let node of _this.offlineNodes) {
                excludes[node] = true
              }
              for (let node in excludes) {
                temp += '!' + node + ','
              }
              req.nodes_to_exclude = trimMark(temp)
            }

            _this.populateRequest.apply(_this.$parent, _this.form).then(r => {
              req = Object.assign(req, r)
              _this.kuboardSprayApi.post(`/clusters/${_this.cluster.name}/${_this.action}`, req).then(resp => {
                let pid = resp.data.data.pid
                resolve(pid)
              }).catch(e => {
                reject(e.response.data.message)
              })
            }).catch(e => {
              console.error(e)
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
