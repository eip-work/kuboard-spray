<i18n>
en: 
  createByOffline: Load Resource Package without internet connection
zh:
  createByOffline: 离线加载资源包
</i18n>

<template>
  <div>
    <div style="text-align: right;">
      <el-button type="primary" icon="el-icon-download" @click="compute_dialogVisible = true">{{$t('createByOffline')}}</el-button>
    </div>
    <el-dialog v-model="compute_dialogVisible" :title="$t('createByOffline')" width="80%" top="10vh" :close-on-click-modal="false">
      <div style="height: calc(80vh - 180px);">
        <el-alert :closable="false" type="warning" class="app_margin_bottom">
          <div style="line-height: 20px;">
            如果您 KuboardSpray 所在机器不能联网，请从如下链接中的任意一个获取离线导入文件：
            <li><el-link href="https://www.kuboard-spray.cn/support" target="_blank">https://www.kuboard-spray.cn/support</el-link></li>
            <li><el-link href="https://www.kuboard.cn/support/kuboard-spray" target="_blank">https://www.kuboard.cn/support/kuboard-spray</el-link></li>
          </div>
        </el-alert>
        <div v-if="contentVisible">
          <Codemirror v-model:value="content" :options="cmOptions" class="create_resource_offline_codemirror"></Codemirror>
        </div>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="compute_dialogVisible = false" icon="el-icon-close">{{$t('msg.close')}}</el-button>
          <el-button type="primary" @click="downloadOffline" icon="el-icon-check" :disabled="cannotSave">{{$t('msg.ok')}}</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import Codemirror from "codemirror-editor-vue3"
import "codemirror/theme/darcula.css"
import "codemirror/mode/yaml/yaml.js"
import yaml from 'js-yaml'

export default {
  props: {
  },
  data() {
    return {
      dialogVisible: false,
      contentVisible: false,
      content: "",
      cmOptions: {
        mode: "yaml", // Language mode
        theme: "darcula", // Theme
        lineNumbers: true, // Show line number
        lineWrapping: true,
        smartIndent: true, // Smart indent
        indentUnit: 2, // The smart indent unit is 2 spaces in length
        foldGutter: false, // Code folding
        styleActiveLine: true, // Display the style of the selected row
      },
    }
  },
  computed: {
    cannotSave () {
      if (this.content) {
        return false
      }
      return true
    },
    compute_dialogVisible: {
      get () { return this.dialogVisible },
      set (v) {
        if (v) {
          this.content = ''
          this.dialogVisible = true
          setTimeout(() => {
            this.contentVisible = true
          }, 200)
        } else {
          this.contentVisible = false
          this.dialogVisible = false
        }
      }
    }
  },
  components: { Codemirror },
  mounted () {
  },
  methods: {
    downloadOffline () {
      let req = yaml.load(this.content)
      if (!req.downloadFrom) {
        this.$message.error('缺少 downloadFrom 字段，请正确复制文件内容')
        return
      }
      if (!req.data) {
        this.$message.error('缺少 data 字段，请正确复制文件内容')
        return
      }
      if (!req.metadata) {
        this.$message.error('缺少 metadata 字段，请正确复制文件内容')
        return
      }
      if (!req.metadata.version) {
        this.$message.error('缺少 metadata.version 字段，请正确复制文件内容')
        return
      }
      let request = {
        package: req,
        downloadFrom: req.downloadFrom
      }
      this.kuboardSprayApi.post(`/resources/${req.metadata.version}/download`, request).then(resp => {
        this.openUrlInBlank(`/#/tail/resource/${req.metadata.version}/history/${resp.data.data.pid}/execute.log`)
        this.$router.push(`/settings/resources/${req.metadata.version}`)
      }).catch(e => {
        console.log(e)
        if (e.response && e.response.data.message)
        this.$message.error('加载失败：' + e.response.data.message)
      })
    }
  }
}
</script>

<style>

.create_resource_offline_codemirror .CodeMirror {
  height: calc(80vh - 260px) !important;
}

.create_resource_offline_codemirror .CodeMirror-wrap pre.CodeMirror-line, .CodeMirror-wrap pre.CodeMirror-line-like {
  word-break: break-all;
}

.create_resource_offline_codemirror .CodeMirror-scroll {
  height: calc(80vh - 260px) !important;
  overflow-y: hidden;
  overflow-x: auto;
}
/* 
.create_resource_offline_codemirror .CodeMirror-line {
  left: 40px !important;
} */
</style>

<style scoped lang="scss">

</style>
