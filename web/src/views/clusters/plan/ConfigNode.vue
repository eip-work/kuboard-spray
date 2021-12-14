<i18n>
en:
  sshcommon: SSH Params (apply to node {nodeName})
zh:
  sshcommon: SSH 连接参数（适用范围：节点 {nodeName}）
  etcd: "ETCD 参数（适用范围：节点 {nodeName}）"
  roles: 节点角色
  roleDescription: 节点角色（适用范围：节点 {nodeName}）
  requiresAtLeastOneRole: 至少需要一个角色
</i18n>

<template>
  <div>
    <SshParams :holder="inventory.all.hosts[nodeName]" :prop="`all.hosts.${nodeName}`" :clusterName="clusterName" :description="$t('sshcommon', {nodeName: nodeName})" isNode></SshParams>
    <ConfigSection v-model:enabled="enabledRoles" :label="$t('roles')" :description="$t('roleDescription', {nodeName: nodeName})" disabled>
      <el-form-item :label="$t('roles')">
        <div class="roles">
          <NodeRoleTag :enabled="isKubeControlPlane" role="kube_control_plane" @click="isKubeControlPlane = !isKubeControlPlane"></NodeRoleTag>
          <NodeRoleTag :enabled="isKubeNode" role="kube_node" @click="isKubeNode = !isKubeNode"></NodeRoleTag>
          <NodeRoleTag :enabled="isEtcd" role="etcd" @click="isEtcd = !isEtcd"></NodeRoleTag>
        </div>
      </el-form-item>
    </ConfigSection>
    <ConfigSection v-if="enabledEtcd" v-model:enabled="enabledEtcd" label="ETCD" :description="$t('etcd', {nodeName: nodeName})">
      <FieldString :holder="inventory.all.children.etcd.hosts[nodeName]" fieldName="etcd_member_name"
        :prop="`all.children.etcd.hosts.${nodeName}`" required></FieldString>
    </ConfigSection>
  </div>
</template>

<script>
import SshParams from './common/SshParams.vue'
import ConfigSection from './ConfigSection.vue'
import NodeRoleTag from './common/NodeRoleTag.vue'

export default {
  props: {
    inventory: { type: Object, required: true },
    clusterName: { type: String, required: true },
    nodeName: { type: String, required: true },
    resourcePackage: { type: Object, required: false, default: undefined },
  },
  data() {
    return {
      enabledRoles: true
    }
  },
  computed: {
    inventoryRef: {
      get () {
        return this.inventory
      },
      set (v) {
        console.log(v)
      }
    },
    enabledEtcd: {
      get () {
        for (let key in this.inventory.all.children.etcd.hosts) {
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
        for (let key in this.inventory.all.children.k8s_cluster.children.kube_control_plane.hosts) {
          if (key === this.nodeName) {
            return true
          }
        }
        return false
      },
      set (v) {
        console.log('setKubeControlPlane', v)
        if (v) {
          this.inventoryRef.all.children.k8s_cluster.children.kube_control_plane.hosts[this.nodeName] = {}
        } else {
          if (!this.isEtcd && !this.isKubeNode) {
            this.$message.error(this.$t('requiresAtLeastOneRole'))
          } else {
            delete this.inventoryRef.all.children.k8s_cluster.children.kube_control_plane.hosts[this.nodeName]
          }
        }
      }
    },
    isKubeNode: {
      get () {
        for (let key in this.inventory.all.children.k8s_cluster.children.kube_node.hosts) {
          if (key === this.nodeName) {
            return true
          }
        }
        return false
      },
      set (v) {
        console.log('setKubeNode', v)
        if (v) {
          this.inventoryRef.all.children.k8s_cluster.children.kube_node.hosts[this.nodeName] = {}
        } else {
          if (!this.isEtcd && !this.isKubeControlPlane) {
            this.$message.error(this.$t('requiresAtLeastOneRole'))
          } else {
            delete this.inventoryRef.all.children.k8s_cluster.children.kube_node.hosts[this.nodeName]
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
          this.inventoryRef.all.children.etcd.hosts[this.nodeName] = {}
        } else {
          if (!this.isKubeControlPlane && !this.isKubeNode) {
            this.$message.error(this.$t('requiresAtLeastOneRole'))
          } else {
            delete this.inventoryRef.all.children.etcd.hosts[this.nodeName]
          }
        }
      }
    },
  },
  components: { SshParams, ConfigSection, NodeRoleTag },
  mounted () {
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
