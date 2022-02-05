<i18n>
en:
  connectivity_check: Network connectivity check
  connectivity_check_desc: Check if nodeport, hostport, containerport and service works in kubernetes cluster.
  selectAPod: Please select a pod first.
  pass: Passed network connectivity check
  fail: Failed network connectivity check
  netcheckerNotInstalled: "You didn't install netchecker."
  installNetchecker: Install netchecker
  podList: Pod list
  podDetails: Pod details
  events: Events
  logs: Logs
zh:
  connectivity_check: 网络连通性检查
  connectivity_check_desc: 检查访问节点端口、访问 hostPort、从容器访问其他容器、从容器访问 Service 等 K8S 各组件间的网络连通性。
  selectAPod: 请先选择一个容器组。
  pass: 网络连通性检查通过
  fail: 网络连通性检查不通过
  netcheckerNotInstalled: 未安装 “网络检查” 组件，请到集群规划页
  installNetchecker: 安装 “网络检查” 组件
  podList: 容器组列表
  podDetails: 容器组详情
  events: 事 件
  logs: 日 志
</i18n>

<template>
  <div>
    <div class="app_block_title">{{ $t('connectivity_check') }}</div>
    <div class="app_description">
      {{ $t('connectivity_check_desc') }}
      <KuboardSprayLink href="https://kuboard-spray.cn/guide/addons/netchecker.html" style="margin-left: 20px;" :size="13"></KuboardSprayLink>
    </div>
    <div v-if="installed">
      <el-skeleton v-if="loading" animated class="app_margin_top"></el-skeleton>
      <div v-else-if="connectivityCheck">
        <el-alert :type="status === 'pass' ? 'success' : 'error'" effect="dark" :closable="false" :title="$t(status)" show-icon>
          <span class="app_text_mono">{{ errorMsg || connectivityCheck.stdout_obj.Message }}</span>
        </el-alert>
        <div v-if="status === 'fail'" style="display: flex;" class="app_description">
          <el-scrollbar style="margin-bottom: 10px; min-width: 285px;" height="450px">
            <div v-for="(item, index) in connectivityCheck.stdout_obj.Absent" :key="'absent' + index" style="margin-bottom: 10px;">
              <el-button type="danger" style="width: 280px;" icon="el-icon-check" :plain="connectivityCheckLog.pod !== item"
                @click="connectivityCheckLog.pod = item" :loading="connectivityCheckLog.pod === item && connectivityLogsLoading">
                <span class="app_text_mono">{{ item }}</span>
              </el-button>
            </div>
            <div v-for="(item, index) in connectivityCheck.stdout_obj.Outdated" :key="'outdated' + index" style="margin-bottom: 10px;">
              <el-button type="danger" style="width: 280px;" icon="el-icon-clock" :plain="connectivityCheckLog.pod !== item"
                @click="connectivityCheckLog.pod = item" :loading="connectivityCheckLog.pod === item && connectivityLogsLoading">
                <span class="app_text_mono">{{ item }}</span>
              </el-button>
            </div>
          </el-scrollbar>
          <div style="flex-grow: 1; margin-left: 10px; margin-bottom: 38px;">
            <el-radio-group v-model="connectivityCheckLog.type" size="default">
              <el-radio-button label="podList">
                <el-icon style="vertical-align: top; margin-right: 5px;">
                  <component :is="connectivityCheckLog.type === 'podList' && connectivityLogsLoading ? 'el-icon-loading' : 'el-icon-copy-document'"></component>
                </el-icon>
                {{ $t('podList') }}
              </el-radio-button>
              <el-radio-button label="podDetails">
                <el-icon style="vertical-align: top; margin-right: 5px;">
                  <component :is="connectivityCheckLog.type === 'podList' && connectivityLogsLoading ? 'el-icon-loading' : 'el-icon-document-copy'"></component>
                </el-icon>
                {{ $t('podDetails') }}
              </el-radio-button>
              <el-radio-button label="events">
                <el-icon style="vertical-align: top; margin-right: 5px;">
                  <component :is="connectivityCheckLog.type === 'podList' && connectivityLogsLoading ? 'el-icon-loading' : 'el-icon-files'"></component>
                </el-icon>
                {{ $t('events') }}
              </el-radio-button>
              <el-radio-button label="logs">
                <el-icon style="vertical-align: top; margin-right: 5px;">
                  <component :is="connectivityCheckLog.type === 'podList' && connectivityLogsLoading ? 'el-icon-loading' : 'el-icon-document'"></component>
                </el-icon>
                {{ $t('logs') }}
              </el-radio-button>
            </el-radio-group>
            <Codemirror class="connectivity_check_code_mirror" v-model:value="connectivityLogsCompute" :options="cmOptions"></Codemirror>
          </div>
        </div>
      </div>
      <div v-if="errorMsg">
        <el-alert :title="'' + errorMsg.response.status" :closable="false" type="error" effect="dark">
          <div class="app_text_mono">
            <div v-if="errorMsg.response.status === 500">
              {{ errorMsg.response.data.message }}
            </div>
            <div v-else>
              {{ errorMsg }}
            </div>
          </div>
        </el-alert>
      </div>
    </div>
    <el-alert v-else type="warning" :closable="false" :title="$t('installNetchecker')">
      {{ $t('netcheckerNotInstalled') }}
      <KuboardSprayLink href="https://kuboard-spray.cn/guide/addons/netchecker.html" style="margin-left: 20px;" :size="12">{{ $t('installNetchecker') }}</KuboardSprayLink>
    </el-alert>
  </div>
