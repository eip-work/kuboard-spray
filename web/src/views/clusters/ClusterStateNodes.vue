<i18n>
en:
  title: Nodes in cluster
  nodeCount: '{count} nodes, {etcdCount} etcd members in k8s cluster'
  unreachable: Cannot reach cluster
  loading: Loading
zh:
  title: 集群中的节点
  nodeCount: K8S 集群中已有 {count} 个节点，ETCD 集群中已有 {etcdCount} 个成员
  unreachable: 不能访问集群
  loading: 加载中
</i18n>

<template>
  <el-popover placement="bottom-start" :title="$t('title')" :width="480" trigger="click">
    <template #reference>
      <el-button v-if="state.code === undefined" type="success" round icon="el-icon-loading">{{$t('loading')}}</el-button>
      <el-button v-else-if="state.code === 200" type="success" round icon="el-icon-success">{{$t('nodeCount', { count, etcdCount: state.etcd_members_count })}}</el-button>
      <el-button v-else type="danger" round icon="el-icon-info">{{$t('unreachable')}}</el-button>
    </template>
    <div>
      <el-scrollbar max-height="45vh">
        <div v-if="state.code === 200 && state.etcd_code === 200">
          <div v-for="(node, name) in onlineNodes" :key="'n' + name" class="nodeInfo app_text_mono">
            <div style="min-width: 120px;">
              {{name}}
            </div>
            <div style="flex-grow: 1;">
              <template v-if="node.k8s_node">
                <el-tag type="primary" effect="dark" style="margin-left: 10px;" v-if="node.k8s_node.metadata.labels['node-role.kubernetes.io/control-plane'] !== undefined">控制节点</el-tag>
                <el-tag type="success" effect="dark" style="margin-left: 10px;" v-if="isKubeNode(name)">工作节点</el-tag>
              </template>
              <template v-if="node.etcd_member">
                <el-tag type="warning" effect="dark" style="margin-left: 10px;">{{node.etcd_member.name}}</el-tag>
                <!-- <el-tag type="primary" effect="light" style="margin-left: 20px;">{{ node.etcd_member.clientURLs && node.etcd_member.clientURLs.length > 0 ? node.etcd_member.clientURLs[0] : '' }}</el-tag> -->
              </template>
              <template v-if="node.k8s_node">
                <template v-for="(addr, index) in node.k8s_node.status.addresses" :key="name + index">
                  <el-tag type="primary" style="float: right;" v-if="addr.type === 'InternalIP'">{{addr.address}}</el-tag>
                </template>
              </template>
            </div>
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
  inject: ['onlineNodes'],
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
    },
    isKubeNode (nodeName) {
      let node = this.onlineNodes[nodeName]
      if (node !== undefined && node.k8s_node !== undefined) {
        node = node.k8s_node
        for (let key in node.spec.taints) {
          let taint = node.spec.taints[key]
          if (taint.effect === "NoSchedule" && taint.key === "node-role.kubernetes.io/master") {
            return false
          }
        }
        return true
      }
      return false
    }
  }
}
</script>

<style scoped lang="scss">
.nodeInfo {
  padding: 10px 20px;
  margin-bottom: 10px;
  background-color: var(--el-color-info-light);
  display: flex;
}
</style>
