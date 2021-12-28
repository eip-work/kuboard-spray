<i18n>
en:
  sshcommon: SSH Params (apply to node {nodeName})
  etcd: "ETCD params (scope: node {nodeName})"
  etcd_member_name_required: Please input etcd_member_name
  etcd_member_name_conflict: "etcd_member_name conflict with that in node {nodeName}"
  roles: Node Role
  roleDescription: 'Node Role (scope: node {nodeName})'
  requiresAtLeastOneRole: Requires at least one role
zh:
  sshcommon: SSH 连接参数（适用范围：节点 {nodeName}）
  etcd: "ETCD 参数（适用范围：节点 {nodeName}）"
  etcd_member_name_required: 请填写 ETCD 成员名称
  etcd_member_name_conflict: "ETCD成员名称与节点 {nodeName} 的ETCD成员名称冲突"
  roles: 节点角色
  roleDescription: 节点角色（适用范围：节点 {nodeName}）
  requiresAtLeastOneRole: 至少需要一个角色
</i18n>

<template>
  <el-form ref="form" :model="inventory" label-width="120px" label-position="left" @submit.enter.prevent>
    <SshParamsNode :cluster="cluster" :holder="inventory.all.hosts[nodeName]" :prop="`all.hosts.${nodeName}`" :clusterName="cluster.name" :nodeName="nodeName" :description="$t('sshcommon', {nodeName: nodeName})" isNode>
      <NodeFact ref="nodeFact" class="app_margin_bottom"
        :form="$refs.form"
        node_owner_type="cluster"
        :node_owner="cluster.name"
        :node_name="nodeName"
        :ansible_host="inventory.all.hosts[nodeName].ansible_host || inventory.all.children.target.vars.ansible_host"
        :ansible_port="inventory.all.hosts[nodeName].ansible_port || inventory.all.children.target.vars.ansible_port"
        :ansible_user="inventory.all.hosts[nodeName].ansible_user || inventory.all.children.target.vars.ansible_user"
        :ansible_password="inventory.all.hosts[nodeName].ansible_password || inventory.all.children.target.vars.ansible_password"
        :ansible_ssh_private_key_file="inventory.all.hosts[nodeName].ansible_ssh_private_key_file || inventory.all.children.target.vars.ansible_ssh_private_key_file"
        :ansible_become="inventory.all.hosts[nodeName].ansible_become || inventory.all.children.target.vars.ansible_become"
        :ansible_become_user="inventory.all.hosts[nodeName].ansible_become_user || inventory.all.children.target.vars.ansible_become_user"
        :ansible_become_password="inventory.all.hosts[nodeName].ansible_become_password || inventory.all.children.target.vars.ansible_host"
      ></NodeFact>
    </SshParamsNode>
    <ConfigSection v-model:enabled="enabledRoles" :label="$t('roles')" :description="$t('roleDescription', {nodeName: nodeName})" disabled anti-freeze>
      <FieldCommon :label="$t('roles')">
        <template #edit>
          <div class="roles">
            <NodeRoleTag :enabled="isKubeControlPlane" role="kube_control_plane" @clickTag="isKubeControlPlane = !isKubeControlPlane"></NodeRoleTag>
            <NodeRoleTag :enabled="isKubeNode" role="kube_node" @clickTag="isKubeNode = !isKubeNode"></NodeRoleTag>
            <NodeRoleTag :enabled="isEtcd" role="etcd" @clickTag="isEtcd = !isEtcd"></NodeRoleTag>
          </div>
        </template>
        <template #view>
          <div class="roles">
            <NodeRoleTag :enabled="isKubeControlPlane" role="kube_control_plane"></NodeRoleTag>
            <NodeRoleTag :enabled="isKubeNode" role="kube_node"></NodeRoleTag>
            <NodeRoleTag :enabled="isEtcd" role="etcd"></NodeRoleTag>
          </div>
        </template>
      </FieldCommon>
    </ConfigSection>
    <ConfigSection v-if="enabledEtcd" v-model:enabled="enabledEtcd" label="ETCD" :description="$t('etcd', {nodeName: nodeName})" disabled anti-freeze>
      <FieldString :holder="inventory.all.children.target.children.etcd.hosts[nodeName]" fieldName="etcd_member_name" :rules="etcd_member_name_rules"
        :prop="`all.children.target.children.etcd.hosts.${nodeName}`" required></FieldString>
    </ConfigSection>
  </el-form>
