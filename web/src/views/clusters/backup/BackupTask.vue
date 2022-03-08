<i18n>
en:
  title: Backup ETCD
zh:
  title: 备份 ETCD

</i18n>

<template>
  <div>
    <ClusterTask action="backup_etcd" :cluster="cluster" :title="$t('title')" :populateRequest="populateRequest" @refresh="$emit('refresh')">
      <div>备份 ETCD </div>
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
    populateRequest (form) {
      console.log(form)
      let req = {}

      return req
    }
  }
}
</script>

<style scoped lang="scss">

</style>
