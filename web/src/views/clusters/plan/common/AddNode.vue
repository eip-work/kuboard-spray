<i18n>
en:
  addNode: Add Node
  nodeName: Node Name
  nodeRoles: Node Roles
  conflict: Conflict with a existing node {name}.
zh:
  addNode: 添加节点
  nodeName: 节点名称
  nodeRoles: 节点角色
  conflict: 与已有节点重名 {name}
</i18n>

<template>
  <el-popover placement="right-start" :title="$t('addNode')"
    v-model:visible="addNodeForm.visible" :width="420" trigger="manual">
    <template #reference>
      <el-button icon="el-icon-plus" type="primary" @click="addNodeForm.visible = true"
        :disabled="editMode === 'view'">{{$t('addNode')}}</el-button>
    </template>
    <el-form label-position="left" label-width="80px" ref="addNodeForm" :model="addNodeForm">
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
    </el-form>
    <div style="text-align: right;">
      <el-button icon="el-icon-close" @click="addNodeForm.visible = false">{{$t('msg.cancel')}}</el-button>
      <el-button icon="el-icon-plus" @click="addNode" type="primary">{{$t('msg.ok')}}</el-button>
    </div>
  </el-popover>
</template>

<script>
export default {
  props: {
    inventory: { type: Object, required: true },
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
            if (this.inventory.all.hosts[value] !== undefined) {
              return callback(this.$t('conflict', {name: value}))
            }
            if (!/^[a-zA-Z][a-zA-Z0-9_]{3,21}$/.test(value)) {
              return callback('必须以字母开头，可以包含数字和字母，长度为 [4-20]')
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
    }
  },
  inject: ['editMode'],
  methods: {
    addNode () {
      this.$refs.addNodeForm.validate(flag => {
        if (flag) {
          this.inventoryRef.all.hosts[this.addNodeForm.name] = {}
          for (let role of this.addNodeForm.roles) {
            if (role === 'etcd') {
              this.inventoryRef.all.children.etcd.hosts[this.addNodeForm.name] = {}
            } else {
              this.inventoryRef.all.children.k8s_cluster.children[role].hosts[this.addNodeForm.name] = {}
            }
          }
          this.currentPropertiesTab = 'NODE_' + this.addNodeForm.name
          this.addNodeForm.visible = false
        }
      })
    }
  }
}
</script>

<style>

</style>