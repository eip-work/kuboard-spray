<i18n>
en:
  kube_control_plane: control plane
  kube_node: worker node
  etcd: etcd node
zh:
  kube_control_plane: 控制节点
  kube_node: 工作节点
  etcd: etcd
</i18n>

<template>
  <div :class="active ? 'node active' : 'node'">
    <div class="app_text_mono" style="font-weight: bold;">
      {{ name }}
    </div>
    <div class="app_text_mono" v-if="inventory.all.hosts[name]">
      {{ inventory.all.hosts[name].ansible_host }}
    </div>
    <div :class="role + ' role app_text_mono'" v-for="(_, role) in roles" :key="'r' + role">{{$t(role)}}</div>
  </div>
</template>

<script>

export default {
  props: {
    name: { type: String, required: true },
    inventory: { type: Object, required: true },
    active: { type: Boolean, required: false, default: false },
  },
  computed: {
    roles () {
      let result = {}
      for (let role in this.inventory.all.children.k8s_cluster.children) {
        for (let n in this.inventory.all.children.k8s_cluster.children[role].hosts) {
          if (n === this.name) {
            result[role] = true
          }
        }
      }
      return result
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
  .node {
    border: solid 1px $--border-color-light;
    border-radius: 5px;
    width: 120px;
    margin: 5px;
    padding: 15px;
    font-size: 13px;
    cursor: pointer;
    .role {
      font-size: 12px;
      margin-top: 5px;
      color: $--color-text-primary;
      padding: 5px 10px;
      border-radius: 5px;
      width: calc(100% - 20px);
      text-align: center;
    }
    .kube_control_plane {
      background-color: $--color-primary-light-9;
    }
    .kube_node {
      background-color: $--color-success-lighter;
    }
    .etcd {
      background-color: $--color-warning-lighter;
    }
  }
  .node:hover {
    border-color: $--color-primary;
    background-color: $--color-primary-light-9;
    .role {
      color: white;
    }
    .kube_control_plane {
      background-color: $--color-primary;
    }
    .kube_node {
      background-color: $--color-success;
    }
    .etcd {
      background-color: $--color-warning;
    }
  }
  .node.active {
    border-color: $--color-warning;
    background-color: $--color-warning-lighter;
    .role {
      color: white;
    }
    .kube_control_plane {
      background-color: $--color-primary;
    }
    .kube_node {
      background-color: $--color-success;
    }
    .etcd {
      background-color: $--color-warning;
    }
  }
  .node.active:hover {
    filter: alpha(opacity=80);
    -moz-opacity: 0.8;
    -khtml-opacity: 0.8;
    opacity: 0.8;
  }

</style>
