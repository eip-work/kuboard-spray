<i18n>
en:
  releaseNote: Release Note
  content: Resource Package Content
zh:
  releaseNote: 版本说明
  content: 资源包内容
</i18n>

<template>
  <div>
    <el-tabs v-if="releaseNoteHtml" v-model="activeName" @tab-click="handleClick" type="card">
      <el-tab-pane :label="$t('releaseNote')" name="releaseNote">
        <div v-if="releaseNoteHtml" v-html="releaseNoteHtml" class="markdown"></div>
      </el-tab-pane>
      <el-tab-pane :label="$t('content')" name="content">
        <div style="padding-top: 10px;">
          <ResourcePackage :expandAll="expandAll" :resourcePackage="resourcePackage"></ResourcePackage>
        </div>
      </el-tab-pane>
    </el-tabs>
    <ResourcePackage v-else :expandAll="expandAll" :resourcePackage="resourcePackage"></ResourcePackage>
  </div>
</template>

<script>
import ResourcePackage from './ResourcePackage.vue'
import markdown from 'markdown-it'

export default {
  props: {
    releaseNote: { type: String, required: false },
    resourcePackage: { type: Object, required: true },
    expandAll: { type: Boolean, required: false, default: false },
  },
  data() {
    return {
      activeName: 'releaseNote'
    }
  },
  computed: {
    releaseNoteHtml () {
      if (this.releaseNote) {
        let md = new markdown()
        return md.render(this.releaseNote)
      }
      return undefined
    }
  },
  components: { ResourcePackage },
  mounted () {
  },
  methods: {

  }
}
</script>

<style lang="scss">
.markdown {
  font-family: Consolas,Menlo,Bitstream Vera Sans Mono,Monaco,"微软雅黑",monospace;
  font-size: 13px;
  h1 {
    font-size: 18px;
    display: none;
  }
  h2 {
    font-size: 16px;
    background-color: var(--el-color-primary-light-9);
    padding: 10px 20px;
  }

}
</style>

<style scoped lang="scss">
.package_title {
  font-weight: bolder;
  font-family: Consolas,Menlo,Bitstream Vera Sans Mono,Monaco,"微软雅黑",monospace;
}
.package_info {
  margin-left: 20px;
  font-family: Consolas,Menlo,Bitstream Vera Sans Mono,Monaco,"微软雅黑",monospace;
}
</style>
