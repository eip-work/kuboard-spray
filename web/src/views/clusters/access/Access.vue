<i18n>
en:
  getKubeconfig: Fetch kubeconfig
zh:
  getKubeconfig: 获取 kubeconfig 文件
</i18n>


<template>
  <div class="access">
    <div class="app_block_title">kubeconfig</div>
    <el-button type="primary" icon="el-icon-files" @click="fetchKubeconfig">{{ $t('getKubeconfig') }}</el-button>
    <div class="app_margin_top">
      <Codemirror v-model:value="kubeconfig" :options="cmOptions"></Codemirror>
    </div>
  </div>
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
      this.kuboardSprayApi.get(`/clusters/${this.cluster.name}/access/kubeconfig`).then(resp => {
        let kubeconfig = resp.data.data
        this.kubeconfig = kubeconfig.substring(kubeconfig.indexOf("\n") + 1)
      }).catch(e => {
        console.log(e)
        this.$message.error('failed to get kubeconfig: ' + e.response.data.msg)
      })
    }
  }
}
</script>

<style scoped lang="scss">
.access {
  min-height: calc(100vh - 220px);
}
</style>
