<i18n>
en:
  resourceList: Resources List
  online: 'On Air '
  download: Download
  minVersionRequired: Min version required to kuboardspray
zh:
  resourceList: 资源包列表
  online: 未下载
  download: 下 载
  minVersionRequired: KuboardSpray最低版本要求
</i18n>

<template>
  <div>
    <ControlBar :title="name">
      <template v-if="resourcePackage">
        <ResourceDownload  v-if="meetVersionRequirement" action="download" :resource="{package: resourcePackage, history:{task_type: 'resource', task_name: name, processing: false, success_tasks: []}}"></ResourceDownload>
        <template v-else>
          <el-tag type="danger" effect="dark">{{ $t('minVersionRequired') }}</el-tag>
          <el-tag type="danger" class="app_text_mono">{{resourcePackage.metadata.kuboard_spray_version.min}}</el-tag>
        </template>
      </template>
    </ControlBar>
    <div class="app_block_title">
      {{ $t('obj.resource') }}
      <el-tag type="warning" effect="dark">
        <i class="el-icon-cloudy"></i>
        {{$t('online')}}
      </el-tag>
    </div>
    <el-card>
      <el-skeleton v-if="loading" animated></el-skeleton>
      <ResourceDetails v-else-if="resourcePackage" :releaseNote="releaseNote" :resourcePackage="resourcePackage" expandAll></ResourceDetails>
    </el-card>
  </div>
</template>

<script>
import mixin from '../../mixins/mixin.js'
import ResourceDetails from './details/ResourceDetails.vue'
import axios from 'axios'
import yaml from 'js-yaml'
import ResourceDownload from './ResourceDownload.vue'
import compareVersions from 'compare-versions'

export default {
  mixins: [mixin],
  percentage () {
    return 100
  },
  breadcrumb () {
    return [
      { label: this.$t('resourceList'), to: '/settings/resources' },
      { label: this.name },
    ]
  },
  refresh () {
    this.refresh()
  },
  props: {
    name: { type: String, required: true },
    mode: { type: String, required: false, default: 'view' },
  },
  data() {
    return {
      resourcePackage: undefined,
      releaseNote: undefined,
      loading: false,
    }
  },
  computed: {
    meetVersionRequirement() {
      if (this.resourcePackage === undefined) {
        return false
      }
      return compareVersions(window.KuboardSpray.version.trimed, this.resourcePackage.metadata.kuboard_spray_version.min) >= 0
    },
  },
  components: { ResourceDetails, ResourceDownload },
  mounted () {
    this.refresh()
  },
  methods: {
    async refresh () {
      this.loading = true
      await axios.get(`https://addons.kuboard.cn/v-kuboard-spray/${this.name}/package.yaml?nocache=${new Date().getTime()}`).then(resp => {
        this.resourcePackage = yaml.load(resp.data)
      }).catch(e => {
        console.log(e)
        this.$message.error('离线环境')
      })
      await axios.get(`https://addons.kuboard.cn/v-kuboard-spray/${this.name}/release.md?nocache=${new Date().getTime()}`).then(resp => {
        this.releaseNote = resp.data
      }).catch(e => {
        console.log(e)
      })
      this.loading = false
    },
  }
}
</script>

<style scoped lang="css">

</style>
