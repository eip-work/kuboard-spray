<i18n>
en:
  title: K8s cluster status
  nodeCount: '{count} nodes, {etcdCount} etcd members in k8s cluster'
  unreachable: Cannot reach cluster
zh:
  title: K8S 集群状态
  nodeCount: K8S 集群中已有 {count} 个节点，{etcdCount} 个 ETCD 节点
  unreachable: 不能访问集群

</i18n>

<template>
  <el-popover placement="bottom-start" :title="$t('title')" :width="420" trigger="click">
    <template #reference>
      <el-button v-if="state.code === 200" type="success" round icon="el-icon-success">{{$t('nodeCount', { count, etcdCount: state.etcd_members_count })}}</el-button>
      <el-button v-else type="danger" round icon="el-icon-info">{{$t('unreachable')}}</el-button>
    </template>
    <div>
      <el-scrollbar max-height="45vh">
        <div v-if="state.code === 200">
          <div v-for="(node, name) in state.nodes" :key="'n' + name" class="nodeInfo app_text_mono">
            node: {{name}}
            <template v-for="(addr, index) in node.status.addresses" :key="name + index">
              <el-tag type="primary" style="margin-left: 20px;" v-if="addr.type === 'InternalIP'">{{addr.address}}</el-tag>
            </template>
          </div>
        </div>
        <div v-if="state.etcd_code === 200">
          <div v-for="(node, name) in state.etcd_members" :key="'n' + name" class="nodeInfo app_text_mono">
            etcd: {{node.name}}
            <el-tag type="primary" effect="dark" style="margin-left: 20px;">{{ node.clientURLs && node.clientURLs.length > 0 ? node.clientURLs[0] : '' }}</el-tag>
          </div>
        </div>
        <el-alert v-else type="error" :closable="false" :title="$t('unreachable')" effect="dark" show-icon>
          {{state.msg}}
        </el-alert>
      </el-scrollbar>
    </div>
  </el-popover>
</template>

<script>
export default {
  props: {
    cluster: { type: Object, required: false, default: () => { return {} } },
  },
  data() {
    return {

    }
  },
  computed: {
    state () {
      if (this.cluster && this.cluster.state) {
        return this.cluster.state
      }
      return {}
    },
    count () {
      let c = 0
      for (let k in this.state.nodes) {
        c ++
        console.log(k)
      }
      return c
    },
  },
  components: { },
  mounted () {
  },
  methods: {
    etcdMember (nodeName) {
      let host = this.cluster.inventory.all.children.target.children.etcd.hosts[nodeName]
      if (host) {
        return host.etcd_member_name
      }
      return undefined
    }
  }
}
</script>

<style scoped lang="scss">
.nodeInfo {
  padding: 10px 20px;
  margin-bottom: 10px;
  background-color: var(--el-color-success-light)
}
</style>
