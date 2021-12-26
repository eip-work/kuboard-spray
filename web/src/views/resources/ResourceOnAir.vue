<i18n>
en:
  resourceList: Resources List
  online: 'On Air '
zh:
  resourceList: 资源包列表
  online: 未下载
</i18n>

<template>
  <div>
    <div class="app_block_title">
      {{ $t('obj.resource') }}
      <el-tag type="warning" effect="dark">
        <i class="el-icon-cloudy"></i>
        {{$t('online')}}
      </el-tag>
    </div>
    <el-card>
      <ResourceDetails v-if="resourcePackage" :resourcePackage="resourcePackage" expandAll></ResourceDetails>
    </el-card>
  </div>
</template>

<script>
import mixin from '../../mixins/mixin.js'
import ResourceDetails from './details/ResourceDetails.vue'
import axios from 'axios'
import yaml from 'js-yaml'

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
    }
  },
  computed: {
  },
  components: { ResourceDetails },
  mounted () {
    this.refresh()
  },
  methods: {
    async refresh () {
      await axios.get(`https://addons.kuboard.cn/v-kuboard-spray-resources/${this.name}/package.yaml`).then(resp => {
        this.resourcePackage = yaml.load(resp.data)
      }).catch(e => {
        console.log(e)
        this.$message.error('离线环境')
      })
    }
  }
}
</script>

<style scoped lang="scss">

</style>
