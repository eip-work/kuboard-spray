<i18n>
en:
  sshcommon: SSH Shared Params (apply to all the k8s nodes)
zh:
  sshcommon: SSH 共享参数（适用范围：所有 k8s 节点）
</i18n>


<template>
  <div>
    <SshParams :clusterName="clusterName" :holder="inventory.all.children.k8s_cluster.vars" prop="all.children.k8s_cluster.vars" :description="$t('sshcommon')"></SshParams>
  </div>
</template>

<script>
import SshParams from './common/SshParams.vue'

export default {
  props: {
    clusterName: { type: String, required: true },
    inventory: { type: Object, required: true },
  },
  data() {
    return {

    }
  },
  computed: {
  },
  components: { SshParams },
  mounted () {
  },
  methods: {
  }
}
</script>

<style scoped lang="scss">

</style>
