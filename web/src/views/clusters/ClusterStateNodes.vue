<i18n>
en:
  title: K8s cluster status
  nodeCount: '{count} nodes in k8s cluster'
  unreachable: Cannot reach cluster
zh:
  title: K8S 集群状态
  nodeCount: K8S 集群中已有 {count} 个节点
  unreachable: 不能访问集群

</i18n>

<template>
  <el-popover placement="bottom-start" :title="$t('title')" :width="320" trigger="click">
    <template #reference>
      <el-button v-if="state.code === 200" type="success" round icon="el-icon-success">{{$t('nodeCount', {count: count})}}</el-button>
      <el-button v-else type="danger" round icon="el-icon-info">{{$t('unreachable')}}</el-button>
    </template>
    <div>
      <el-scrollbar max-height="45vh">
        <div v-if="state.code === 200">
          <div v-for="(node, name) in state.nodes" :key="'n' + name" class="nodeInfo app_text_mono">
            {{name}}
            <template v-for="(addr, index) in node.status.addresses" :key="name + index">
              <el-tag type="primary" v-if="addr.type === 'InternalIP'">{{addr.address}}</el-tag>
            </template>
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
    state: { type: Object, required: false, default: () => { return {} } }
  },
  data() {
    return {

    }
  },
  computed: {
    count () {
      let c = 0
      for (let k in this.state.nodes) {
        c ++
        console.log(k)
      }
      return c
    }
  },
  components: { },
  mounted () {
  },
  methods: {

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
