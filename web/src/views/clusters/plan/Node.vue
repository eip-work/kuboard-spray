<i18n>
en:
  confirmDelete: Are you sure to delete this node ?
zh:
  confirmDelete: 您是否确认要删除此节点？
</i18n>

<template>
  <div class="node_wrapper">
    <div class="status" v-if="nodes[name]">
      <el-button type="success" circle icon="el-icon-cloudy"></el-button>
    </div>
    <div class="delete_button" v-if="!hideDeleteButton && editMode !== 'view'">
        <el-popconfirm icon="el-icon-info" icon-color="red" :title="$t('confirmDelete')" @confirm="$emit('delete_button')" placement="right-start">
          <template #reference>
            <el-button icon="el-icon-delete" type="danger" circle @submit.prevent.stop></el-button>
          </template>
        </el-popconfirm>
    </div>
    <div :class="nodeClass">
      <div class="app_text_mono" style="font-weight: bold;">
        <span v-if="name !== 'localhost' && name !== 'bastion'">{{ name }}</span>
        <span v-else>{{ $t('obj.' + name) }}</span>
      </div>
      <div class="app_text_mono" v-if="inventory.all.hosts[name]">
        {{ inventory.all.hosts[name].ansible_host }}
      </div>
      <div :class="role + ' role app_text_mono'" v-for="(_, role) in roles" :key="'r' + role">{{ roleName(role) }}</div>
      <slot></slot>
    </div>
  </div>
</template>

<script>

export default {
  props: {
    name: { type: String, required: true },
    inventory: { type: Object, required: true },
    active: { type: Boolean, required: false, default: false },
    hideDeleteButton: { type: Boolean, required: false, default: false },
    nodes: { type: Object, required: false, default: () => {return {}} },
  },
  inject: ['editMode'],
  computed: {
    nodeClass () {
      let result = 'node'
      if (this.active) {
        result += ' active'
      }
      if (this.nodes[this.name]) {
        result += ' online'
      }
      return result
    },
    roles () {
      let result = {}
      for (let role in this.inventory.all.children.target.children.k8s_cluster.children) {
        for (let n in this.inventory.all.children.target.children.k8s_cluster.children[role].hosts) {
          if (n === this.name) {
            result[role] = true
          }
        }
      }
      for (let n in this.inventory.all.children.target.children.etcd.hosts) {
        if (n === this.name) {
          result['etcd'] = true
        }
      }
      return result
    }
  },
  mounted () {
  },
  methods: {
    roleName(role) {
      if (role === 'etcd') {
        return this.inventory.all.children.target.children.etcd.hosts[this.name].etcd_member_name || this.$t('node.' + role)
      } else {
        return this.$t('node.' + role)
      }
    }
  }
}
</script>

<style scoped lang="scss">
.node_wrapper{
  .status {
    height: 0px;
    text-align: left;
  }
  .delete_button {
    height: 0px;
    text-align: right;
    display: none;
  }
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
  .node.online {
    border-color: $--color-success;
  }
  .node.online.active {
    background-color: $--color-success-lighter;
  }
}

.node_wrapper:hover .delete_button {
  display: block;
}

</style>
