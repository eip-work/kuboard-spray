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
  <el-popover placement="bottom-start" :title="$t('title')" :width="600" trigger="click">
    <template #reference>
      <el-button v-if="state.code === undefined" type="info" round icon="el-icon-loading">{{$t('loading')}}</el-button>
      <template v-else-if="state.code === 200">
        <el-button v-if="healthy" type="success" round icon="el-icon-success">{{$t('nodeCount', { count, etcdCount: state.etcd_members_count })}}</el-button>
        <el-button v-else type="danger" round icon="el-icon-success">{{$t('nodeCount', { count, etcdCount: state.etcd_members_count })}}</el-button>
      </template>
      <el-button v-else type="danger" round icon="el-icon-info">{{$t('unreachable')}}</el-button>
    </template>
    <div>
      <el-scrollbar max-height="45vh">
        <div v-if="state.code === 200 && state.etcd_code === 200">
          <div v-for="(node, name) in onlineNodes" :key="'n' + name">
            <ClusterStateNodesItem :node="node" :name="name"></ClusterStateNodesItem>
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
import ClusterStateNodesItem from './ClusterStateNodesItem.vue'

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
    healthy () {
      for (let k in this.state.etcd_members) {
        let etcd = this.state.etcd_members[k]
        if (etcd.health.health !== true) {
          return false
        }
      }
      for (let k in this.state.nodes) {
        let node = this.state.nodes[k]
        for (let condition of node.status.conditions) {
          if (condition.type === 'Ready' && condition.status !== 'True') {
            return false
          }
          if (condition.type !== 'Ready' && condition.status !== 'False') {
            return false
          }
        }
      }
      return true
    }
  },
  components: { ClusterStateNodesItem },
  mounted () {
  },
  methods: {
  }
}
</script>

<style scoped lang="scss">
</style>
