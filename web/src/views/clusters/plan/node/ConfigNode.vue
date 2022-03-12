<i18n>
en:
  terminal: Open ssh terminal
  sshcommon: SSH Params (apply to node {nodeName})
  etcd: "ETCD params (scope: node {nodeName})"
  etcd_member_name_required: Please input etcd_member_name
  etcd_member_name_conflict: "etcd_member_name conflict with that in node {nodeName}"
  invalidName: Hostname must consist of lower case alphanumeric characters or digit, and must start with an alphanumeric character
  roles: Node Role
  roleDescription: 'Node Role (scope: node {nodeName})'
  requiresAtLeastOneRole: Requires at least one role
  pendingDelete: Pending Delete.
  pendingDeleteAction: Click Execute button to do the actural deletion.
  pendingAdd: Pending Add.
  pendingAddAction: Clieck Execute button to do the actural addition.
zh:
  terminal: 打开 SSH 终端
  sshcommon: SSH 连接参数（适用范围：节点 {nodeName}）
  etcd: "ETCD 参数（适用范围：节点 {nodeName}）"
  etcd_member_name_required: 请填写 ETCD 成员名称
  etcd_member_name_conflict: "ETCD成员名称与节点 {nodeName} 的ETCD成员名称冲突"
  invalidName: 必须由小写字母、数字组成，且必须以字母开头，以字母/数字结尾
  roles: 节点角色
  roleDescription: 节点角色（适用范围：节点 {nodeName}）
  requiresAtLeastOneRole: 至少需要一个角色
  pendingDelete: 等待删除
  pendingDeleteAction: 点击控制栏的 “节点待删除” 按钮，执行删除操作
  pendingAdd: 等待添加
  pendingAddAction: 点击控制栏的 “安装/设置集群” 或 “添加节点” 按钮，执行添加操作
</i18n>

<template>
  <el-form ref="form" :model="inventory" label-width="120px" label-position="left" @submit.enter.prevent>
    <el-alert v-if="editMode === 'view'" :type="pingpongType" :closable="false" class="app_margin_bottom app_alert_block">
      <template #title>
        <span class="app_text_mono">PING : {{nodeName}}</span>
      </template>
      <div v-if="pingpong[nodeName]" class="app_text_mono">
        <div v-if="pingpong[nodeName].ping === 'pong'">
          {{ pingpong[nodeName].ping }}
          <el-button style="float: right; margin-right: -8px;" @click="openUrlInBlank(`#/ssh/cluster/${cluster.name}/${nodeName}`)" icon="el-icon-monitor" type="primary">{{ $t('terminal')}}</el-button>
        </div>
        <span v-else>
          <div class="app_margin_bottom">
            <pre>{{ pingpong[nodeName].message }}</pre>
            <pre v-if="pingpong[nodeName].stdout"><el-tag style="margin-right: 5px;">stdout</el-tag>{{pingpong[nodeName].stdout}}</pre>
            <pre v-if="pingpong[nodeName].stderr"><el-tag style="margin-right: 5px;" type="danger" effect="dark">stderr</el-tag>{{pingpong[nodeName].stderr}}</pre>
          </div>
          <el-button type="primary" @click="$emit('ping')" :loading="pingpongLoading" icon="el-icon-lightning">PING</el-button>
          <el-button v-if="pingpong[nodeName].unreachable === false" style="float: right; margin-right: -8px;" @click="openUrlInBlank(`#/ssh/cluster/${cluster.name}/${nodeName}`)" icon="el-icon-monitor" type="primary">{{ $t('terminal')}}</el-button>
        </span>
      </div>
      <div v-else>
        <el-button type="primary" @click="$emit('ping')" :loading="pingpongLoading" icon="el-icon-lightning">PING</el-button>  
      </div>
    </el-alert>
    <el-alert v-if="pendingDelete" :title="$t('pendingDelete')" type="error" :closable="false" effect="dark" class="app_margin_bottom" show-icon>
      {{$t('pendingDeleteAction')}}
      <KuboardSprayLink href="https://kuboard-spray.cn/guide/maintain/add-replace-node.html" :size="12" style="margin-left: 20px; color: white;"></KuboardSprayLink>
    </el-alert>
    <el-alert v-if="pendingAdd" :title="$t('pendingAdd')" type="warning" :closable="false" effect="dark" class="app_margin_bottom" show-icon>
      {{$t('pendingAddAction')}}
      <KuboardSprayLink href="https://kuboard-spray.cn/guide/maintain/add-replace-node.html" :size="12" style="margin-left: 20px; color: white;"></KuboardSprayLink>
    </el-alert>
    <StateNode v-if="onlineNodes[nodeName]" :cluster="cluster" :nodeName="nodeName"></StateNode>
    <SshParamsNode :cluster="cluster" v-if="inventory.all.hosts[nodeName]"
      :holder="inventory.all.hosts[nodeName]" :prop="`all.hosts.${nodeName}`" :clusterName="cluster.name" :nodeName="nodeName" :description="$t('sshcommon', {nodeName: nodeName})">
      <NodeFact ref="nodeFact" class="app_margin_bottom" v-if="inventory.all.hosts[nodeName]"
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
        :ansible_become_password="inventory.all.hosts[nodeName].ansible_become_password || inventory.all.children.target.vars.ansible_become_password"
        :ansible_python_interpreter="inventory.all.hosts[nodeName].ansible_python_interpreter || inventory.all.children.target.vars.ansible_python_interpreter"
        :ip="inventory.all.hosts[nodeName].ip"
        :ansible_ssh_common_args="inventory.all.children.target.vars.ansible_ssh_common_args"
      ></NodeFact>
    </SshParamsNode>
    <ConfigSection v-model:enabled="enabledRoles" :label="$t('roles')" :description="$t('roleDescription', {nodeName: nodeName})" disabled anti-freeze>
      <FieldCommon :label="$t('roles')" :anti-freeze="inventory.all.hosts[nodeName].kuboardspray_node_action === 'add_node'">
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
        :anti-freeze="onlineNodes[nodeName] === undefined || inventory.all.hosts[nodeName].kuboardspray_node_action === 'add_node'"
        :prop="`all.children.target.children.etcd.hosts.${nodeName}`" required></FieldString>
    </ConfigSection>
  </el-form>
