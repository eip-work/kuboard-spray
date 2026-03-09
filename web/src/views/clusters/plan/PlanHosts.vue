<i18n>
en:
  singleNode: Specific Node
  enabledBation: Enabled
  disabledBation: Disabled
  selectANode: Please select a node from the diagram to the left.
  defaultSshParams: Default SSH Params
  sshcommon: SSH Shared Params (apply to all the k8s nodes)
zh:
  singleNode: 单个节点
  enabledBation: 使用跳板机
  disabledBation: 不使用跳板机
  selectANode: 请从左侧图中选择一个节点
  defaultSshParams: 默认 SSH 参数
  sshcommon: SSH 共享参数（适用范围：所有 k8s 节点）
</i18n>

<template>
  <div style="display: flex; height: 100%; overflow: visible">
    <div class="plan">
      <div class="left">
        <div style="padding: 5px; font-weight: bolder; font-size: 14px">Pangee Cluster</div>
        <div>
          <Node class="localhost" name="localhost" :cluster="cluster" hideDeleteButton
            :active="computedCurrentPropertiesTab === 'localhost'" @click="computedCurrentPropertiesTab = 'localhost'">
          </Node>
        </div>
        <div>
          <div v-if="bastionEnabled" class="verticalConnection"></div>
          <div v-else class="horizontalConnection"></div>
        </div>
        <div>
          <Node class="bastion" name="bastion" :cluster="cluster" hideDeleteButton
            :active="computedCurrentPropertiesTab === 'bastion'" @click="computedCurrentPropertiesTab = 'bastion'">
            <div style="margin-top: 10px">
              <el-tag v-if="bastionEnabled" type="danger" effect="dark" size="small"
                style="width: 100px; text-align: center">{{
                  t("enabledBation")
                }}</el-tag>
              <el-tag v-else type="info" effect="light" size="small" style="width: 100px; text-align: center">{{
                t("disabledBation")
                }}</el-tag>
            </div>
          </Node>
          <div class="horizontalConnection" :style="bastionEnabled ? '' : 'border-color: white;'"></div>
        </div>
        <div style="line-height: 28px">
          <PangeeClusterLink href="https://github.com/opencmit/pangee-cluster/blob/main/docs/guide/maintain/add-replace-node.md" :size="12">
            添加/删除节点？
          </PangeeClusterLink>
        </div>
      </div>
      <div class="right">
        <div
          style="padding: 5px; font-weight: bolder; font-size: 14px; height: 28px; line-height: 28px; margin-bottom: 10px">
          <span style="margin-right: 20px">Kubernetes Cluster</span>
          <AddNode :inventory="inventory" v-model:currentPropertiesTab="computedCurrentPropertiesTab"></AddNode>
          <el-button v-if="mode == 'view'" type="primary" plain icon="el-icon-lightning" @click="ping"
            :loading="pingpong_loading">
            <span class="app_text_mono">PING</span>
          </el-button>
        </div>
        <el-scrollbar height="calc(100vh - 343px)">
          <div class="masters">
            <Node
              v-for="(item, index) in inventory.all.children.target.children.k8s_cluster.children.kube_control_plane.hosts"
              :key="'control_plane' + index" @deleted="computedCurrentPropertiesTab = 'node_nodes'"
              @click="computedCurrentPropertiesTab = 'NODE_' + index" :pingpong="pingpong"
              :pingpong_loading="pingpong_loading" :active="nodeRoles(index)[computedCurrentPropertiesTab] ||
                'NODE_' + index === computedCurrentPropertiesTab
                " :name="index" :cluster="cluster"></Node>
            <template v-for="(item, index) in inventory.all.children.target.children.etcd.hosts" :key="'etcd' + index">
              <Node v-if="isEtcdAndNotControlPlane(index)" :name="index" :cluster="cluster"
                @deleted="computedCurrentPropertiesTab = 'node_nodes'"
                @click="computedCurrentPropertiesTab = 'NODE_' + index" @delete_button="deleteNode(index)"
                :pingpong="pingpong" :pingpong_loading="pingpong_loading" :active="nodeRoles(index)[computedCurrentPropertiesTab] ||
                  'NODE_' + index === computedCurrentPropertiesTab
                  "></Node>
            </template>
          </div>
          <div class="workers">
            <template
              v-for="(item, index) in inventory.all.children.target.children.k8s_cluster.children.kube_node.hosts"
              :key="'node' + index">
              <Node v-if="isNode(index)" :name="index" :cluster="cluster" :pingpong="pingpong"
                :pingpong_loading="pingpong_loading" @click="computedCurrentPropertiesTab = 'NODE_' + index"
                @deleted="computedCurrentPropertiesTab = 'node_nodes'" :active="nodeRoles(index)[computedCurrentPropertiesTab] ||
                  'NODE_' + index === computedCurrentPropertiesTab
                  "></Node>
            </template>
          </div>
          <div class="harbors">
            <template v-for="(item, index) in inventory.all.children.target.children.harbor.hosts" :key="'node' + index">
              <Node
                v-if="isHarborAndNothingElse(index)"
                :key="'harbor' + index"
                @deleted="computedCurrentPropertiesTab = 'node_nodes'"
                @click="computedCurrentPropertiesTab = 'NODE_' + index"
                :pingpong="pingpong"
                :pingpong_loading="pingpong_loading"
                :active="nodeRoles(index)[computedCurrentPropertiesTab] || 'NODE_' + index === computedCurrentPropertiesTab"
                :name="index"
                :cluster="cluster"
                hideDeleteButton
              ></Node>
            </template>

          </div>
          <div class="workers">
            <template v-for="(item, index) in nodeGap.inventory.all.hosts" :key="'gap_node' + index">
              <Node :name="index" :cluster="nodeGap" :pingpong="pingpong" :pingpong_loading="pingpong_loading"
                @deleted="computedCurrentPropertiesTab = 'node_nodes'"
                @click="computedCurrentPropertiesTab = 'GAP_NODE_' + index"
                :active="'GAP_NODE_' + index === computedCurrentPropertiesTab"></Node>
            </template>
          </div>
        </el-scrollbar>
      </div>
    </div>
    <div class="properties" style="height: 100%">
      <el-tabs type="card" v-model="computedCurrentPropertiesTab" class="app_scrollable_tabs"
        :before-leave="beforeLeaveCurrentTab">
        <el-tab-pane name="localhost">
          <template #label>
            {{ t("defaultSshParams") }}
          </template>
          <div v-if="computedCurrentPropertiesTab === 'localhost'">
            <SshParamsCluster :cluster="cluster" :holder="cluster.inventory.all.vars" prop="all.vars"
              :description="t('sshcommon')"></SshParamsCluster>
            <HttpProxy v-if="showHttpProxy" :cluster="cluster"></HttpProxy>
          </div>
        </el-tab-pane>
        <el-tab-pane name="bastion">
          <template #label>
            {{ $t("obj.bastion") }}
          </template>
          <SshParamsBastion v-if="cluster && cluster.inventory && computedCurrentPropertiesTab == 'bastion'"
            :cluster="cluster" nodeName="bastion" :holder="inventory.all.hosts.bastion || {}" prop="all.hosts.bastion">
          </SshParamsBastion>
        </el-tab-pane>
        <el-tab-pane :name="computedCurrentPropertiesTab"
          v-if="computedCurrentPropertiesTab.indexOf('NODE_') === 0 || computedCurrentPropertiesTab.indexOf('GAP_NODE_') === 0">
          <template #label>
            <div v-if="computedCurrentPropertiesTab.indexOf('NODE_') === 0" style="width: 100px; text-align: center">
              {{ computedCurrentPropertiesTab.slice(5) }}
            </div>
            <div v-else style="width: 100px; text-align: center">{{ computedCurrentPropertiesTab.slice(9) }}</div>
          </template>
          <div class="tab_content">
            <ConfigNode v-if="computedCurrentPropertiesTab.indexOf('NODE_') === 0" :cluster="cluster"
              :nodeName="computedCurrentPropertiesTab.slice(5)" :pingpong="pingpong" :pingpongLoading="pingpong_loading"
              @ping="ping"></ConfigNode>
            <CopyGapNodeToInventory v-else :cluster="nodeGap" :nodeName="computedCurrentPropertiesTab.slice(9)"
              :pingpong="pingpong" :pingpongLoading="pingpong_loading" @ping="ping"></CopyGapNodeToInventory>
          </div>
        </el-tab-pane>
        <el-tab-pane v-else name="node_nodes">
          <template #label>
            <div style="width: 100px; text-align: center">{{ t("singleNode") }}</div>
          </template>
          <el-scrollbar max-height="calc(100vh - 306px)">
            <div class="tab_content">
              <el-alert type="warning" :closable="false" :title="t('selectANode')"> </el-alert>
            </div>
          </el-scrollbar>
        </el-tab-pane>
      </el-tabs>
    </div>
  </div>
