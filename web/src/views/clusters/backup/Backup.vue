<i18n>
en:
  node: Node name
  etcd_member_name: ETCD member name
  backup_name: Backup name
  size: Size
  refresh: Refresh
  delete_message: Delete selected backups.
zh:
  node: 节点名称
  etcd_member_name: ETCD 成员名称
  backup_name: 备份名称
  size: 备份文件大小
  refresh: 刷 新
  delete_message: 删除选中的备份文件。
</i18n>

<template>
  <div>
    <el-scrollbar height="calc(100vh - 220px)">
      <div v-if="cluster">
        <template v-if="cluster.resourcePackage && cluster.resourcePackage.data.supported_playbooks.backup_etcd">
          <div style="display: flex;">
            <ConfirmButton type="danger" :text="$t('msg.delete')" icon="el-icon-delete" @confirm="remove_backups" placement="right" :message="$t('delete_message')"
              :disabled="selection.length === 0" style="margin-right: 10px;">
            </ConfirmButton>
            <el-button type="primary" plain icon="el-icon-refresh" style="margin-right: 10px;" @click="list">{{ $t('refresh') }}</el-button>
            <BackupTask :cluster="cluster" style="margin-right: 10px;" :loading="loading"></BackupTask>
          </div>
          <div class="app_margin_bottom"></div>
          <el-alert v-if="backups.length === 0" type="warning" :closable="false">
            当前没有备份
          </el-alert>
          <el-skeleton v-else-if="loading" animated></el-skeleton>
          <div v-else>
            <el-table :data="backups" style="width: 100%" @selection-change="handleSelectionChange">
              <el-table-column type="selection" width="55"></el-table-column>
              <el-table-column prop="node_name" :label="$t('node_name')" sortable />
              <el-table-column prop="etcd_member_name" :label="$t('etcd_member_name')" sortable />
              <el-table-column prop="backup_name" :label="$t('backup_name')" sortable />
              <el-table-column prop="size" :label="$t('size')" sortable align="right" width="200">
                <template #default="scope">
                  <span style="margin-right: 10px;">
                    {{ scope.row.size ? String(scope.row.size).replace(/(\d)(?=(\d{3})+$)/g, "$1,") : ''}}
                  </span>
                </template>
              </el-table-column>
            </el-table>
          </div>
        </template>
        <el-alert v-else :closable="false">当前资源包不支持备份 ETCD 数据，请升级到最新的资源包。</el-alert>
      </div>
    </el-scrollbar>
  </div>
</template>

<script>
import BackupTask from './BackupTask.vue'

export default {
  props: {
    cluster: { type: Object, required: true },
  },
  data() {
    return {
      backups: [],
      backupToRestore: undefined,
      selection: [],
      loading: false,
    }
  },
  computed: {
  },
  components: { BackupTask },
  mounted () {
    this.list()
  },
  methods: {
    list () {
      this.loading = true
      this.kuboardSprayApi.get(`/clusters/${this.cluster.name}/backup`).then(resp => {
        this.backups = resp.data.data
        this.loading = false
      }).catch(e => {
        this.loading = false
        console.log(e)
      })
    },
    handleSelectionChange(val) {
      this.selection = val
    },
    remove_backups () {
      let req = {
        backups_to_remove: []
      }
      for (let item of this.selection) {
        req.backups_to_remove.push(`${item.node_name}/${item.etcd_member_name}/${item.backup_name}`)
      }
      this.loading = true
      this.kuboardSprayApi.post(`/clusters/${this.cluster.name}/backup/remove`, req).then(resp => {
        console.log(resp.data)
        this.selection = []
        this.list()
      }).catch(e => {
        console.log(e)
        this.loading = false
      })
    }
  }
}
</script>

<style scoped lang="scss">

</style>
