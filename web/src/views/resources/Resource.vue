<i18n>
en:
  resourceList: Resources List
  confirmToDelete: Do you confirm to delete the resource package ?
zh:
  resourceList: 资源包列表
  confirmToDelete: 是否确认要删除该资源包
</i18n>

<template>
  <div>
    <ControlBar :title="$t('obj.resource')">
      <span class="app_text_mono" style="margin-right: 20px;">{{name}}</span>
      <el-tag type="warning" effect="dark" style="margin-right: 20px;">
        <i class="el-icon-cloudy"></i>
        {{$t('online')}}
      </el-tag>
      <ResourceDownload v-if="resource" :resource="resource" :loading="loading" @refresh="refresh"></ResourceDownload>
      <el-popconfirm :title="$t('confirmToDelete')" @confirm="removeResource">
        <template #reference>
          <el-button type="danger" icon="el-icon-delete">{{ $t('msg.delete') }}</el-button>
        </template>
      </el-popconfirm>
    </ControlBar>
    <el-card>
      <el-skeleton v-if="loading"></el-skeleton>
      <ResourceDetails v-else-if="resource" :resourcePackage="resource.package" expandAll></ResourceDetails>
    </el-card>
  </div>
</template>

<script>
import mixin from '../../mixins/mixin.js'
import ResourceDetails from './details/ResourceDetails.vue'
import ResourceDownload from './ResourceDownload.vue'

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
      resource: undefined,
      loading: false,
    }
  },
  computed: {
  },
  components: { ResourceDetails, ResourceDownload },
  mounted () {
    this.refresh()
  },
  methods: {
    refresh () {
      this.loading = true
      this.kuboardSprayApi.get(`/resources/${this.name}`).then(resp => {
        this.resource = resp.data.data
        this.loading = false
      }).catch(e => {
        this.loading = false
        console.log(e)
      })
    },
    removeResource() {
      this.kuboardSprayApi.delete(`/resources/${this.name}`).then(() => {
        this.$message.success(this.$t('msg.delete_succeeded'))
        this.$router.replace(`/settings/resources`)
      }).catch(e => {
        this.$message.error(this.$t('msg.delete_failed', {msg: e.response.data.message}))
      })
    }
  }
}
</script>

<style scoped lang="scss">

</style>