</template>

<script>
import SshParamsNode from '../common/SshParamsNode.vue'
import NodeRoleTag from '../common/NodeRoleTag.vue'
import NodeFact from '../../../fact/NodeFact.vue'

export default {
  props: {
    cluster: { type: Object, required: true },
    nodeName: { type: String, required: true },
  },
  data() {
    return {
      enabledRoles: true,
      
      etcd_member_name_rules: [
        {
          validator: (rule, value, callback) => {
            for (let key in this.inventory.all.children.target.children.etcd.hosts) {
              if (key !== this.nodeName && this.inventory.all.children.target.children.etcd.hosts[key].etcd_member_name === value) {
                return callback(this.$t('etcd_member_name_conflict', {nodeName: key}))
              }
            }
            return callback()
          },
          trigger: 'blur',
        }
      ]
    }
  },
  computed: {
    inventory: {
      get () {
        return this.cluster.inventory
      },
      set (v) {
        console.log(v)
      }
    },
    enabledEtcd: {
      get () {
        for (let key in this.inventory.all.children.target.children.etcd.hosts) {
          if (key === this.nodeName) {
            return true
          }
        }
        return false
      },
      set (v) {
        console.log(v)
      }
    },
    isKubeControlPlane: {
      get () {
        for (let key in this.inventory.all.children.target.children.k8s_cluster.children.kube_control_plane.hosts) {
          if (key === this.nodeName) {
            return true
          }
        }
        return false
      },
      set (v) {
        console.log('setKubeControlPlane', v)
        if (v) {
          this.inventory.all.children.target.children.k8s_cluster.children.kube_control_plane.hosts[this.nodeName] = {}
        } else {
          if (!this.isEtcd && !this.isKubeNode) {
            this.$message.error(this.$t('requiresAtLeastOneRole'))
          } else {
            delete this.inventory.all.children.target.children.k8s_cluster.children.kube_control_plane.hosts[this.nodeName]
          }
        }
      }
    },
    isKubeNode: {
      get () {
        for (let key in this.inventory.all.children.target.children.k8s_cluster.children.kube_node.hosts) {
          if (key === this.nodeName) {
            return true
          }
        }
        return false
      },
      set (v) {
        console.log('setKubeNode', v)
        if (v) {
          this.inventory.all.children.target.children.k8s_cluster.children.kube_node.hosts[this.nodeName] = {}
        } else {
          if (!this.isEtcd && !this.isKubeControlPlane) {
            this.$message.error(this.$t('requiresAtLeastOneRole'))
          } else {
            delete this.inventory.all.children.target.children.k8s_cluster.children.kube_node.hosts[this.nodeName]
          }
        }
      }
    },
    isEtcd: {
      get () {
        return this.enabledEtcd
      },
      set (v) {
        console.log('setEtcd', v)
        if (v) {
          this.inventory.all.children.target.children.etcd.hosts[this.nodeName] = {}
        } else {
          if (!this.isKubeControlPlane && !this.isKubeNode) {
            this.$message.error(this.$t('requiresAtLeastOneRole'))
          } else {
            delete this.inventory.all.children.target.children.etcd.hosts[this.nodeName]
          }
        }
      }
    },
  },
  components: { SshParamsNode, NodeRoleTag, NodeFact },
  mounted () {
    this.$refs.nodeFact.loadFacts()
  },
  watch: {
    nodeName: function() {
      this.fact = undefined
      setTimeout(() => {
        this.$refs.nodeFact.loadFacts()
      }, 200)
    }
  },
  methods: {
  }
}
</script>

<style scoped lang="scss">
.roles {
  display: flex;
  flex-wrap: wrap;
  .role {
    margin-bottom: 10px;
  }
  margin-bottom: -10px;
}
</style>
