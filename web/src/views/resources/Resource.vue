<i18n>
en:
  resourceList: Resources List
  confirmToDelete: Do you confirm to delete the resource package ?
  loaded: Loaded
  not_load: Not loaded
zh:
  resourceList: 资源包列表
  confirmToDelete: 是否确认要删除该资源包
  loaded: '已成功加载'
  not_load: '未加载'
</i18n>

<template>
  <div>
    <ControlBar :title="$t('obj.resource')">
      <span class="app_text_mono" style="margin-right: 20px;">{{name}}</span>
      <template v-if="resource">
        <el-tag v-if="isInstalled" type="success" effect="dark" style="margin-right: 20px;">
          <el-icon :size="14" style="width: 14px; height: 14px; vertical-align: bottom;">
            <el-icon-cloudy></el-icon-cloudy>
          </el-icon>
          {{$t('loaded')}}
        </el-tag>
        <el-tag v-else type="danger" effect="dark" style="margin-right: 20px;">
          <el-icon :size="14" style="width: 14px; height: 14px; vertical-align: bottom;">
            <el-icon-cloudy></el-icon-cloudy>
          </el-icon>
          {{$t('not_load')}}
        </el-tag>
        <ResourceDownload :resource="resource" action="reload" :loading="loading" @refresh="refresh"></ResourceDownload>
        <el-popconfirm v-if="!resource.history.processing" :title="$t('confirmToDelete')" @confirm="removeResource">
          <template #reference>
            <el-button type="danger" icon="el-icon-delete">{{ $t('msg.delete') }}</el-button>
          </template>
        </el-popconfirm>
      </template>
    </ControlBar>
    <el-card>
      <el-skeleton v-if="loading"></el-skeleton>
      <ResourceDetails v-else-if="resource" :releaseNote="releaseNote" :resourcePackage="resource.package" expandAll></ResourceDetails>
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
      releaseNote: undefined,
      loading: false,
    }
  },
  computed: {
    isInstalled () {
      if (this.resource) {
        return this.resource.history.success_tasks.length > 0
      }
      return false
    }
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
      this.kuboardSprayApi.get(`/resources/${this.name}/release_note`).then(resp => {
        this.releaseNote = resp.data.data.release_note
      }).catch(e => {
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

<style scoped lang="css">

</style>