</template>

<script>
import SshParamsNode from '../common/SshParamsNode.vue'
import NodeRoleTag from '../common/NodeRoleTag.vue'
import NodeFact from '../../../common/fact/NodeFact.vue'
import StateNode from './StateNode.vue'

export default {
  props: {
    cluster: { type: Object, required: true },
    nodeName: { type: String, required: true },
    pingpong: { type: Object, required: false, default: () => {return {}} },
    pingpongLoading: { type: Boolean, required: false, default: false },
  },
  data() {
    return {
      enabledRoles: true,
      
      etcd_member_name_rules: [
        {
          validator: (rule, value, callback) => {
            if (!/^[a-z]([_a-z0-9]*[a-z0-9])?$/.test(value)) {
              return callback(this.$t('invalidName'))
            }
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
  inject: ['editMode', 'onlineNodes'],
  computed: {
    pingpongType () {
      if (this.pingpong[this.nodeName]) {
        return this.pingpong[this.nodeName].ping === 'pong' ? 'success' : 'error'
      }
      return 'info'
    },
    inventory: {
      get () {
        return this.cluster.inventory
      },
      set (v) {
        console.log(v)
      }
    },
    pendingDelete () {
      if (this.inventory.all.hosts[this.nodeName] && this.inventory.all.hosts[this.nodeName].kuboardspray_node_action === 'remove_node') {
        return true
      }
      return false
    },
    pendingAdd () {
      if (this.inventory.all.hosts[this.nodeName] && this.inventory.all.hosts[this.nodeName].kuboardspray_node_action === 'add_node') {
        return true
      }
      return false
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
        // console.log('setKubeNode', v)
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
  components: { SshParamsNode, NodeRoleTag, NodeFact, StateNode },
  mounted () {
    if (this.inventory.all.hosts[this.nodeName].ansible_host && this.$refs.nodeFact) {
      this.$refs.nodeFact.loadFacts()
    }
  },
  watch: {
    nodeName: function(newValue) {
      if (this.$refs.nodeFact) {
        this.$refs.nodeFact.clear()
        if (this.inventory.all.hosts[newValue].ansible_host) {
          setTimeout(() => {
            if (this.$refs.nodeFact) {
              this.$refs.nodeFact.loadFacts()
            }
          }, 200)
        }
      }
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
