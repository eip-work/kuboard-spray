<i18n>
en:
  cert_renew: Renew Certificate
  cert_renew_desc: Renew Certificates of APIServers restart APIServers one bye one.
  desc_label: Task description
zh:
  cert_renew: 更新 APIServer 证书
  cert_renew_desc: 将会逐个更新 APIServer 的证书并重启 APIServer
  desc_label: 任务描述
</i18n>

<template>
  <ClusterTask :history="cluster.history" :title="$t('cert_renew')" :cluster="cluster"
    action="renew_cert" @refresh="$emit('refresh')" :populateRequest="populateRequest">
      <el-form-item :label="$t('desc_label')">
        <span style="font-weight: bold;">{{ $t('cert_renew_desc') }}</span>
      </el-form-item>
  </ClusterTask>
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
    populateRequest () {
      return {}
    },
  }
}
</script>

<style scoped lang="css">
</style>
