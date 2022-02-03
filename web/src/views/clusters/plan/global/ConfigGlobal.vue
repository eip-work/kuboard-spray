<i18n>
en:
  sshcommon: SSH Shared Params (apply to all the k8s nodes)
zh:
  sshcommon: SSH 共享参数（适用范围：所有 k8s 节点）
</i18n>


<template>
  <div>
    <SshParamsCluster :cluster="cluster" :holder="cluster.inventory.all.children.target.vars" prop="all.children.target.vars" :description="$t('sshcommon')"></SshParamsCluster>
    <HttpProxy v-if="showHttpProxy" :cluster="cluster"></HttpProxy>
    <ContainerManager :cluster="cluster"></ContainerManager>
    <OsMirror :cluster="cluster"></OsMirror>
  </div>
</template>

<script>
import SshParamsCluster from '../common/SshParamsCluster.vue'
import OsMirror from './OsMirror.vue'
import ContainerManager from './ContainerManager.vue'
import HttpProxy from './HttpProxy.vue'

export default {
  props: {
    cluster: { type: Object, required: true },
  },
  data() {
    return {

    }
  },
  computed: {
    showHttpProxy () {
      return location.hash.indexOf("showHttpProxy=true") > 0
    }
  },
  components: { SshParamsCluster, OsMirror, ContainerManager, HttpProxy },
  mounted () {
  },
  methods: {
  }
}
</script>

<style scoped lang="css">

</style>
