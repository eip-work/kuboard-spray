<template>
  <div>
    <el-dialog v-model="dialogVisibleCompute" :title="title" width="70%" top="5vh">
      <div class="preview_yaml_code_mirror">
        <Codemirror v-if="dialogVisibleCompute && yamlVisible" v-model:value="content" :options="cmOptions"></Codemirror>
      </div>
      <template #footer>
        <div style="margin-top: -10px;">
          <CopyToClipBoard v-if="yamlVisible" style="float: left;" :value="content"></CopyToClipBoard>
          <el-button type="primary" icon="el-icon-close" @click="dialogVisibleCompute = false">{{ $t('msg.close') }}</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import Codemirror from "codemirror-editor-vue3"
import "codemirror/theme/darcula.css"
import "codemirror/mode/yaml/yaml.js"
import yaml from 'js-yaml'

export default {
  props: {
  },
  data() {
    return {
      dialogVisible: false,
      yamlVisible: false,
      items: [],
      title: this.$t('msg.preview_yaml'),
      cmOptions: {
        mode: "yaml", // Language mode
        theme: "darcula", // Theme
        lineNumbers: true, // Show line number
        lineWrapping: true,
        smartIndent: true, // Smart indent
        indentUnit: 2, // The smart indent unit is 2 spaces in length
        foldGutter: true, // Code folding
        readOnly: true,
        styleActiveLine: true, // Display the style of the selected row
      },
    }
  },
  computed: {
    content () {
      return yaml.dump(this.items)
    },
    dialogVisibleCompute: {
      get () { 
        return this.dialogVisible
      },
      set (v) {
        this.dialogVisible = v
        if (!v) {
          this.yamlVisible = v
        }
        this.items = []
      }
    }
  },
  components: { Codemirror },
  mounted () {
  },
  methods: {
    show(items, title) {
      this.items = items
      this.title = title || this.title
      this.dialogVisible = true
      setTimeout(() => {
        this.yamlVisible = true
      }, 500)
    }
  }
}
</script>

<style lang="css">
.preview_yaml_code_mirror {
  height: calc(90vh - 150px);
}
.preview_yaml_code_mirror .CodeMirror {
  height: calc(90vh - 150px);
}

.preview_yaml_code_mirror .CodeMirror-scroll {
  height: calc(90vh - 150px);
  overflow-y: hidden;
  overflow-x: auto;
}
</style>
