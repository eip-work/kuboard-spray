<i18n>
en:
  architecture: Architecture
  global_config: Configuration
  remote_global_config: Configuration
  hosts: Hosts
zh:
  architecture: 部署架构
  global_config: 参数配置
  remote_global_config: 参数配置
  hosts: 节点列表
</i18n>

<template>
  <el-form ref="form" label-width="120px" label-position="left" @submit.enter.prevent :model="inventory"
    style="height: 100%">
    <div style="display: flex; height: 100%">
      <el-tabs v-model="currentTab" type="border-card" tab-position="left" class="app_scrollable_tabs"
        style="width: 100%; margin-bottom: -5px">
        <el-tab-pane name="architecture">
          <template #label> {{ t("architecture") }} </template>
          <PlanArchitecture :cluster="cluster"></PlanArchitecture>
        </el-tab-pane>
        <el-tab-pane name="config_remote">
          <template #label> {{ t("remote_global_config") }} </template>
          <PlanConfigDynamic :cluster="cluster" class="app_scroll_content"></PlanConfigDynamic>
        </el-tab-pane>
        <el-tab-pane name="hosts" style="overflow: visible">
          <template #label> {{ t("hosts") }} </template>
          <PlanHosts :cluster="cluster" :mode="mode"></PlanHosts>
        </el-tab-pane>
      </el-tabs>
    </div>
  </el-form>
</template>

<script>
import { computed } from "vue";
import PlanConfigDynamic from "./PlanConfigDynamic.vue";
import PlanHosts from "./PlanHosts.vue";
import PlanArchitecture from "./PlanArchitecture.vue";

export default {
  props: {
    cluster: { type: Object, required: true },
    mode: { type: String, required: false, default: "view" }
  },
  data() {
    return {
      pingpong: {},
      pingpong_loading: false
    };
  },
  provide() {
    return {
      editMode: computed(() => {
        if (this.mode === "view") {
          return "view";
        }
        if (this.cluster && this.cluster.history.success_tasks.length > 0) {
          return "frozen";
        }
        return this.mode;
      }),
      currentTab: computed(() => {
        return this.currentPropertiesTab;
      }),
      validateFormFunction: callback => {
        this.$refs.form.validate(flag => {
          console.log("validateForm: " + flag);
          if (callback) {
            callback(flag);
          }
        });
      }
    };
  },
  computed: {
    currentTab: {
      get() {
        if (this.$store.state.cluster[this.cluster.name] == undefined) {
          return "architecture";
        }
        if (this.$store.state.cluster[this.cluster.name].plan == undefined) {
          return "architecture";
        }
        return this.$store.state.cluster[this.cluster.name].plan.currentTab || "architecture"
      },
      set(v) {
        this.$store.commit("cluster/CHANGE_CLUSTER_STATE",
          {
            cluster: this.cluster.name,
            tab: "plan",
            key: "currentTab",
            value: v
          }
        );
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
    bastionEnabled() {
      return this.cluster.inventory.all.children.target.children.bastion !== undefined;
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
        if (this.cluster.inventory.all.hosts[nodeName] === undefined) {
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
    PlanArchitecture,
    PlanConfigDynamic,
    PlanHosts
  },
  mounted() {
    let temp = sessionStorage.getItem(this.cluster.name + "_plan_tab") || "global_config";
    if (temp.indexOf("NODE_") === 0) {
      if (this.inventory.all.hosts[temp.slice(5)] === undefined) {
        temp = "node_nodes";
      }
    }
    this.currentPropertiesTab = temp;
  },
  watch: {
    currentPropertiesTab: function () {
      sessionStorage.setItem(this.cluster.name + "_plan_tab", this.currentPropertiesTab);
    },
    "inventory.all.hosts": function () {
      if (this.currentPropertiesTab.indexOf("NODE_") !== 0) {
        return;
      }
      if (this.inventory.all.hosts[this.currentPropertiesTab.slice(5)] === undefined) {
        this.currentPropertiesTab = "node_nodes";
      }
    }
  },
  methods: {
    validate(callback) {
      this.$refs.form.validate(callback);
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
      return roles;
    },
    isNode(name) {
      let roles = this.nodeRoles(name);
      let flag = false;
      for (let key in roles) {
        if (key === "kube_node") {
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
    showBastion() {
      this.currentPropertiesTab = "global_config";
      setTimeout(() => {
        this.$refs.configPangeeClusterScroll.setScrollTop(0);
      }, 400);
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
  }
}

.properties {
  flex-grow: 1;
  max-width: 45%;
  min-width: 45%;
  margin-left: 20px;

  .tab_content {
    padding: 10px;
    min-height: calc(100vh - 296px);
    background-color: #f1f4fa;
  }
}
</style>
