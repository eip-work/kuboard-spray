<i18n>
en:
  "logs": "Logs"
  "find": "Find"
  "connecting": "Connecting"
  "connected": "Connected"
  "closing": "Closing"
  "closed": "Closed"
  "unknown": "Unknown"
  "confirmToClearLogs": "You are about to clear the lines displayed, do you confirm ?"
  "succeedInClear": "Succeed in clear logs."
zh:
  "logs": "日志"
  "find": "查 找"
  "connecting": "正在连接"
  "connected": "日志追踪进行中"
  "closing": "正在关闭"
  "closed": "已断开"
  "unknown": "未知状态"
  "confirmToClearLogs": "此操作将清空当前显示的日志, 是否继续？"
  "succeedInClear": "清除日志成功"
</i18n>

<template>
  <div>
    <div style="background-color: #3a3333; color: #fffeff; font-weight: 500; padding: 8px 20px 8px 20px; line-height: 22px; text-align: center;">
      {{this.$t('logs')}} - <span class="app_text_mono">{{ownerType}}/{{ownerName}}/{{pid}}</span>
      <span style="float:left;">
        <ChangeFontSize class="button" :terminal="xterm" :fitAddon="fitAddon"></ChangeFontSize>
        <ChangeColor class="button"></ChangeColor>
        <el-button type="info" icon="el-icon-search" @click="$refs.find.show()">{{$t('find')}}</el-button>
        <Find ref="find" class="button" :terminal="xterm"></Find>
      </span>
      <span :style="`float: right; font-size: 14px; font-weight: 600; color: ${socketReadyState === 1 ? '#33FF33' : '#FF6600'};`">
        {{ stateStr }}
      </span>
      <el-tag size="medium" style="float: right; margin-right: 20px;" type="primary">{{ $t('finished') }}</el-tag>
    </div>
    <div id="terminal" :style="`width: 100%; height: calc(100vh - 39px); background-color: black;`"></div>
    <K8sTerminalErrorHint ref="errorHint"/>
  </div>
</template>

<script>
import { Terminal } from 'xterm'
import { FitAddon } from 'xterm-addon-fit'
// import { AttachAddon } from 'xterm-addon-attach'
import 'xterm/css/xterm.css'
import K8sTerminalErrorHint from './K8sTerminalErrorHint'
import ChangeFontSize from './ChangeFontSize'
import Find from './Find'
import ChangeColor from './ChangeColor'
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
    pid: { type: String, required: true },
    file: { type: String, required: true },
    tailLines: { type: Number, required: false, default: 500 },
    autoTrace: { type: Boolean, required: false, default: true }
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
    }
  },
  computed: {
    wsUrl() {
      let wsHost = location.host
      let protocol = location.protocol === 'http:' ? 'ws:' : 'wss:'
      let str = `${protocol}//${wsHost}${trimSlash(location.pathname)}/api/ssh/cluster/c1/node1`
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
  components: { K8sTerminalErrorHint, ChangeFontSize, Find, ChangeColor },
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
      this.socket = new WebSocket(this.wsUrl)

      this.xterm.open(document.getElementById('terminal'))
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
        console.log('send:', '0' + data)
        // _this.socket.send(data)
      })


      this.socket.onmessage = function (event) {
        // let data = decode(event.data)
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
      }
      this.socket.onclose = function () {
        _this.socketReadyState = _this.socket.readyState
      }

      // this.xterm.attach(this.socket)
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

      document.title = `${this.$t('logs')} - ${this.ownerType} / ${this.ownerName} / ${this.pid}`
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
