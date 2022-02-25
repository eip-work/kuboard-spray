<i18n>
en:
  addNode: Add Node
  nodeName: Node Name
  nodeRoles: Node Roles
  conflict: Conflict with a existing node {name}.
  removeNodeFirst: Please remove node or cancel removing node first.
  invalidName: Hostname must consist of lower case alphanumeric characters or '-', and must start with an alphanumeric character
  cannotUseKeyword: Cannot use {keyword} as node name.

  nodeRequirement1: The machine on which KuboardSpray runs cannot be used as a k8s node.
  nodeRequirement2: Please make sure to use a clean machine as a node (haven't tried to install k8s on it, and no other programs run on it).
zh:
  addNode: 添加节点
  nodeName: 节点名称
  nodeRoles: 节点角色
  conflict: 与已有节点重名 {name}
  removeNodeFirst: 请先删除或者取消删除节点
  invalidName: 必须由小写字母、数字、减号组成，且必须以字母开头，以字母/数字结尾
  cannotUseKeyword: 不能使用关键字 {keyword} 作为节点名称

  nodeRequirement1: KuboardSpray 所在机器不能用作 k8s 节点
  nodeRequirement2: 请确保您的节点服务器是一个干净的机器（没有安装过 k8s，不用来运行其他程序）
</i18n>

<template>
  <el-popover placement="right-start" :title="$t('addNode')" v-if="editMode !== 'view'"
    v-model:visible="visible" :width="480" trigger="manual">
    <template #reference>
      <el-button icon="el-icon-plus" type="primary" @click="visible = true"
        :disabled="editMode === 'view'">{{$t('addNode')}}</el-button>
    </template>
    <el-form label-position="left" label-width="80px" ref="addNodeForm" :model="addNodeForm" @submit.enter.prevent>
      <el-form-item :label="$t('nodeName')" prop="name" :rules="nodeNameRules">
        <el-input v-model.trim="addNodeForm.name"></el-input>
      </el-form-item>
      <el-form-item :label="$t('nodeRoles')" prop="roles" :rules="nodeRoleRules">
        <el-checkbox-group v-model="addNodeForm.roles">
          <el-checkbox label="kube_control_plane">{{$t('node.kube_control_plane')}}</el-checkbox>
          <el-checkbox label="etcd">{{$t('node.etcd')}}</el-checkbox>
          <el-checkbox label="kube_node">{{$t('node.kube_node')}}</el-checkbox>
        </el-checkbox-group>
      </el-form-item>
      <div class="app_margin_bottom">
        <li class="desc">{{ $t('nodeRequirement1') }}</li>
        <li class="desc">{{ $t('nodeRequirement2') }}</li>
      </div>
    </el-form>
    <div style="text-align: right;">
      <KuboardSprayLink href="https://kuboard-spray.cn/guide/maintain/add-replace-node.html" style="float: left;"></KuboardSprayLink>
      <el-button icon="el-icon-close" @click="addNodeForm.visible = false">{{$t('msg.cancel')}}</el-button>
      <el-button icon="el-icon-plus" @click="addNode" type="primary">{{$t('msg.ok')}}</el-button>
    </div>
  </el-popover>
</template>

<script>
export default {
  props: {
    inventory: { type: Object, required: true },
    currentPropertiesTab: { type: String, required: false},
  },
  data () {
    return {
      addNodeForm: {
        name: '',
        roles: [],
        visible: false,
      },
      nodeNameRules: [
        { required: true, message: 'Required', trigger: 'blur' },
        {
          validator: (rule, value, callback) => {
            let keywords = ['bastion', 'target', 'kube_node', 'kube_control_plane', 'etcd', 'calico_rr', 'k8s_cluster', 'localhost', 'all']
            if (keywords.indexOf(value) > 0) {
              return callback(this.$t('cannotUseKeyword', { keyword: value }))
            }
            if (this.inventory.all.hosts[value] !== undefined) {
              return callback(this.$t('conflict', {name: value}))
            }
            if (!/^[a-z]([-a-z0-9]*[a-z0-9])?$/.test(value)) {
              return callback(this.$t('invalidName'))
            }
            callback()
          },
          trigger: 'blur',
        }
      ],
      nodeRoleRules: [
        { required: true, message: 'Required', trigger: 'blur' },
      ]
    }
  },
  computed: {
    inventoryRef: {
      get () {return this.inventory},
      set () {},
    },
    visible: {
      get () { return this.addNodeForm.visible },
      set (v) {
        if (v) {
          if (this.isClusterInstalled && this.isClusterOnline) {
            if (this.pendingRemoveNodes.length > 0) {
              this.$message.error(this.$t('removeNodeFirst'))
              return
            }
          }
        }
        this.addNodeForm.visible = v
      },
    }
  },
  inject: ['editMode', 'isClusterInstalled', 'isClusterOnline', 'pendingRemoveNodes'],
  methods: {
    addNode () {
      this.$refs.addNodeForm.validate(flag => {
        if (flag) {
          this.inventoryRef.all.hosts[this.addNodeForm.name] = {
            'kuboardspray_node_action': 'add_node',
          }
          for (let role of this.addNodeForm.roles) {
            if (role === 'etcd') {
              this.inventoryRef.all.children.target.children.etcd.hosts[this.addNodeForm.name] = {
                'kuboardspray_node_action': 'add_node',
              }
            } else {
              this.inventoryRef.all.children.target.children.k8s_cluster.children[role].hosts[this.addNodeForm.name] = {
                'kuboardspray_node_action': 'add_node',
              }
            }
          }
          this.$emit('update:currentPropertiesTab', 'NODE_' + this.addNodeForm.name)
          this.addNodeForm.visible = false
        }
      })
    }
  }
}
</script>

<style>
.desc {
  color: var(--el-text-color-secondary);
  font-size: 12px;
  line-height: 24px;
}
</style>