<i18n>
en:
  title: Restore ETCD
  operation: Operation
  restore_desc: Restore backup data to ETCD cluster. You should not continue unless you know what you are doing.
zh:
  title: 恢复 ETCD 节点
  operation: 操作
  restore_desc: 将备份数据恢复到 ETCD 集群。请谨慎操作！此操作将会把 ETCD 集群所有节点的数据恢复到选定的快照。
</i18n>

<template>
  <div>
    <ClusterTask v-if="!cluster.history.processing"
      action="restore_etcd" :cluster="cluster" :title="$t('title')" :populateRequest="populateRequest"
      :disabled="disabled" type="primary" @refresh="$emit('refresh')" :width="600">
      <el-form-item :label="$t('operation')">
        <el-alert type="warning" show-icon :closable="false" style="font-weight: bolder; line-height: 28px; margin: 10px 0 20px 0;">
          {{ $t('restore_desc') }}
          <kuboard-spray-link href="https://kuboard-spray.cn/guide/maintain/backup.html#影响" style="float: right;" size="12px;" type="warning"></kuboard-spray-link>
        </el-alert>
        已选定从 ETCD 成员 <el-tag class="app_text_mono">{{backupFile.etcd_member_name}}</el-tag> 备份的快照文件 <el-tag class="app_text_mono">{{backupFile.backup_name}}</el-tag>
      </el-form-item>
    </ClusterTask>
  </div>
</template>

<script>
import KuboardSprayLink from '../../../components/KuboardSprayLink.vue'
import ClusterTask from '../../common/task/ClusterTask.vue'

export default {
  props: {
    cluster: { type: Object, required: true },
    disabled: { type: Boolean, required: false, default: false },
    backupFile: { type: Object, required: false, default: undefined },
  },
  data() {
    return {

    }
  },
  computed: {
  },
  components: { ClusterTask, KuboardSprayLink },
  emits: ['refresh'],
  mounted () {
  },
  methods: {
    async populateRequest () {
      let req = {
        backup_file_path: `${this.backupFile.node_name}/${this.backupFile.etcd_member_name}`,
        backup_file_name: this.backupFile.backup_name
      }

      return req
    }
  }
}
</script>

<style scoped lang="scss">

</style>