</template>

<script>
import Node from "./Node.vue";
import ConfigNode from "./node/ConfigNode.vue";
import CopyGapNodeToInventory from "./node/CopyGapNodeToInventory.vue";
import AddNode from "./common/AddNode.vue";
import SshParamsBastion from "./node/SshParamsBastion.vue";
import SshParamsCluster from "./common/SshParamsCluster.vue";
import HttpProxy from "./node/HttpProxy.vue";

export default {
  props: {
    cluster: { type: Object, required: true },
    mode: { type: String, required: false, default: "view" }
  },
  data() {
    return {
      currentPropertiesTab: "localhost",
      pingpong: {},
      pingpong_loading: false
    };
  },
  inject: ["validateFormFunction", "editMode"],
  computed: {
    computedCurrentPropertiesTab: {
      get() {
        return this.currentPropertiesTab;
      },
      set(v) {
        console.log("=---", v);
        if (this.editMode == "view") {
          this.currentPropertiesTab = v;
          return;
        }
        this.validateFormFunction(flag => {
          if (flag) {
            this.currentPropertiesTab = v;
          } else {
            this.$message({ message: this.$t("msg.fix_the_form"), type: "error" });
          }
        });
      }
    },
    inventory: {
      get() {
        return this.cluster.inventory;
      },
      set(v) {
        console.log(v);
      }
    },
    showHttpProxy() {
      return (
        location.hash.indexOf("showHttpProxy=true") > 0 ||
        this.cluster.inventory.all.children.target.vars.http_proxy !== undefined
      );
    },
    bastionEnabled() {
      return this.cluster.inventory.all.hosts.bastion !== undefined;
    },
    ipSet() {
      let ips = [];
      for (let i in this.cluster.inventory.all.hosts) {
        ips.push(this.cluster.inventory.all.hosts[i].ansible_host)
      }
      return ips;
    },
    nodeGap() {
      let temp = {
        name: this.cluster.name,
        type: "gap",
        inventory: {
          all: {
            hosts: {},
            children: {
              target: {
                children: {
                  calico_rr: {
                    hosts: {}
                  },
                  etcd: {
                    hosts: {}
                  },
                  k8s_cluster: {
                    children: {
                      kube_control_plane: { hosts: {} },
                      kube_node: { hosts: {} }
                    }
                  }
                },
                vars: {}
              }
            }
          }
        }
      };
      if (this.cluster.state === undefined || this.cluster.state.nodes === undefined) {
        return temp;
      }
      temp.inventory.all.children.target.vars = this.cluster.inventory.all.children.target.vars;
      for (let nodeName in this.cluster.state.nodes) {
        let node = this.cluster.state.nodes[nodeName];
        if (this.cluster.inventory.all.hosts[nodeName] === undefined && this.ipSet.indexOf(nodeName) < 0) {
          let ip = "";
          for (let i in node.status.addresses) {
            if (node.status.addresses[i].type === "InternalIP") {
              ip = node.status.addresses[i].address;
            }
          }
          temp.inventory.all.hosts[nodeName] = {
            ansible_host: ip
          };
        }
      }
      return temp;
    }
  },
  components: {
    Node,
    ConfigNode,
    AddNode,
    CopyGapNodeToInventory,
    SshParamsBastion,
    SshParamsCluster,
    HttpProxy
  },
  mounted() {
    let temp = sessionStorage.getItem(this.cluster.name + "_plan_tab") || "localhost";
    if (temp.indexOf("NODE_") === 0) {
      if (this.inventory.all.hosts[temp.slice(5)] === undefined) {
        temp = "node_nodes";
      }
    }
    this.computedCurrentPropertiesTab = temp;
  },
  watch: {
    computedCurrentPropertiesTab: function () {
      sessionStorage.setItem(this.cluster.name + "_plan_tab", this.computedCurrentPropertiesTab);
    },
    "inventory.all.hosts": function () {
      if (this.computedCurrentPropertiesTab.indexOf("NODE_") !== 0) {
        return;
      }
      if (this.inventory.all.hosts[this.computedCurrentPropertiesTab.slice(5)] === undefined) {
        this.computedCurrentPropertiesTab = "node_nodes";
      }
    }
  },
  methods: {
    validate(callback) {
      this.$refs.form.validate(callback);
    },
    beforeLeaveCurrentTab(activeName, oldActiveName) {
      if (this.editMode == "view") {
        return true;
      }
      return new Promise((resolve, reject) => {
        this.validateFormFunction(flag => {
          if (flag) {
            resolve();
          } else {
            reject();
            this.$message({ message: this.$t("msg.fix_the_form"), type: "error" });
          }
        });
      });
    },
    nodeRoles(name) {
      let roles = {};
      for (let role in this.inventory.all.children.target.children.k8s_cluster.children) {
        for (let n in this.inventory.all.children.target.children.k8s_cluster.children[role].hosts) {
          if (n === name) {
            roles[role] = true;
          }
        }
      }
      for (let n in this.inventory.all.children.target.children.etcd.hosts) {
        if (n === name) {
          roles.etcd = true;
        }
      }
      for (let n in this.inventory.all.children.target.children.harbor.hosts) {
        if (n === name) {
          roles.harbor_node = true;
        }
      }
      return roles;
    },
    isNode(name) {
      let roles = this.nodeRoles(name);
      let flag = false;
      for (let key in roles) {
        if (key === "kube_node" || key === "harbor_node") {
          flag = true;
        } else {
          return false;
        }
      }
      return flag;
    },
    isHarborAndNothingElse(name) {
      let roles = this.nodeRoles(name);
      let flag = false;
      for (let key in roles) {
        if (key === "harbor_node") {
          flag = true;
        } else {
          return false;
        }
      }
      return flag;
    },
    isEtcdAndNotControlPlane(name) {
      let roles = this.nodeRoles(name);
      return roles.etcd && !roles.kube_control_plane;
    },
    ping() {
      this.pingpong = {};
      this.pingpong_loading = true;
      let req = { nodes: "target" };
      let count = 0;
      for (let key in this.cluster.inventory.all.hosts) {
        if (key !== "localhost" && key !== "bastion") {
          count++;
        }
      }
      if (count == 0) {
        this.pingpong_loading = false;
        return;
      }
      this.pangeeClusterApi
        .post(`/clusters/${this.cluster.name}/state/ping`, req)
        .then(resp => {
          this.pingpong = resp.data.data.items;
          this.pingpong_loading = false;
        })
        .catch(e => {
          if (e.response && e.response.data) {
            this.$message.error("不能测试节点是否在线: " + e.response.data.message);
          } else {
            this.$message.error("不能测试节点是否在线: " + e);
          }
          this.pingpong_loading = false;
        });
    }
  }
};
</script>

