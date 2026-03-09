<i18n>
en:
  singleNode: Specific Node
  global_config: Global Config
  addons: Addons
  enabledBation: Enabled
  disabledBation: Disabled
  selectANode: Please select a node from the diagram to the left.
  resourcePackage: Resource Package
  add_node_count: "{count} nodes to add"
  delete_node_count: "{count} nodes to delete"
  edit: Edit
zh:
  singleNode: 单个节点
  global_config: 全局设置
  addons: 可选组件
  enabledBation: 使用跳板机
  disabledBation: 不使用跳板机
  selectANode: 请从左侧图中选择一个节点
  resourcePackage: 资源包
  add_node_count: 共 {count} 个待添加节点
  delete_node_count: 共 {count} 个待删除节点
  edit: 编辑
</i18n>

<template>
  <div style="height: calc(100% - 2px); display: flex; flex-direction: column">
    <div style="margin-bottom: 10px">
      <el-radio-group v-model="currentOperation" size="default">
        <template v-for="(operation, index) in cluster.resourcePackage.data.operations" :key="'op_' + index">
          <el-radio-button :label="operation.title[locale]" :value="index"
            :disabled="isDisabled(operation)"></el-radio-button>
        </template>
      </el-radio-group>
      <div class="operation-params"
        v-if="cluster.resourcePackage.data.operations[currentOperation] && cluster.resourcePackage.data.operations[currentOperation].pangeecluster_node_action">
        {{ t(cluster.resourcePackage.data.operations[currentOperation].pangeecluster_node_action + "_count", {
          count:
        pendingNodesLength}) }}
        <el-link type="primary" icon="el-icon-edit" @click="$emit('go-to-plan-hosts')">
          <span style="margin-left: 5px">{{ t("edit") }}</span>
        </el-link>
        <OperationPendingNode v-for="(node, name) in pendingNodes" :node="node" :node-name="name"
          :inventory="cluster.inventory">
        </OperationPendingNode>
      </div>
    </div>
    <div style="display: flex; gap: 10px; flex: 1; overflow-y: hidden;">
      <div class="operation-card" style="min-width: 240px; flex-grow: 0;">
        <div class="noselect operation-steps">
          <div>
            <el-steps direction="vertical" :active="currentStep">
              <template v-for="(step, index) in cluster.resourcePackage.data.operations[currentOperation].steps"
                :key="'step_' + index">
                <el-step :title="step.name" :description="step.title[locale]" @click="currentStep = index"
                  :status="currentStep >= index || true ? 'success' : ''"
                  :class="currentStep == index ? 'is-selected-step' : ''" />
              </template>
            </el-steps>
          </div>
        </div>
      </div>
      <div class="operation-card" style="max-width: 50%; min-width: 50%;" :body-style="{ height: 'calc(100% - 40px)' }">
        <div class="markdown-title" v-if="cluster.resourcePackage.data.operations[currentOperation].steps[currentStep]">
          <div style="flex-grow: 1; font-weight: bolder;">
            {{ cluster.resourcePackage.data.operations[currentOperation].steps[currentStep].title[locale] }}
          </div>
          <el-button style="float: right" @click="showFileBrowser" type="primary"
            icon="el-icon-pointer">查看代码</el-button>
        </div>
        <div class="markdown-content">
          <OperationStepMarkdown :cluster="cluster" :operation-index="currentOperation" :step-index="currentStep">
          </OperationStepMarkdown>
        </div>
      </div>
      <div class="operation-card">
        <div class="markdown-title">
          <OperationStepExecute :cluster="cluster" :currentOperation="currentOperation" :currentStep="currentStep"
            @refresh="$emit('refresh')">
          </OperationStepExecute>
          <OperationStepStatus v-if="!loadingStepHistory && $refs.history && $refs.history.history && $refs.history.history.length > 0"
            ref="stepStatus" :cluster="cluster" :currentOperation="currentOperation" :currentStep="currentStep"
            @refresh="$refs.history.refresh()"></OperationStepStatus>
        </div>
        <div class="operation-history">
          <OperationStepHistory ref="history" :cluster="cluster" :currentOperation="currentOperation"
            :currentStep="currentStep" @loadingStatusChanged="loadingStepHistory = $event">
          </OperationStepHistory>
        </div>
      </div>
    </div>
    <FileBrowser :package-name="`${cluster.resourcePackage.metadata.version}`" ref="filebrowser"></FileBrowser>
  </div>
</template>

<script>
import FileBrowser from "./filebrowser/FileBrowser.vue";
import OperationStepMarkdown from "./OperationStepMarkdown.vue"
import OperationStepExecute from "./OperationStepExecute.vue";
import OperationStepHistory from "./OperationStepHistory.vue";
import OperationStepStatus from "./OperationStepStatus.vue";
import OperationPendingNode from "./OperationPendingNode.vue";

function getValue(obj, path, defaultValue = undefined) {
  // 将路径字符串转换为数组
  const keys = Array.isArray(path) ? path : path.split('.');

  // 遍历路径
  let result = obj;
  for (const key of keys) {
    if (result == null || result[key] === undefined) {
      return defaultValue;
    }
    result = result[key];
  }
  return result;
}

