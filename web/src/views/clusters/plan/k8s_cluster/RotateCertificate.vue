<i18n>
en:
  label: Certificate
  description: Auto rotate certificates
zh:
  label: 证书更新
  description: 自动更新证书
</i18n>

<template>
  <ConfigSection v-model:enabled="enabled" :label="$t('label')" :description="$t('description')" disabled>
    <template #more>
      是否自动更新 kubelet / kube-apiserver 证书
    </template>
    <FieldBool labelWidth="180px" :holder="vars" :prop="prop" fieldName="auto_renew_certificates">
      <template #edit_desc>
        <div v-if="vars.auto_renew_certificates">
          每月第一个星期一自动更新 apiserver 证书
        </div>
      </template>
      <template #view_desc>
        <div v-if="vars.auto_renew_certificates">
          每月第一个星期一自动更新 apiserver 证书
        </div>
      </template>
    </FieldBool>
    <FieldBool labelWidth="180px" :holder="vars" :prop="prop" fieldName="kubelet_rotate_certificates">
      <template #edit_desc>
        <div v-if="vars.kubelet_rotate_certificates">
          kubelect client 证书过期时自动更新
        </div>
      </template>
      <template #view_desc>
        <div v-if="vars.kubelet_rotate_certificates">
          kubelect client 证书过期时自动更新
        </div>
      </template>
    </FieldBool>
    <!-- <FieldBool labelWidth="180px" :holder="vars" :prop="prop" fieldName="kubelet_rotate_server_certificates">
      <template #edit_desc>
        <div v-if="vars.kubelet_rotate_server_certificates" style="color: var(--el-color-danger);">
          自动申请证书后，需要您手动审批证书更新申请
          <li>kubectl get csr</li>
          <li>kubectl certificate approve</li>
        </div>
        <div v-else style="color: var(--el-color-danger);">
          将触发 CIS 扫描告警
        </div>
      </template>
      <template #view_desc>
        <div v-if="vars.kubelet_rotate_server_certificates" style="color: var(--el-color-danger);">
          自动申请证书后，需要您手动审批证书更新申请
          <li>kubectl get csr</li>
          <li>kubectl certificate approve</li>
        </div>
      </template>
    </FieldBool> -->
  </ConfigSection>
</template>

<script>
export default {
  props: {
    cluster: { type: Object, required: true },
  },
  data() {
    return {

    }
  },
  computed: {
    prop () {
      return 'all.children.target.children.k8s_cluster.vars'
    },
    vars: {
      get () { return this.cluster.inventory.all.children.target.children.k8s_cluster.vars },
      set () {}
    },
    enabled: {
      get () {return true},
      set () {},
    },
  },
  components: { },
  mounted () {
  },
  methods: {
  }
}
</script>

<style scoped lang="css">

</style>