<style lang="css">
.properties .el-tabs__item.is-active {
  padding-left: 0 20px !important;
  padding-right: 0 20px !important;
}
</style>

<style scoped lang="scss">
.plan {
  flex-grow: 2;
  display: flex;
  height: 100%;

  .left {
    width: 140px;
    min-height: 300px;
    border-radius: 10px;
    margin-left: 5px;

    .horizontalConnection {
      flex-grow: 1;
      width: calc(100% - 30px);
      height: 20px;
      border-bottom-left-radius: 5px;
      border-bottom: dashed 2px var(--el-color-primary);
      border-left: dashed 2px var(--el-color-primary);
      margin-left: 60px;
      margin-bottom: 23px;
      margin-right: -40px;
    }

    .verticalConnection {
      margin-left: 60px;
      border-left: dashed 2px var(--el-color-primary);
      height: 40px;
    }
  }

  .right {
    flex-grow: 1;
    min-height: 300px;
    border: solid 1.5px var(--el-color-primary-light-5);
    border-radius: 10px;
    margin-left: 30px;
    padding: 10px;

    .masters {
      display: flex;
      flex-wrap: wrap;
      padding-bottom: 10px;
      margin-bottom: 10px;
      border-bottom: dotted 1.5px var(--el-color-primary-light-5);
    }

    .workers {
      display: flex;
      flex-wrap: wrap;
    }

    .harbors {
      display: flex;
      flex-wrap: wrap;
    }
  }
}

.properties {
  flex-grow: 1;
  max-width: 45%;
  min-width: 45%;
  margin-left: 20px;

  .tab_content {
    padding: 10px;
    // min-height: calc(100vh - 296px);
    background-color: #f1f4fa;
  }
}
</style>