</template>

<script>
import Codemirror from "codemirror-editor-vue3"
import "codemirror/theme/darcula.css"
import "codemirror/mode/yaml/yaml.js"
import "codemirror/mode/shell/shell.js"

export default {
  props: {
    cluster: { type: Object, required: true },
  },
  data() {
    return {
      errorMsg: undefined,
      loading: false,
      connectivityCheck: undefined,
      connectivityLogs: this.$t('selectAPod'),
      connectivityLogsLoading: false,
      connectivityCheckLog: {
        type: 'events',
        pod: '',
      },
      cmOptions: {
        mode: "shell",
        theme: "darcula",
        lineNumbers: true,
        lineWrapping: true,
        smartIndent: true,
        indentUnit: 2,
        foldGutter: true,
        styleActiveLine: true,
      }
    }
  },
  computed: {
    installed () {
      return this.cluster.state && this.cluster.state.addons.netchecker && this.cluster.state.addons.netchecker.is_installed
    },
    status () {
      if (this.connectivityCheck && this.connectivityCheck.stdout_obj) {
        if (!this.connectivityCheck.stdout_obj.Absent && !this.connectivityCheck.stdout_obj.Outdated) {
          return 'pass'
        }
      }
      return 'fail'
    },
    connectivityLogsCompute: {
      get () {
        return this.connectivityLogs + ''
      },
      set () {}
    }
  },
  components: { Codemirror },
  mounted () {
    this.refresh()
  },
  watch: {
    connectivityCheckLog: {
      handler: function (newValue) {
        if (newValue.pod || newValue.type === 'podList') {
          this.connectivityLogs = 'loading ...'
          this.connectivityLogsLoading = true
          let namespace = this.cluster.inventory.all.children.target.children.k8s_cluster.vars.netcheck_namespace || 'default'
          this.kuboardSprayApi.get(`/clusters/${this.cluster.name}/health_check/details?namespace=${namespace}&&pod=${newValue.pod}&&type=${newValue.type}`).then(resp => {
            console.log(resp.data.data)
            this.connectivityLogs = resp.data.data.stdout || resp.data.data.stderr
            this.connectivityLogsLoading = false
          }).catch(e => {
            this.connectivityLogsLoading = false
            if (e.response.status === 500) {
              this.connectivityLogs = e.response.data.message
            } else if (e.response.status === 400){
              this.connectivityLogs = e.response.data
            }
          })
        } else {
          this.connectivityLogs = this.$t('selectAPod')
        }
      },
      deep: true,
    }
  },
  methods: {
    async refresh () {
      if (!this.installed) {
        return
      }
      this.loading = true
      this.errorMsg = undefined
      await this.kuboardSprayApi.get(`/clusters/${this.cluster.name}/health_check/connectivity_check`).then(resp => {
        console.log(resp.data)
        resp.data.data.stdout_obj = JSON.parse(resp.data.data.stdout)
        this.connectivityCheck = resp.data.data
      }).catch(e => {
        this.errorMsg = e
      })
      this.loading = false
    }
  }
}
</script>

<style lang="css">
.connectivity_check_code_mirror .CodeMirror {
  height: 380px;
}

.connectivity_check_code_mirror .CodeMirror-scroll {
  height: 380px;
  overflow-y: hidden;
  overflow-x: auto;
}
</style>