export default {
  props: {
    cluster: { type: Object, required: true },
  },
  data() {
    return {
      loadingStepHistory: false
    };
  },
  computed: {
    currentOperation: {
      get() {
        if (this.$store.state.cluster[this.cluster.name] == undefined) {
          return 0;
        }
        if (this.$store.state.cluster[this.cluster.name].operation == undefined) {
          return 0;
        }
        return this.$store.state.cluster[this.cluster.name].operation.currentOperation || 0
      },
      set(v) {
        this.$store.commit("cluster/CHANGE_CLUSTER_STATE",
          {
            cluster: this.cluster.name,
            tab: "operation",
            key: "currentOperation",
            value: v
          }
        );
      }
    },
    currentStep: {
      get() {
        if (this.$store.state.cluster[this.cluster.name] == undefined) {
          return 0;
        }
        if (this.$store.state.cluster[this.cluster.name].operation == undefined) {
          return 0;
        }
        let t = this.$store.state.cluster[this.cluster.name].operation.currentStep;
        if (this.cluster.resourcePackage.data.operations[this.currentOperation].steps[t] == undefined) {
          return 0;
        }
        return this.$store.state.cluster[this.cluster.name].operation.currentStep || 0
      },
      set(v) {
        this.$store.commit("cluster/CHANGE_CLUSTER_STATE",
          {
            cluster: this.cluster.name,
            tab: "operation",
            key: "currentStep",
            value: v
          }
        );
      }
    },
    pendingNodes() {
      let result = {};
      if (this.cluster.resourcePackage.data.operations[this.currentOperation] && this.cluster.resourcePackage.data.operations[this.currentOperation].pangeecluster_node_action) {
        let node_action = this.cluster.resourcePackage.data.operations[this.currentOperation].pangeecluster_node_action;
        for (let i in this.cluster.inventory.all.hosts) {
          let host = this.cluster.inventory.all.hosts[i]
          if (host.pangeecluster_node_action == node_action) {
            result[i] = host;
          }
        }
      }
      return result;
    },
    pendingNodesLength() {
      let length = 0;
      for (let i in this.pendingNodes) {
        length++;
      }
      return length;
    }
  },
  components: {
    FileBrowser,
    OperationStepExecute,
    OperationStepMarkdown,
    OperationStepHistory,
    OperationStepStatus,
    OperationPendingNode
  },
  methods: {
    stepClass(step) {
      if (this.currentStep == step) {
        return "step active";
      }
      return "step";
    },
    showFileBrowser() {
      let path = "/operations/" + this.cluster.resourcePackage.data.operations[this.currentOperation].name;
      path += "/" + this.cluster.resourcePackage.data.operations[this.currentOperation].steps[this.currentStep].name;
      this.$refs.filebrowser.show([{
        isDir: true,
        path: path
      }]);
    },
    isDisabled(operation) {
      console.log(operation);
      if (operation.enabled_on == undefined) {
        return false;
      }
      for (let condition of operation.enabled_on) {
        let variableValue = getValue(this.cluster.inventory, condition.variable);
        console.log(variableValue, condition.operator, condition.value)
        if (variableValue == undefined) {
          return true;
        }
        if (condition.operator == "equal") {
          if (variableValue != condition.value) {
            return true;
          }
        } else if (condition.operator == "not_equal") {
          if (variableValue == condition.value) {
            return true;
          }
        }
      }
      return false;
    }
  }
};
</script>

<style lang="scss" scoped>
.operation-params {
  padding: 10px;
  border: 1px solid var(--el-color-info-light-9);
  border-radius: 5px;
  border-top-left-radius: 0;
  background-color: var(--el-fill-color-extra-light);
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
  align-items: center;
}

.operation-card {
  flex-grow: 1;
  height: 100%;
  overflow: hidden;
  padding: 20px;
  border: 1px solid var(--el-border-color);
  border-radius: 8px;

  .operation-steps {
    height: 100%;
    overflow: hidden;
    overflow-y: auto;
  }

  .markdown-title {
    height: 36px;
    margin-bottom: 10px;
    background-color: var(--el-color-primary-light-9);
    display: flex;
    align-items: center;
    justify-content: stretch;
    padding: 0 10px;
  }

  .markdown-content {
    height: calc(100% - 46px);
    overflow: hidden;
    overflow-y: auto;
  }

  .operation-history {
    height: calc(100% - 46px);
    overflow: hidden;
    overflow-y: auto;
  }
}
</style>

<style lang="scss">
.operation-card .el-step__main {
  margin-bottom: 20px;
  cursor: pointer;

  .el-step__title {
    font-size: 13px;
    font-family: Consolas, Menlo, "Bitstream Vera Sans Mono", Monaco, "微软雅黑", monospace;
    font-weight: normal;
  }

  .el-step__description {
    border-bottom: 2px solid var(--el-color-primary-light-7);
    font-size: 14px;
    font-weight: bolder;
  }
}

.operation-card .is-selected-step {
  .el-step__icon {
    border: 2px solid var(--el-color-primary);

    svg {
      color: var(--el-color-primary);
    }
  }

  .el-step__description {
    border-bottom: 3px solid var(--el-color-primary);
  }
}
</style>