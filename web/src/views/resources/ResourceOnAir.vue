<i18n>
en:
  resourceList: Resources List
  online: 'On Air '
  download: Download
  minVersionRequired: Min version required to pangeecluster
  useProxy: "Enable Proxy"
  proxy: Proxy
  proxyPlaceHolder: Proxy URL that PangeeCluster server can use, e.g. http://proxy.example.com:8080
  prompt: Download ResourcePackage to PangeeCluster server
  timeout: Timeout (seconds)
zh:
  resourceList: 资源包列表
  online: 未下载
  download: 下 载
  minVersionRequired: PangeeCluster最低版本要求
  useProxy: "使用代理"
  proxy: 代理
  proxyPlaceholder: PangeeCluster 服务器可以使用的代理地址，例如 http://proxy.example.com:8080
  prompt: 将资源包下载到服务器
  timeout: 超时时间（秒）
</i18n>

<template>
  <div>
    <ControlBar :title="version">
      <template v-if="resourcePackage">
        <el-popover v-if="meetVersionRequirement" trigger="click" width="630px" placement="bottom-start">
          <template #reference>
            <el-button type="warning" icon="el-icon-download" :loading="downloading">
              {{ t('download') }}
            </el-button>
          </template>
          <el-form @submit.prevent.stop label-position="left" label-width="110px">
            <div style="font-weight: bolder; line-height: 48px; font-size: 16px;">{{ t("prompt") }}</div>
            <el-form-item :label="t('timeout')" prop="timeout">
              <el-input-number v-model="form.timeout" :max="1200" :min="10" :step="10"></el-input-number>
            </el-form-item>
            <el-form-item :label="t('useProxy')">
              <el-switch v-model="enableProxyOnDownload" />
            </el-form-item>
            <el-form-item :label="t('proxy')">
              <el-input v-model="httpProxy" :placeholder="t('proxyPlaceholder')"></el-input>
            </el-form-item>
          </el-form>
          <div style="text-align: right;">
            <el-button type="warning" icon="el-icon-download" @click="downloadResource" :loading="downloading">
              {{ t('download') }}
            </el-button>
          </div>
        </el-popover>
        <template v-else>
          <el-tag type="danger" effect="dark">{{ t('minVersionRequired') }}</el-tag>
          <el-tag type="danger" class="app_text_mono">{{ resourcePackage.metadata.pangee_cluster_version.min }}</el-tag>
        </template>
      </template>
    </ControlBar>
    <div class="app_block_title app_margin_bottom">
      {{ $t('obj.resource') }}
      <el-tag type="warning" effect="dark">
        <el-icon :size="14" style="width: 14px; height: 14px; vertical-align: bottom;">
          <el-icon-cloudy></el-icon-cloudy>
        </el-icon>
        {{ t('online') }}
      </el-tag>
    </div>
    <el-card style="margin: 20px;">
      <el-skeleton v-if="loading" animated></el-skeleton>
      <ResourceDetails v-else-if="resourcePackage" :releaseNote="releaseNote" :resourcePackage="resourcePackage"
        :releaseNoteBaseUrl="`https://raw.githubusercontent.com/opencmit/pangee-cluster-resource-package/${this.tag}`"
        expandAll>
      </ResourceDetails>
    </el-card>
  </div>
</template>

<script>
import mixin from '../../mixins/mixin.js'
import ResourceDetails from './details/ResourceDetails.vue'
import yaml from 'js-yaml'
import ResourceDownload from './ResourceDownload.vue'
import compareVersions from 'compare-versions'
import axios from "axios"

export default {
  mixins: [mixin],
  percentage() {
    return 100
  },
  breadcrumb() {
    return [
      { label: this.t('resourceList'), to: '/settings/resources' },
      { label: this.version },
    ]
  },
  refresh() {
    this.refresh()
  },
  props: {
    tag: { type: String, required: true },
    version: { type: String, required: true },
    mode: { type: String, required: false, default: 'view' },
  },
  data() {
    return {
      resourcePackage: undefined,
      releaseNote: undefined,
      loading: false,
      downloading: false,
      form: {
        timeout: 30
      }
    }
  },
  computed: {
    meetVersionRequirement() {
      if (this.resourcePackage === undefined) {
        return false
      }
      return compareVersions(window.PangeeCluster.version.trimed, this.resourcePackage.metadata.pangee_cluster_version.min) >= 0
    },
    enableProxyOnDownload: {
      get () {
        return this.$store.state.header.proxy.enableProxyOnDownload
      },
      set (value) {
        this.$store.commit('header/CHANGE_PROXY', { key: 'enableProxyOnDownload', value })
      }
    },
    httpProxy: {
      get () {
        return this.$store.state.header.proxy.httpProxy
      },
      set (value) {
        this.$store.commit('header/CHANGE_PROXY', { key: 'httpProxy', value })
      }
    },
  },
  components: { ResourceDetails, ResourceDownload },
  mounted() {

    this.refresh()
  },
  methods: {
    async refresh() {
      this.loading = true
      axios.get(`https://raw.githubusercontent.com/opencmit/pangee-cluster-resource-package/${this.tag}/package.yaml?nocache=${Math.random()}`).then(resp => {
        this.resourcePackage = yaml.load(resp.data)
      }).catch(e => {
        console.log(e)
        this.$message.error('离线环境')
      })
      axios.get(`https://raw.githubusercontent.com/opencmit/pangee-cluster-resource-package/${this.tag}/release.md?nocache=${Math.random()}`).then(resp => {
        this.releaseNote = resp.data
      }).catch(e => {
        console.log(e)
      })
      this.loading = false
    },
    async downloadResource() {
      let source = `https://github.com/opencmit/pangee-cluster-resource-package/archive/refs/tags/${this.tag}.zip`
      const formData = new FormData();
      formData.append('sourceUrl', source);
      if (this.enableProxyOnDownload && this.httpProxy) {
        formData.append('proxyUrl', this.httpProxy);
      }
      this.downloading = true
      await this.pangeeClusterApi
        .post('/resources/upload', formData, { headers: { 'Content-Type': 'multipart/form-data' } })
        .then(resp => {
          this.$router.replace(`/settings/resources/${resp.data.data.version}`);
        })
        .catch(e => {
          this.$message.error(e.response.data.message)
        });
      this.downloading = false
    },
  }
}
</script>

<style scoped lang="css"></style>
