<i18n>
en:
  cert_expiration_check: Cert expiration check
  cert_expiration_check_desc: Check certificate expiration date of APIServer
zh:
  cert_expiration_check: 证书有效期检查
  cert_expiration_check_desc: 检查 APIServer 的证书有效期
</i18n>

<template>
  <div>
    <div class="app_block_title">{{ $t('cert_expiration_check') }}</div>
    <div class="app_description">
      {{ $t('cert_expiration_check_desc') }}
      <el-button type="primary" style="margin-left: 20px;" icon="el-icon-promotion" @click="checkCertExpiration">{{$t('cert_expiration_check')}}</el-button>
      <CertRenew v-if="cluster.resourcePackage.data.supported_playbooks.renew_cert && !cluster.history.processing" :cluster="cluster" @refresh="$emit('refresh')"></CertRenew>
    </div>
    <el-skeleton v-if="loading" animated class="app_margin_top"></el-skeleton>
    <div v-else-if="checkResult">
      <el-tabs v-model="activeName" type="card">
        <el-tab-pane v-for="(result, nodeName) in certExpiration" :key="'cert' + nodeName" :label="nodeName">
          <el-scrollbar class="result" :max-height="285">
            <pre class="app_code"><span style="color: var(--el-color-success); font-weight: bolder;">[root@{{ nodeName }}]</span># kubeadm certs check-expiration

{{result.stdout}}</pre>
          </el-scrollbar>
        </el-tab-pane>
      </el-tabs>
    </div>
  </div>
</template>

<script>
import CertRenew from './CertRenew.vue'

export default {
  props: {
    cluster: { type: Object, required: true },
  },
  data() {
    return {
      checkResult: undefined,
      loading: false,
      activeName: undefined,
    }
  },
  computed: {
    certExpiration () {
      if (this.checkResult) {
        return this.checkResult.plays[0].tasks[0].hosts
      }
      return {}
    }
  },
  components: { CertRenew },
  mounted () {
  },
  methods: {
    checkCertExpiration () {
      this.loading = true
      this.kuboardSprayApi.get(`/clusters/${this.cluster.name}/state/check_cert_expiration`).then(resp => {
        this.checkResult = resp.data.data
        this.loading = false
      }).catch(e => {
        console.error(e)
        this.loading = false
        if (e.response && e.response.data.message)
        this.$message.error(e.response.data.message)
      })
    },
    renewCert () {

    }
  }
}
</script>

<style scoped lang="scss">

</style>
