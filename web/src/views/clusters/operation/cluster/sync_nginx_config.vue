<i18n>
en:
  title: "Update apiserver list in loadbalancer"
  sync_nginx_config_desc: "Control_plane list in k8s cluster changed, this task is going to update 'upstream kube_apiserver' block in /etc/nginx/nginx.conf on all the kube_node nodes."
zh:
  title: "更新负载均衡中 apiserver 列表"
  sync_nginx_config_desc: "集群中控制节点的列表发生变化，此操作将更新所有工作节点上 /etc/nginx/nginx.conf 文件中 upstream kube_apiserver 的列表"
</i18n>

<template>
  <el-form-item prop="action" :rules="rules">
    <div class="form_description" style="margin-bottom: 10px;">
      {{ $t('sync_nginx_config_desc') }}
    </div>
    <template v-for="(item, key) in cluster.inventory.all.children.target.children.k8s_cluster.children.kube_node.hosts" :key="'node' + key">
      <el-tag v-if="pingpong[key] && pingpong[key].status === 'SUCCESS'
          && cluster.inventory.all.children.target.children.k8s_cluster.children.kube_node.hosts[key] !== undefined
          && cluster.inventory.all.hosts[key].kuboardspray_node_action !== 'add_node'" 
          style="margin-right: 10px; margin-bottom: 10px;" effect="dark" type="primary">
        <span class="app_text_mono">{{key}}</span>
      </el-tag>
    </template>
  </el-form-item>
</template>

<script>
export default {
  props: {
    cluster: { type: Object, required: true },
    name: { type: String, required: true },
    form: { type: Object, required: true },
    pingpong: { type: Object, required: false, default: () => { return {} } },
    pingpong_loading: { type: Boolean, required: false, default: false },
    offlineNodes: { type: Array, required: false, default: () => { return [] } },
  },
  data() {
    return {
      rules: [{
        validator (rule, value, callback) {
          return callback()
        },
        trigger: 'change',
      }]
    }
  },
  inject: ['isClusterInstalled', 'isClusterOnline', 'pendingRemoveNodes', 'pendingAddNodes'],
  computed: {
  },
  components: { },
  mounted () {
  },
  methods: {
    populateRequest() {
      return {}
    },
  }
}
</script>

<style scoped lang="scss">
.form_description {
  font-size: 12px;
  // color: var(--el-text-color-placeholder);
  color: #aaa;
  max-width: 700px;
}
</style>
