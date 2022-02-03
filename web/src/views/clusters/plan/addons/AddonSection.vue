<i18n>
en:
  is_installed_true: Not Installed
  is_installed_false: Installed
  install_addon: Install addon
  remove_addon: Remove addon
  install: Install
  uninstall: Uninstall

  offlineNodes: Offline nodes
  offlineNodesDesc: Exclude the following offline nodes.

  verbose: Include task params
  verbose_: 
  verbose_v: May include sensitive data in the trace, e.g. path to files, user name, password.
  verbose_vvv: includes more information in log, only used in development.
  v_: Normal
  v_v: Details
  v_vvv: More
  fork: ansible fork
  fork_more: Max number of nodes can be operated in the installation.

  saveAndInstall: Please fill in params for the addon, save and then execute install task.
zh:
  is_installed_true: 已安装
  is_installed_false: 未安装
  install_addon: 安装可选组件
  remove_addon: 卸载可选组件
  install: 安 装
  uninstall: 卸 载

  offlineNodes: 离线节点
  offlineNodesDesc: 排除以下不在线的节点：

  verbose: 显示任务参数
  verbose_: 正常输出的日志，通常选用此选项
  verbose_v: 日志中会包含部分敏感信息，例如：文件路径、用户名密码等
  verbose_vvv: 日志中会包含最详细的信息，通常只在开发阶段使用
  v_: 正常
  v_v: 详细
  v_vvv: 更多
  fork: 并发数量
  fork_more: 安装过程中可以同时操作的目标节点的最大数量。ansible fork.

  saveAndInstall: 请先填写该可选组件所需的参数，保存后再执行安装动作
</i18n>

<template>
  <ConfigSection ref="configSection" v-model:enabled="enabledRef" :label="label" :helpLink="`https://kuboard-spray.cn/guide/addons/${addonName}.html`" :anti-freeze="computedAntiFreeze">
    <template #header>
      {{ label }}
      <div v-if="cluster && addonState" style="float: right; margin-right: 10px; margin-top: -1px;">
        <el-tag v-if="addonState.is_installed" type="success" style="margin-right: 10px;">
          <i class="el-icon-check" style="margin-right: 5px;"></i>
          {{$t('is_installed_true')}}
        </el-tag>
        <el-tag v-else type="warning" style="margin-right: 10px;">
          <i class="el-icon-close" style="margin-right: 5px;"></i>
          {{$t('is_installed_false')}}
        </el-tag>
        <ExecuteTask v-if="showExecute" :history="cluster.history" placement="left"
          :disabled="executeDisabled"
          :title="$t(action, {label})" :startTask="applyPlan" @refresh="$emit('refresh')" @visibleChange="onVisibleChange">
          <template #reference>
            <el-button v-if="addonState.is_installed" icon="el-icon-remove" type="danger" plain>{{$t('uninstall')}}</el-button>
            <el-button v-else icon="el-icon-download" type="warning" plain>{{$t('install')}}</el-button>
          </template>
          <div v-if="pingpong_loading" style="display: block; min-width: 420px;">
            <el-skeleton animated></el-skeleton>
          </div>
          <el-form v-else ref="form" :model="form" @submit.prevent.stop label-position="left" label-width="120px">
            <div class="app_block_title">{{label}}</div>
            <el-form-item :label="$t('verbose')">
              <el-radio-group v-model="form.verbose">
                <el-radio-button label="">{{$t('v_')}}</el-radio-button>
                <el-radio-button label="v">{{$t('v_v')}}</el-radio-button>
                <el-radio-button label="vvv">{{$t('v_vvv')}}</el-radio-button>
              </el-radio-group>
              <div style="width: 350px; margin-left: 0; color: #aaa; font-size: 12px;">{{$t('verbose_' + form.verbose)}}</div>
            </el-form-item>
            <el-form-item :label="$t('fork')" style="margin-top: 10px;">
              <el-input-number v-model="form.fork" :step="2" style="width: 166px;"></el-input-number>
              <div style="width: 350px; margin-left: 0; color: #aaa; font-size: 12px;">{{$t('fork_more')}}</div>
            </el-form-item>
            <el-form-item :label="$t('offlineNodes')" class="app_margin_top" v-if="offlineNodes.length > 0">
            <div class="form_description">{{ $t('offlineNodesDesc') }}</div>
            <template v-for="node in offlineNodes" :key="'exclude' + node">
              <el-tooltip class="box-item" effect="dark" :content="pingpong[node].message" placement="top-end">
                <el-tag type="danger" effect="dark" style="margin: 0 10px 10px 0;">
                  <span class="app_text_mono" style="font-size: 14px; margin-right: 10px;">{{ node }}</span> <i class="el-icon-question"></i>
                </el-tag>
              </el-tooltip>
            </template>
          </el-form-item>
          </el-form>
        </ExecuteTask>
        <el-button v-else-if="isClusterOnline && editMode === 'view' && !addonState.is_installed" type="primary" plain icon="el-icon-download" @click="prepareForInstall">{{ $t('install_addon') }}</el-button>
      </div>
    </template>
    <template #more>
      <slot name="more"></slot>
    </template>
    <slot></slot>
  </ConfigSection>
