<i18n>
en:
  sshcommon: SSH Shared Params (apply to all the k8s nodes)
zh:
  sshcommon: SSH 共享参数（适用范围：所有 k8s 节点）
</i18n>


<template>
  <div>
    <SshParamsBastion v-if="cluster && cluster.inventory" :cluster="cluster" nodeName="bastion" :holder="inventory.all.hosts.bastion || {}" 
      prop="all.hosts.bastion"></SshParamsBastion>
    <SshParamsCluster :cluster="cluster" :holder="cluster.inventory.all.children.target.vars" prop="all.children.target.vars" :description="$t('sshcommon')"></SshParamsCluster>
    <HttpProxy v-if="showHttpProxy" :cluster="cluster"></HttpProxy>
    <ContainerManager :cluster="cluster" @refresh="$emit('refresh')"></ContainerManager>
    <OsMirror :cluster="cluster"></OsMirror>
  </div>
</template>

<script>
import SshParamsCluster from '../common/SshParamsCluster.vue'
import OsMirror from './OsMirror.vue'
import ContainerManager from './ContainerManager.vue'
import HttpProxy from './HttpProxy.vue'
import SshParamsBastion from './SshParamsBastion.vue'

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
      return location.hash.indexOf("showHttpProxy=true") > 0 || this.cluster.inventory.all.children.target.vars.http_proxy !== undefined
    },
    inventory: {
      get () {
        return this.cluster.inventory
      },
      set () {}
    },
  },
  components: { SshParamsCluster, OsMirror, ContainerManager, HttpProxy, SshParamsBastion },
  mounted () {
  },
  methods: {
  }
}
</script>

<style scoped lang="css">

</style>
