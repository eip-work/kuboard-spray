<i18n>
en:
  is_installed_true: Not Installed
  is_installed_false: Installed
  install_addon: Install addon
  remove_addon: Remove addon
  install: Install
  uninstall: Uninstall

  verbose: Include task params
  verbose_: 
  verbose_v: May include sensitive data in the trace, e.g. path to files, user name, password.
  verbose_vvv: includes more information in log, only used in development.
  v_: Normal
  v_v: Details
  v_vvv: More
  fork: ansible fork
  fork_more: Max number of nodes can be operated in the installation.
zh:
  is_installed_true: 已安装
  is_installed_false: 未安装
  install_addon: 安装可选组件
  remove_addon: 卸载可选组件
  install: 安 装
  uninstall: 卸 载

  verbose: 显示任务参数
  verbose_: 正常输出的日志，通常选用此选项
  verbose_v: 日志中会包含部分敏感信息，例如：文件路径、用户名密码等
  verbose_vvv: 日志中会包含最详细的信息，通常只在开发阶段使用
  v_: 正常
  v_v: 详细
  v_vvv: 更多
  fork: 并发数量
  fork_more: 安装过程中可以同时操作的目标节点的最大数量。ansible fork.
</i18n>

<template>
  <ConfigSection v-model:enabled="enabledRef" :label="label" :helpLink="`https://kuboard-spray.cn/guide/addons/${addonName}.html`" :anti-freeze="computedAntiFreeze">
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
          <el-form ref="form" :model="form" @submit.prevent.stop label-position="left" label-width="120px">
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
          </el-form>
        </ExecuteTask>
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

export default {
  props: {
    enabled: { prop: Boolean, required: true },
    label: { prop: String, require: false, default: "" },
    cluster: { prop: Object, required: true },
    addonName: { prop: String, required: true },
    antiFreeze: { prop: Boolean, required: false, default: undefined },
    pingpong: { prop: Object, required: false, default: () => { return {} }},
  },
  data() {
    return {
      form: {
        verbose: '',
        fork: 5,
        action: 'install_addon',
      },
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
      if (this.cluster && this.cluster.state && this.addonState.is_installed) {
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
    onVisibleChange (v) {
      if (v && this.cluster.state && this.addonState) {
        this.action = this.addonState.is_installed ? 'remove_addon' : 'install_addon'
      }
    },
    async applyPlan () {
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

<style scoped lang="scss">

</style>
