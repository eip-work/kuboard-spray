<i18n>
en:
  sshcommon: SSH Shared Params (apply to all the k8s nodes)
zh:
  sshcommon: SSH 共享参数（适用范围：所有 k8s 节点）
</i18n>


<template>
  <div>
    <SshParamsCluster :cluster="cluster" :holder="cluster.inventory.all.children.target.vars" prop="all.children.target.vars" :description="$t('sshcommon')"></SshParamsCluster>
    <ContainerManager :cluster="cluster"></ContainerManager>
    <OsMirror :cluster="cluster"></OsMirror>
  </div>
</template>

<script>
import SshParamsCluster from '../common/SshParamsCluster.vue'
import OsMirror from './OsMirror.vue'
import ContainerManager from './ContainerManager.vue'

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
  components: { SshParamsCluster, OsMirror, ContainerManager },
  mounted () {
  },
  methods: {
  }
}
</script>

<style scoped lang="scss">

</style>
