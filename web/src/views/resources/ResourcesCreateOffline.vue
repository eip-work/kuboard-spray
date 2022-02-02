<i18n>
en: 
  createByOffline: Load Resource Package without internet connection
zh:
  createByOffline: 离线加载资源包
</i18n>

<template>
  <div style="display: inline-block; text-align: left;">
    <div style="text-align: right;">
      <el-button type="primary" icon="el-icon-download" @click="compute_dialogVisible = true">{{$t('createByOffline')}}</el-button>
    </div>
    <el-dialog v-model="compute_dialogVisible" :title="$t('createByOffline')" width="80%" top="5vh" :close-on-click-modal="false">
      <div style="height: calc(90vh - 180px);">
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
import clone from 'clone'
import compareVersions from 'compare-versions'

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
    resourcePackage () {
      if (this.content) {
        try {
          return yaml.load(this.content)
        } catch (e) {
          this.$message.error(e)
        }
      }
      return undefined
    },
    meetVersionRequirement() {
      if (this.resourcePackage === undefined) {
        return false
      }
      if (this.resourcePackage.metadata && this.resourcePackage.metadata.kuboard_spray_version && this.resourcePackage.metadata.kuboard_spray_version.min) {
        return compareVersions(window.KuboardSpray.version.trimed, this.resourcePackage.metadata.kuboard_spray_version.min) >= 0
      }
      return false
    },
    meetArch() {
      if (this.resourcePackage === undefined) {
        return false
      }
      if (this.resourcePackage.data && this.resourcePackage.data.kubernetes && this.resourcePackage.data.kubernetes.image_arch) {
        return window.KuboardSpray.version.arch ==  this.resourcePackage.data.kubernetes.image_arch
      }
      return false
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
      let req = clone(this.resourcePackage)
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
      if (!this.meetVersionRequirement) {
        this.$message.error('KuboardSpray 最低版本要求为：' + this.resourcePackage.metadata.kuboard_spray_version.min + '，当前版本为：' + window.KuboardSpray.version.version)
        return
      }
      if (!this.meetArch) {
        this.$message.error('当前 KuboardSpray 只能导入 ' + window.KuboardSpray.version.arch + ' 格式的资源包')
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
  height: calc(90vh - 260px) !important;
}

.create_resource_offline_codemirror .CodeMirror-wrap pre.CodeMirror-line, .CodeMirror-wrap pre.CodeMirror-line-like {
  word-break: break-all;
}

.create_resource_offline_codemirror .CodeMirror-scroll {
  height: calc(90vh - 260px) !important;
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