</template>

<script>
import { computed } from 'vue'
import ExecuteTask from '../../../common/task/ExecuteTask.vue'

function trimMark(str) {
  if (str[str.length - 1] === ',') {
    return str.slice(0, str.length - 1)
  }
  return str
}

export default {
  props: {
    enabled: { prop: Boolean, required: true },
    label: { prop: String, require: false, default: "" },
    cluster: { prop: Object, required: true },
    addonName: { prop: String, required: true },
    antiFreeze: { prop: Boolean, required: false, default: undefined },
    // pingpong: { prop: Object, required: false, default: () => { return {} }}, 
  },
  data() {
    return {
      form: {
        verbose: '',
        fork: 5,
        action: 'install_addon',
      },
      pingpong: {},
      pingpong_loading: true,
    }
  },
  inject: ['editMode', 'isClusterOnline', 'isClusterInstalled'],
  provide () {
    return {
      // 根据可选组件的安装状态决定是否冻结组件属性的编辑
      editMode: computed(() => {
        if (this.editMode == 'frozen') {
          return this.computedAntiFreeze ? 'edit' : 'frozen'
        }
        return this.editMode
      })
    }
  },
  computed: {
    offlineNodes () {
      let result = []
      for (let key in this.cluster.inventory.all.hosts) {
        if (this.pingpong[key] && this.pingpong[key].status !== 'SUCCESS') {
          result.push(key)
        }
      }
      return result
    },
    addonInResourcePackage () {
      if (this.cluster.resourcePackage) {
        for (let addon of this.cluster.resourcePackage.data.addon) {
          if (addon.name === this.addonName) {
            return addon
          }
        }
      }
      return {}
    },
    addonState () {
      if (this.cluster.state) {
        return this.cluster.state.addons[this.addonName]
      }
      return undefined
    },
    showExecute () {
      if (this.cluster && this.cluster.history && this.cluster.history.processing) {
        return false
      }
      if (this.editMode === 'view' && this.isClusterInstalled && this.isClusterOnline) {
        return this.addonState.intend_to_install
      }
      return false
    },
    executeDisabled () {
      return this.cluster.resourcePackage.data.supported_playbooks[this.action] === undefined 
        || this.addonInResourcePackage.lifecycle == undefined
        || this.addonInResourcePackage.lifecycle[this.action + '_tags'] === undefined
    },
    enabledRef: {
      get () { return this.enabled },
      set (v) { this.$emit('update:enabled', v) }
    },
    computedAntiFreeze () {
      if (this.antiFreeze !== undefined) {
        return this.antiFreeze
      }
      if (this.cluster && this.cluster.state && this.addonState && this.addonState.is_installed) {
        return false
      }
      return true
    },
    action () {
      if (this.cluster.state && this.addonState) {
        return this.addonState.is_installed ? 'remove_addon' : 'install_addon'
      }
      return ''
    }
  },
  components: { ExecuteTask },
  mounted () {
    this.onVisibleChange(true)
  },
  methods: {
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
        this.pingpong_loading = false
      })
    },
    prepareForInstall () {
      this.enabledRef = true
      this.$router.replace('?mode=edit')
      setTimeout(() => {
        this.$refs.configSection.expand(true)
        this.$message.warning(this.$t('saveAndInstall'))
      }, 500)
    },
    onVisibleChange (v) {
      if (v && this.cluster.state && this.addonState) {
        this.action = this.addonState.is_installed ? 'remove_addon' : 'install_addon'
        this.testPingPong('target')
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
              addon_name: this.addonName,
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
              for (let node in excludes) {
                temp += '!' + node + ','
              }
              req.nodes_to_exclude = trimMark(temp)
            }
            console.log(req)
            this.kuboardSprayApi.post(`/clusters/${this.cluster.name}/${this.action}`, req).then(resp => {
              let pid = resp.data.data.pid
              resolve(pid)
            }).catch(e => {
              reject(e.response.data.message || e)
            })
          } else {
            reject('请验证表单')
          }
        })
      })
    }
  }
}
</script>

<style scoped lang="css">

</style>
