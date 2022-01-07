<i18n>
en:
  getKubeconfig: Fetch kubeconfig
  accessFromControlPlane: Use kubectl on control_plane
  accessMethods: You can use differenct ways to access the cluster
  controlPlanes: SSH connect to any of the following nodes, and use kubectl command to administrate the cluster.
  proposeKuboard: Kuboard is a popular K8S cluster management UI, you can refer to its website to learn how to install and use it.
  switchToPlan: Swith to cluster plan view.
zh:
  getKubeconfig: 获取 kubeconfig 文件
  accessFromControlPlane: 在主节点上使用 kubectl
  accessMethods: 您可以使用多种方式对集群进行管理控制
  controlPlanes: 您可以 ssh 到如下节点中的任意一个，直接执行 kubectl 命令可以管理集群。
  proposeKuboard: Kuboard 是一款非常值得推荐的 K8S 集群管理界面，请参考 Kuboard 网站，安装改管理界面。
  switchToPlan: 切换到集群规划页
</i18n>


<template>
  <el-skeleton v-if="cluster.state === undefined" animated></el-skeleton>
  <el-scrollbar height="calc(100vh - 220px)" v-else-if="cluster.state.code === 200">
    <el-alert type="info" :title="$t('accessMethods')" :closable="false"></el-alert>
    <div class="app_block_title">{{ $t('accessFromControlPlane') }}</div>
    <div class="access_details" v-if="cluster">
      <el-alert :title="$t('controlPlanes')" :closable="false" type="success"></el-alert>
      <div class="details">
        <template v-for="(item, key) in cluster.inventory.all.children.target.children.k8s_cluster.children.kube_control_plane.hosts" :key="key">
          <div v-if="cluster.state && cluster.state.nodes[key]" class="app_margin_top">
            <el-tag style="margin-right: 10px;" size="medium">
              <span class="app_text_mono">{{key}}</span>
            </el-tag>
            <el-tag effect="dark" size="medium">
              <span class="app_text_mono">{{cluster.inventory.all.hosts[key].ansible_host}}</span>
            </el-tag>
          </div>
        </template>
      </div>
    </div>
    <div class="app_block_title">kubeconfig</div>
    <div class="access_details">
      <el-button type="primary" icon="el-icon-files" @click="fetchKubeconfig" :loading="kubeconfigLoading">{{ $t('getKubeconfig') }}</el-button>
      <CopyToClipBoard v-if="kubeconfig" :value="kubeconfig"></CopyToClipBoard>
      <el-skeleton class="app_margin_top" v-if="kubeconfigLoading" animated></el-skeleton>
      <div v-if="kubeconfig && !kubeconfigLoading" class="app_margin_top app_codemirror_auto_height">
        <Codemirror v-model:value="kubeconfig" :options="cmOptions"></Codemirror>
      </div>
    </div>
    <div class="app_block_title">kuboard</div>
    <div class="access_details">
      <el-alert :closable="false" type="success" effect="dark" :title="$t('proposeKuboard')"></el-alert>
      <div class="details">
        <el-link class="app_margin_top" href="https://www.kuboard.cn/" target="_blank">https://www.kuboard.cn</el-link>
      </div>
    </div>
  </el-scrollbar>
  <el-alert v-else-if="cluster.state.code === 500" type="error" :closable="false" effect="dark" show-icon>
    {{cluster.state.msg}}
    <div style="margin-top: 20px;">
      <el-button type="primary" round icon="el-icon-arrow-left" @click="$emit('switch', 'plan')">{{$t('switchToPlan')}}</el-button>
    </div>
  </el-alert>
</template>

<script>
import Codemirror from "codemirror-editor-vue3"
import "codemirror/theme/darcula.css"
import "codemirror/mode/yaml/yaml.js"

export default {
  props: {
    cluster: {type: Object, required: true,}
  },
  data() {
    return {
      kubeconfig: undefined,
      kubeconfigLoading: false,
      cmOptions: {
        mode: "yaml", // Language mode
        theme: "darcula", // Theme
        lineNumbers: true, // Show line number
        lineWrapping: true,
        smartIndent: true, // Smart indent
        indentUnit: 2, // The smart indent unit is 2 spaces in length
        foldGutter: true, // Code folding
        styleActiveLine: true, // Display the style of the selected row
      },
    }
  },
  computed: {
  },
  components: { Codemirror },
  mounted () {
  },
  methods: {
    fetchKubeconfig () {
      this.kubeconfigLoading = true
      this.kubeconfig = undefined
      this.kuboardSprayApi.get(`/clusters/${this.cluster.name}/access/kubeconfig`).then(resp => {
        let kubeconfig = resp.data.data
        this.kubeconfig = kubeconfig.stdout
        this.kubeconfigLoading = false
      }).catch(e => {
        console.log(e)
        this.kubeconfigLoading = false
        this.$message.error('failed to get kubeconfig: ' + e.response.data.msg)
      })
    }
  }
}
</script>

<style scoped lang="scss">
.access_details {
  padding-left: 40px;
  margin-bottom: 20px;
}
.details {
  background-color: $--color-info-lighter;
  padding: 10px 20px 20px 20px;
}
</style>
