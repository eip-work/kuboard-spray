<i18n>
en:
  gapNode: Node in cluster, not in inventory
  message: Please fulfill its info here, and add to pangeecluster.
  addToPangeeCluster: Add to KubaordSpray
zh:
  gapNode: 集群中存在，而 PangeeCluster 中未记录的节点
  message: 该节点在集群中存在，但是却不存在于 PangeeCluster 记录的信息中，可能是因为您通过其他途径向集群添加了该节点。
  addToPangeeCluster: 添加到 PangeeCluster
</i18n>

<template>
  <div class="app_margin_bottom">
    <el-alert :title="t('gapNode')" type="warning" :closable="false" class="app_margin_bottom">
      <div style="width: 100%;">
        <div>{{ t('message') }}</div>
        <!-- <div style="text-align: right;">
          <el-button type="primary" icon="el-icon-plus">{{ t('addToPangeeCluster') }}</el-button>
        </div> -->
      </div>
    </el-alert>
    <StateNode v-if="nodes[nodeName]" :cluster="cluster" :nodeName="nodeName" :nodes="nodes"></StateNode>
  </div>
</template>

<script>
import StateNode from './StateNode.vue'

export default {
  props: {
    cluster: { type: Object, required: true },
    nodeName: { type: String, required: true },
    nodes: { type: Object, required: false, default: () => { return {} } },
    realCluster: { type: Object, required: false, default: undefined },
  },
  data() {
    return {

    }
  },
  computed: {
  },
  components: { StateNode },
  mounted() {
  },
  methods: {

  }
}
</script>

<style scoped lang="css"></style>
