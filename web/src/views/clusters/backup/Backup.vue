<i18n>
en:
  node: Node name
  etcd_member_name: ETCD member name
  backup_name: Backup name
  size: Size
zh:
  node: 节点名称
  etcd_member_name: ETCD 成员名称
  backup_name: 备份名称
  size: 备份文件大小
</i18n>

<template>
  <div>
    <el-scrollbar height="calc(100vh - 220px)">
      <div v-if="cluster">
        <template v-if="cluster.resourcePackage && cluster.resourcePackage.data.supported_playbooks.backup_etcd">
          <BackupTask :cluster="cluster"></BackupTask>
          <div class="app_margin_bottom"></div>
          <el-alert v-if="backups.length === 0" type="warning" :closable="false">
            当前没有备份
          </el-alert>
          <div v-else>
            <el-table :data="backups" style="width: 100%" @selection-change="handleSelectionChange">
              <el-table-column type="selection" width="55"></el-table-column>
              <el-table-column prop="node_name" :label="$t('node_name')" sortable />
              <el-table-column prop="etcd_member_name" :label="$t('etcd_member_name')" sortable />
              <el-table-column prop="backup_name" :label="$t('backup_name')" sortable />
              <el-table-column prop="size" :label="$t('size')" sortable />
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
      this.kuboardSprayApi.get(`/clusters/${this.cluster.name}/backup`).then(resp => {
        this.backups = resp.data.data
      }).catch(e => {
        console.log(e)
      })
    },
    handleSelectionChange(val) {
      this.selection = val
    }
  }
}
</script>

<style scoped lang="scss">

</style>
