<i18n>
en:
  title: "Load resource package {name}."
  selectSource: "Select a source"
  download: "Download resource package"
  upload: "Upload Resource Package"
  downloadArchitecture: "Select CPU architecture"
  useProxy: " Use proxy"
  retries: "Retries"
zh:
  title: "加载资源包 {name}"
  selectSource: "选择一个源"
  download: "加载资源包"
  upload: "加载离线资源包"
  downloadArchitecture: "选择 CPU 架构"
  useProxy: "使用代理"
  retries: "重试次数"
</i18n>

<template>
  <ExecuteTask :history="resource.history" :startTask="startTask" :label="t('download')"
    :title="t('title', { name: resource.package.metadata.version })" :loading="loading" @refresh="$emit('refresh')">
    <el-form @submit.prevent.stop :model="form" ref="form" label-position="left" label-width="120px">
      <el-form-item :label="t('downloadArchitecture')">
        <el-radio-group v-model="form.downloadArchitecture">
          <el-radio-button label="amd64" value="amd64"></el-radio-button>
          <el-radio-button label="arm64" value="arm64"></el-radio-button>
          <el-radio-button label="both" value="amd64,arm64"></el-radio-button>
        </el-radio-group>
      </el-form-item>
      <el-form-item :label="t('retries')">
        <el-input-number v-model="form.retries" :max="5" :min="0"></el-input-number>
      </el-form-item>
      <el-form-item :label="t('useProxy')">
        <el-switch v-model="enableProxyOnDownload" />
      </el-form-item>
      <el-form-item v-if="enableProxyOnDownload" label="Http proxy">
        <el-input v-model="httpProxy" />
      </el-form-item>
    </el-form>
  </ExecuteTask>
</template>

<script>
import ExecuteTask from "../common/task/ExecuteTask.vue";
import { useI18n } from "vue-i18n";

export default {
  props: {
    resource: { type: Object, required: true },
    loading: { type: Boolean, required: false },
  },
  data() {
    return {
      form: {
        retries: 0,
        downloadArchitecture: 'amd64',
        downloadFrom: '',
        enableProxyOnDownload: false,
        httpProxy: ''
      },
      sourceRules: [{ required: true, message: this.i18n("selectSource"), trigger: "change" }]
    };
  },
  setup() {
    const { t } = useI18n({
      inheritLocale: true,
      useScope: "local"
    });
    return { i18n: t };
  },
  computed: {
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
  components: { ExecuteTask },
  emits: ["refresh"],
  methods: {
    startTask() {
      console.log(this.resource.package.metadata.version)
      return new Promise((resolve, reject) => {
        let request = {
          version: this.resource.package.metadata.version,
          retries: this.form.retries + "",
          downloadArchitecture: this.form.downloadArchitecture,
          enableProxyOnDownload: this.enableProxyOnDownload.toString(), // 转为 "true" 或 "false"
          httpProxy: this.httpProxy,
        };
        this.pangeeClusterApi
          .post(`/resources/${request.version}/download`, request)
          .then(resp => {
            this.$router.replace(`/settings/resources/${request.version}`);
            resolve(resp.data.data.pid);
          })
          .catch(e => {
            reject(e);
          });
      });
    },
  }
};
</script>

<style scoped lang="css">
.confirmText {
  font-size: 15px;
  color: var(--el-color-danger);
  font-weight: bold;
}

.hidden {
  display: none;
}
</style>
