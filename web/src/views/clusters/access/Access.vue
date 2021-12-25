<i18n>
en:
  getKubeconfig: Fetch kubeconfig
  accessFromControlPlane: Use kubectl on control_plane
  accessMethods: You can use differenct ways to access the cluster
zh:
  getKubeconfig: 获取 kubeconfig 文件
  accessFromControlPlane: 在主节点上使用 kubectl
  accessMethods: 您可以使用多种方式对集群进行管理控制
</i18n>


<template>
  <el-scrollbar height="calc(100vh - 220px)">
    <el-alert type="info" :title="$t('accessMethods')" :closable="false"></el-alert>
    <div class="app_block_title">{{ $t('accessFromControlPlane') }}</div>
    <div class="access_details">

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
      
    </div>
  </el-scrollbar>
</template>

<script>
import Codemirror from "codemirror-editor-vue3"
import "codemirror/theme/dracula.css"
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
        mode: "text/yaml", // Language mode
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
        this.kubeconfig = kubeconfig.substring(kubeconfig.indexOf("\n") + 1)
        this.kubeconfigLoading = false
      }).catch(e => {
        console.log(e)
        this.$message.error('failed to get kubeconfig: ' + e.response.data.msg)
        this.kubeconfigLoading = false
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
</style>
