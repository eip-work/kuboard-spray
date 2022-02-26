<i18n>
en:

zh:

</i18n>

<template>
  <div>
    <el-scrollbar height="calc(100vh - 220px)">
      <ConnectivityCheck v-if="cluster.inventory && cluster.inventory.all.children.target.children.k8s_cluster.vars.deploy_netchecker" :cluster="cluster" class="app_margin_bottom"></ConnectivityCheck>
      <CertExpirationCheck :cluster="cluster" @refresh="$emit('refresh')"></CertExpirationCheck>
    </el-scrollbar>
  </div>
</template>

<script>
import ConnectivityCheck from './ConnectivityCheck.vue'
import CertExpirationCheck from './CertExpirationCheck.vue'

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
  components: { ConnectivityCheck, CertExpirationCheck },
  mounted () {

  },
  methods: {

  }
}
</script>

<style scoped lang="css">

</style>
