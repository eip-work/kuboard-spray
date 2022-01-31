<i18n>
en:
  "terminal": "terminal"
  "find": "Find"
  "connecting": "Connecting"
  "connected": "Connected"
  "closing": "Closing"
  "closed": "Closed"
  "unknown": "Unknown"
zh:
  "terminal": "终端"
  "find": "查 找"
  "connecting": "正在连接"
  "connected": "已连接"
  "closing": "正在关闭"
  "closed": "已断开"
  "unknown": "未知状态"
</i18n>

<template>
  <div>
    <div style="background-color: #3a3333; color: #fffeff; font-weight: 500; padding: 8px 20px 8px 20px; line-height: 22px; text-align: center;">
      {{this.$t('terminal')}} - <span class="app_text_mono">{{ownerType}}/{{ownerName}}/{{nodeName}}</span>
      <span style="float:left;">
        <ChangeFontSize class="button" :terminal="xterm" :fitAddon="fitAddon"></ChangeFontSize>
        <ChangeColor class="button"></ChangeColor>
        <el-button type="info" icon="el-icon-search" @click="$refs.find.show()">{{$t('find')}}</el-button>
        <Find ref="find" class="button" :terminal="xterm"></Find>
      </span>
      <span :style="`float: right; font-size: 14px; font-weight: 600; color: ${socketReadyState === 1 ? '#33FF33' : '#FF6600'};`">
        {{ stateStr }}
      </span>
    </div>
    <div id="terminal" :style="`width: 100%; height: calc(100vh - 67px); background-color: black;`"></div>
    <SshQuickCommands :roles="roles" :socket="socket" :cluster="cluster"></SshQuickCommands>
    <K8sTerminalErrorHint ref="errorHint"/>
  </div>
</template>

<script>
import { Terminal } from 'xterm'
import { FitAddon } from 'xterm-addon-fit'
import 'xterm/css/xterm.css'
import K8sTerminalErrorHint from './K8sTerminalErrorHint'
import ChangeFontSize from './ChangeFontSize'
import Find from './Find'
import ChangeColor from './ChangeColor'
import SshQuickCommands from './SshQuickCommands.vue'
// import {encode} from 'js-base64'

function trimSlash(str) {
  if (str[str.length - 1] === '/') {
    return str.slice(0, str.length - 1)
  }
  return str
}

