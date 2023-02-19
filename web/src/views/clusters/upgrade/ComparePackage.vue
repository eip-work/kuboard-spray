<i18n>
en:
  releaseNote: Release Note
  content: Resource Package Content
  compare: Compare installed version with target
  compareTitle: Compare Resource Packages
zh:
  releaseNote: 版本说明
  content: 资源包内容
  compare: 对比已安装版本与目标版本
  compareTitle: 对比资源包：
</i18n>

<template>
  <div>
    <el-dialog v-model="dialogVisible" width="80%" top="5vh" :close-on-click-modal="false">
      <template #header>
        <div style="background-color: var(--el-color-primary-light-9); margin-right: 25px; padding: 10px;">
          {{$t('compareTitle')}}
          <el-tag>
            <span class="app_text_mono" style="font-size: 14px;">{{cluster.resourcePackage.metadata.version}}</span>
          </el-tag>
          <span style="margin: 0 10px; font-size: 12px;">v.s.</span>
          <el-tag type="warning" effect="dark">
            <span class="app_text_mono" style="font-size: 14px; font-weight: bolder;">{{target ? target.yaml.metadata.version : ''}}</span>
          </el-tag>
        </div>
      </template>
      <div v-if="dialogVisible">
        <el-tabs v-if="releaseNoteHtml" v-model="activeName" type="card">
          <el-tab-pane :label="$t('compare')" name="compare">
            <ComparePackageTable :cluster="cluster" :target="target.yaml"></ComparePackageTable>
          </el-tab-pane>
          <el-tab-pane :label="$t('releaseNote')" name="releaseNote">
            <el-scrollbar height="calc(90vh - 245px)">
              <el-skeleton v-if="releaseNoteLoading"></el-skeleton>
              <div v-else-if="releaseNoteHtml" v-html="releaseNoteHtml" class="markdown"></div>
            </el-scrollbar>
          </el-tab-pane>
          <el-tab-pane :label="$t('content')" name="content">
            <el-scrollbar height="calc(90vh - 245px)">
              <div style="padding-top: 10px;">
                <ResourcePackage :expandAll="true" :resourcePackage="target.yaml"></ResourcePackage>
              </div>
            </el-scrollbar>
          </el-tab-pane>
        </el-tabs>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false" icon="el-icon-close" type="primary">{{$t('msg.close')}}</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import ResourcePackage from '../../resources/details/ResourcePackage.vue'
import ComparePackageTable from './ComparePackageTable.vue'
import markdown from 'markdown-it'
import axios from 'axios'

export default {
  props: {
    cluster: { type: Object, required: true },
  },
  data() {
    return {
      dialogVisible: false,
      activeName: 'compare',
      releaseNote: undefined,
      releaseNoteLoading: false,
      target: undefined,
    }
  },
  computed: {
    releaseNoteHtml () {
      if (this.releaseNote) {
        let md = new markdown()
        return md.render(this.releaseNote)
      }
      return "# hello"
    },
  },
  components: { ResourcePackage, ComparePackageTable },
  mounted () {
  },
  methods: {
    show(target) {
      this.dialogVisible = true
      this.releaseNote = undefined
      this.releaseNoteLoading = true
      this.target = target
      if (this.target.imported || this.target.isOffline) {
        this.kuboardSprayApi.get(`/resources/${this.target.version}/release_note`).then(resp => {
          this.releaseNote = resp.data.data.release_note
        }).catch(e => {
          console.log(e)
        })
      } else {
        axios.get(`https://addons.kuboard.cn/v-kuboard-spray/${this.target.version}/release.md?nocache=${new Date().getTime()}`).then(resp => {
          this.releaseNote = resp.data
        }).catch(e => {
          console.log(e)
        })
      }
      this.releaseNoteLoading = false
    }
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

<style scoped lang="css">
.package_title {
  font-weight: bolder;
  font-family: Consolas,Menlo,Bitstream Vera Sans Mono,Monaco,"微软雅黑",monospace;
}
.package_info {
  margin-left: 20px;
  font-family: Consolas,Menlo,Bitstream Vera Sans Mono,Monaco,"微软雅黑",monospace;
}
</style>
