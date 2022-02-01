<i18n>
en:
  resourceVersion: Resource Package Version

zh:
  resourceVersion: 资源包版本

</i18n>

<template>
  <el-scrollbar height="calc(100vh - 220px)">
    <div>
      <el-form label-position="left" label-width="120px">
        <el-form-item :label="$t('resourceVersion')">
          <div class="app_text_mono version_no">{{ resource.metadata.version }}</div>
        </el-form-item>
      </el-form>
    </div>
    <CompareVersion :cluster="cluster" :version="version"></CompareVersion>
  </el-scrollbar>
</template>

<script>
import CompareVersion from './CompareVersion.vue'

export default {
  props: {
    cluster: { type: Object, required: true },
  },
  data() {
    return {
      version: undefined,
      loading: false,
    }
  },
  computed: {
    
    resource () {
      return this.cluster.resourcePackage
    }
  },
  components: { CompareVersion },
  mounted () {
    this.loadClusterVersion()
  },
  methods: {
    async loadClusterVersion() {
      this.loading = true
      await this.kuboardSprayApi.get(`/clusters/${this.cluster.name}/state/version`).then(resp => {
        this.version = resp.data.data
      }).catch(e => {
        console.log(e)
      })
      this.loading = false
    },
  }
}
</script>

<style scoped lang="scss">
.version_no {
  font-weight: bolder;
  font-size: 14px;
}
</style>