export default {
  props: {
    ownerType: { type: String, required: true },
    ownerName: { type: String, required: true },
    nodeName: { type: String, required: true },
  },
  data() {
    return {
      socketReadyState: 0,
      ttySize: {rows: 100},
      xterm: undefined,
      fitAddon: undefined,
      fontSize: parseInt(localStorage.getItem('terminal-font-size')) || 14,
      isRunning: true,
      dialogVisible: false,
      roles: undefined,
      socket: undefined,
      cluster: undefined,
    }
  },
  computed: {
    wsUrl() {
      let wsHost = location.host
      let protocol = location.protocol === 'http:' ? 'ws:' : 'wss:'
      let str = `${protocol}//${wsHost}${trimSlash(location.pathname)}/api/ssh/${this.ownerType}/${this.ownerName}/${this.nodeName}`
      return str
    },
    stateStr() {
      if (this.socketReadyState === 0) {
        return this.$t('connecting')
      } else if (this.socketReadyState === 1) {
        return this.$t('connected')
      } else if (this.socketReadyState === 2) {
        return this.$t('closing')
      } else if (this.socketReadyState === 3) {
        return this.$t('closed')
      }
      return this.$t('unknown')
    },
  },
  components: { K8sTerminalErrorHint, ChangeFontSize, Find, ChangeColor, SshQuickCommands },
  mounted() {
    window.addEventListener('resize', this.handleResize)
    this.refresh()
  },
  beforeUnmount() {
    this.xterm.dispose()
    window.removeEventListener('resize', this.handleResize)
  },
  methods: {
    handleResize() {
      this.$nextTick(() => {
        this.fitAddon.fit()
      })
    },
    refresh() {
      let _this = this
      this.xterm = new Terminal({
        fontFamily: 'Consolas,Menlo,Bitstream Vera Sans Mono,Monaco,"微软雅黑",monospace',
        fontSize: this.fontSize,
        scrollback: 100000,
        theme: { selection: 'rgba(255,36,36,0.5)', foreground: localStorage.getItem('terminalForegroundColor') || '#ccc' }
      })
      this.fitAddon = new FitAddon();
      this.xterm.loadAddon(this.fitAddon);
      this.xterm.open(document.getElementById('terminal'))
      _this.fitAddon.fit()
      console.log('init', this.xterm.cols, this.xterm.rows)

      this.socket = new WebSocket(this.wsUrl + `?cols=${this.xterm.cols}&&rows=${this.xterm.rows}`)

      this.xterm.onResize(data => {
        this.ttySize = data
        console.log(data)
        setTimeout(() => {
          let req = {
            cols: data.cols,
            rows: data.rows,
            x: window.innerWidth,
            y: window.innerHeight - 45,
          }
          _this.socket.send('1' + JSON.stringify(req))
        }, 1500)
      })
      this.xterm.onData(data => {
        _this.socket.send('0' + data)
      })


      this.socket.onmessage = function (event) {
        let data = event.data
        _this.xterm.write(data)
      }
      this.socket.onerror = function (event) {
        console.log(event)
        _this.$refs.errorHint.show(_this.wsUrl)
      }
      this.socket.onopen = function () {
        _this.fitAddon.fit()
        _this.socketReadyState = _this.socket.readyState
        _this.kuboardSprayApi.get(`/${_this.ownerType}s/${_this.ownerName}`).then(resp => {
          let inventory = resp.data.data.inventory
          _this.cluster = resp.data.data
          _this.roles = {
            kube_control_plane: inventory.all.children.target.children.k8s_cluster.children.kube_control_plane.hosts[_this.nodeName] !== undefined,
            etcd: inventory.all.children.target.children.etcd.hosts[_this.nodeName] !== undefined,
            kube_node: inventory.all.children.target.children.k8s_cluster.children.kube_node.hosts[_this.nodeName] !== undefined,
          }
          if (inventory.all.children.target.children.etcd.hosts[_this.nodeName]) {
            console.log('this is a etcd node.')
            setTimeout(() => {
              _this.socket.send('0' + `export ETCDCTL_API=3
export ETCDCTL_CERT=/etc/ssl/etcd/ssl/admin-$(hostname).pem
export ETCDCTL_KEY=/etc/ssl/etcd/ssl/admin-$(hostname)-key.pem
export ETCDCTL_CACERT=/etc/ssl/etcd/ssl/ca.pem
# 此处开始，执行您想要执行的 etcdctl 命令
`)
              if (inventory.all.hosts[_this.nodeName].kuboardspray_node_action === undefined) {
                _this.socket.send('0etcdctl member list\n')
                _this.socket.send('0\n')
              }
              if (inventory.all.children.target.children.k8s_cluster.children.kube_control_plane.hosts[_this.nodeName]) {
                _this.socket.send('0# 或者 kubectl 命令\n')
                if (inventory.all.hosts[_this.nodeName].kuboardspray_node_action === undefined) {
                  _this.socket.send('0kubectl get nodes -o wide\n\n')
                }
              }
              
            }, 2000)
          }
        }).catch(e => {
          console.log('cannot read inventory', e)
        })
      }
      this.socket.onclose = function () {
        _this.socketReadyState = _this.socket.readyState
      }

      let interval = setInterval(() => {
        try {
          if (this.socketReadyState === 1) {
            this.socket.send('')
            window.console.log('anti idle, send empty string')
          } else {
            clearInterval(interval)
            console.log('WebSocket 已关闭，停止 anti idle')
          }
        } catch (e) {
          clearInterval(interval)
          window.console.error('anti idle, failed ' + e)
        }
      }, 30000)
      setInterval(() => {
        this.socketReadyState = this.socket.readyState
      }, 2000)

      document.title = `${this.$t('terminal')} - ${this.ownerType} / ${this.ownerName} / ${this.pid}`
    }
  }
}
</script>

<style lang="scss" scoped>
.button {
  font-weight: 500;
  display: inline-block;
  margin-right: 10px;
  vertical-align: top;
  background-color: gray;
}
.confirmText {
  font-size: 15px;
  font-weight: bolder;
  margin-bottom: 5px;
}
</style>
