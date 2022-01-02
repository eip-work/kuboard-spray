<i18n>
en:
  "logs": "Logs"
  "find": "Find"
  "initialLines": "Initial Lines"
  "displayLines": "Display Lines"
  "autoTail": "Auto Tail"
  "stopAutoTail": "Stop Tail"
  "cachedLines": "Cached Lines"
  "cachedPages": "Cached Pages"
  "linesPerPage": "Lines per page"
  "redLinesExplain": "The first line in a page is in RED."
  "tab": "Tab - view by page"
  "overFlow": "There are more than 10,000 lines in cache, disconnected to server, you can continue to view the content in cache."
  "overFlowTitle": "Cache overflow"
  "connecting": "Connecting"
  "connected": "Connected"
  "closing": "Closing"
  "closed": "Closed"
  "unknown": "Unknown"
  "confirmToClearLogs": "You are about to clear the lines displayed, do you confirm ?"
  "succeedInClear": "Succeed in clear logs."
  failedToKill: Failed to kill process.
  killed: Succeeded in kill process.
  confirmToKill: This is going to kill the ansible task running in the background, may introduce unexpected effect, do you confirm to continue?
  confirm: I do confirm to kill the task.
  forceKill: Force to kill the task
  finished: Task Completed
zh:
  "logs": "日志"
  "find": "查 找"
  "initialLines": "初始行数"
  "displayLines": "显示行数"
  "autoTail": "自动追踪"
  "stopAutoTail": "暂停追踪"
  "cachedLines": "已缓存行数"
  "cachedPages": "已缓存页数"
  "linesPerPage": "每页行数"
  "redLinesExplain": "红色为每页第一行，出现折行时重新计算"
  "tab": "Tab键可逐页查看"
  "overFlow": "缓存内容已经超过 10000 条，已断开与服务器的连接，您可以继续查看缓存中的内容"
  "overFlowTitle": "缓存过大"
  "connecting": "正在连接"
  "connected": "日志追踪进行中"
  "closing": "正在关闭"
  "closed": "已断开"
  "unknown": "未知状态"
  "confirmToClearLogs": "此操作将清空当前显示的日志, 是否继续？"
  "succeedInClear": "清除日志成功"
  failedToKill: 结束进程失败
  killed: 已成功结束进程
  confirmToKill: 将要强制结束此 ansible 进程，有可能会导致不能预测的问题，是否继续？
  confirm: 我确定要强制结束此任务
  forceKill: 强制结束任务
  finished: 任务已结束
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
      <el-button v-if="isRunning" style="float: right; margin-right: 20px;" type="danger" @click="dialogVisible = true">{{ $t('forceKill') }}</el-button>
      <el-tag v-else size="medium" style="float: right; margin-right: 20px;" type="primary">{{ $t('finished') }}</el-tag>
    </div>
    <div id="terminal" :style="`width: 100%; height: calc(100vh - 39px); background-color: black;`"></div>
    <el-dialog :title="$t('msg.prompt')" v-model="dialogVisible" width="60%" top="calc(50vh - 180px)" :close-on-click-modal="false">
      <el-alert type="error" :closable="false" effect="dark" show-icon>
        <div class="confirmText">{{$t('confirmToKill')}}</div>
      </el-alert>
      <template #footer>
        <div>
          <el-button type="default" icon="el-icon-close" @click="dialogVisible = false">{{ $t('msg.close') }}</el-button>
          <el-button type="primary" icon="el-icon-check" @click="killProcess">{{ $t('confirm') }}</el-button>
        </div>
      </template>
    </el-dialog>
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
      let str = `${protocol}//${wsHost}${trimSlash(location.pathname)}/api/execute/${this.ownerType}/${this.ownerName}/tail/${this.pid}/${this.file}`
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
    // clear() {
    //   this.$confirm(this.$t('confirmToClearLogs'), this.$t('message.prompt'), {
    //     type: 'warning'
    //   }).then(() => {
    //     this.xterm.clear()
    //     this.$message({
    //       type: 'success',
    //       message: this.$t('succeedInClear')
    //     });
    //   }).catch(() => {
    //     this.$message({
    //       type: 'info',
    //       message: this.$t('message.canceled')
    //     });          
    //   });
    // },
    handleResize() {
      this.$nextTick(() => {
        this.fitAddon.fit()
      })
    },
    killProcess () {
      this.kuboardSprayApi.delete(`/execute/${this.ownerType}/${this.ownerName}/kill/${this.pid}`).then(resp=> {
        if (resp.data.code === 200) {
          this.$message.success(this.$t('killed'))
          this.dialogVisible = false
        }
      }).catch(e => {
        this.$message.error(this.$t('failedToKill' + e.response.data.msg))
      })
    },
    refresh() {
      let _this = this
      this.xterm = new Terminal({
        fontFamily: 'Consolas,Menlo,Bitstream Vera Sans Mono,Monaco,"微软雅黑",monospace',
        fontSize: this.fontSize,
        theme: { selection: 'rgba(255,36,36,0.5)', foreground: localStorage.getItem('terminalForegroundColor') || '#ccc' }
      })
      this.fitAddon = new FitAddon();
      this.xterm.loadAddon(this.fitAddon);
      this.socket = new WebSocket(this.wsUrl)

      this.xterm.open(document.getElementById('terminal'))
      this.xterm.onResize(data => {
        this.ttySize = data
      })
      this.xterm.onKey(event => {
        if (event.domEvent.code === 'Enter') {
          _this.xterm.writeln('')
        }
      })

      this.fitAddon.fit()

      this.socket.onmessage = function (event) {
        _this.xterm.writeln(event.data)
        if (event.data.indexOf('KUBOARD SPRAY *****************************************************************') >= 0) {
          _this.isRunning = false
        }
      }
      this.socket.onerror = function (event) {
        console.log(event)
        _this.$refs.errorHint.show(_this.wsUrl)
      }
      this.socket.onopen = function () {
        _this.socketReadyState = _this.socket.readyState
      }
      this.socket.onclose = function () {
        _this.socketReadyState = _this.socket.readyState
      }
      let interval = setInterval(() => {
        try {
          if (this.socketReadyState === 1) {
            this.socket.send('0')
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
