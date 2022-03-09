<i18n>
en:
  title: Backup ETCD
  operation: Operation
  backup_desc: Backup ETCD data to kuboard-spray
zh:
  title: 备份 ETCD
  operation: 操作
  backup_desc: 备份 ETCD 数据到 kuboard-spray 所在机器
</i18n>

<template>
  <div>
    <ClusterTask action="backup_etcd" :cluster="cluster" :title="$t('title')" :populateRequest="populateRequest" @refresh="$emit('refresh')">
      <el-form-item :label="$t('operation')">
        <div style="font-weight: bolder; line-height: 28px;">{{ $t('backup_desc') }}</div>
        <el-tag class="app_text_mono">/data/cluster/{{cluster.name}}/backup</el-tag>
      </el-form-item>
    </ClusterTask>
  </div>
</template>

<script>
import ClusterTask from '../../common/task/ClusterTask.vue'

export default {
  props: {
    cluster: { type: Object, required: true },
  },
  data() {
    return {

    }
  },
  computed: {
  },
  components: { ClusterTask },
  mounted () {
  },
  methods: {
    async populateRequest (form) {
      console.log(form)
      let req = {}

      return req
    }
  }
}
</script>

<style scoped lang="scss">

</style>
