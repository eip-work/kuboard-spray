<i18n>
en:
  sshcommon: SSH Params (apply to node {nodeName})
  etcd: "ETCD params (scope: node {nodeName})"
  etcd_member_name_required: Please input etcd_member_name
  etcd_member_name_conflict: "etcd_member_name conflict with that in node {nodeName}"
  roles: Node Role
  roleDescription: 'Node Role (scope: node {nodeName})'
  requiresAtLeastOneRole: Requires at least one role
  validate: Validate Connection
  facts: Host information
  baseInfo: Basic information
  network: Network information
  disk: Disk information
  ansible_machine: CPU Arch
  os: Operation System
  cpumem: CPU/Mem
  no_cached_facts: no cached facts found, please click the Validate Connection button above.
zh:
  sshcommon: SSH 连接参数（适用范围：节点 {nodeName}）
  etcd: "ETCD 参数（适用范围：节点 {nodeName}）"
  etcd_member_name_required: 请填写 ETCD 成员名称
  etcd_member_name_conflict: "ETCD成员名称与节点 {nodeName} 的ETCD成员名称冲突"
  roles: 节点角色
  roleDescription: 节点角色（适用范围：节点 {nodeName}）
  requiresAtLeastOneRole: 至少需要一个角色
  validate: 验证连接
  facts: 主机信息
  baseInfo: 基本信息
  network: 网络信息
  disk: 磁盘信息
  ansible_machine: CPU 架构
  os: 操作系统
  cpumem: CPU/内存
  no_cached_facts: 未找到该节点的缓存信息，请点击上方的 "验证连接" 按钮
</i18n>

<template>
  <el-form ref="form" :model="inventory" label-width="120px" label-position="left" @submit.enter.prevent>
    <SshParamsNode :cluster="cluster" :holder="inventory.all.hosts[nodeName]" :prop="`all.hosts.${nodeName}`" :clusterName="cluster.name" :nodeName="nodeName" :description="$t('sshcommon', {nodeName: nodeName})" isNode>
      <el-form-item>
        <div style="text-align: right;">
          <el-button type="primary" :loading="loadingFact" icon="el-icon-refresh-left" @click="loadFacts(false)">{{$t('validate')}}</el-button>
        </div>
      </el-form-item>
      <el-skeleton v-if="loadingFact" animated></el-skeleton>
      <div v-if="fact" class="app_form_mini app_margin_bottom">
        <el-form v-if="fact.ansible_facts" label-width="160px" label-position="left">
          <div style="text-align: center; margin-bottom: 10px; margin-top: -10px; font-weight: bold;">[ {{$t('facts')}} ]</div>
          <el-collapse v-model="activeNames">
            <el-collapse-item name="1">
              <template #title>
                <span class="package_title">{{$t('baseInfo')}}</span>
              </template>
              <div class="package_info">
                <PackageContentField v-if="fact.ansible_facts.ansible_lsb && fact.ansible_facts.ansible_lsb.description" :holder="fact.ansible_facts.ansible_lsb" fieldName="description" :label="$t('os')"></PackageContentField>
                <el-form-item v-else :label="$t('os')">
                  <span class="field_value app_text_mono">{{fact.ansible_facts.ansible_distribution}} {{fact.ansible_facts.ansible_distribution_version}}</span>
                </el-form-item>
                <PackageContentField :holder="fact.ansible_facts" fieldName="ansible_machine" :label="$t('ansible_machine')"></PackageContentField>
                <el-form-item :label="$t('cpumem')">
                  <span class="field_value app_text_mono">{{fact.ansible_facts.ansible_processor_vcpus}}core / {{Math.round(fact.ansible_facts.ansible_memtotal_mb * 10 / 1024)/10}}Gi</span>
                </el-form-item>
              </div>
            </el-collapse-item>
            <el-collapse-item name="2">
              <template #title>
                <span class="package_title">{{$t('network')}}</span>
              </template>
              <div class="package_info">
                <PackageContentField :holder="fact.ansible_facts.ansible_default_ipv4" fieldName="type"></PackageContentField>
                <PackageContentField :holder="fact.ansible_facts.ansible_default_ipv4" fieldName="interface"></PackageContentField>
                <PackageContentField :holder="fact.ansible_facts.ansible_default_ipv4" fieldName="address"></PackageContentField>
                <PackageContentField :holder="fact.ansible_facts.ansible_default_ipv4" fieldName="network"></PackageContentField>
                <PackageContentField :holder="fact.ansible_facts.ansible_default_ipv4" fieldName="gateway"></PackageContentField>
                <PackageContentField :holder="fact.ansible_facts.ansible_default_ipv4" fieldName="broadcast"></PackageContentField>
                <PackageContentField :holder="fact.ansible_facts.ansible_default_ipv4" fieldName="netmask"></PackageContentField>
                <PackageContentField :holder="fact.ansible_facts.ansible_default_ipv4" fieldName="macaddress"></PackageContentField>
                <PackageContentField :holder="fact.ansible_facts.ansible_default_ipv4" fieldName="mtu"></PackageContentField>
              </div>
            </el-collapse-item>
            <el-collapse-item name="3">
              <template #title>
                <span class="package_title">{{$t('disk')}}</span>
              </template>
              <div class="package_info">
                <template v-for="(item, key) in fact.ansible_facts.ansible_devices" :key="'disk' + key">
                  <el-form-item :label="key" v-if="item.model">
                    <el-form-item label="vendor" label-width="80px">{{item.vendor}}</el-form-item>
                    <el-form-item label="model" label-width="80px">{{item.model}}</el-form-item>
                    <el-form-item label="size" label-width="80px">{{item.size}}</el-form-item>
                  </el-form-item>
                </template>
                <el-form-item></el-form-item>
              </div>
            </el-collapse-item>
          </el-collapse>
        </el-form>
        <div v-else>
          <el-alert :closable="false" type="error">
            <div class="app_text_mono">
              {{fact.msg}}
            </div>
          </el-alert>
        </div>
      </div>
    </SshParamsNode>
    <ConfigSection v-model:enabled="enabledRoles" :label="$t('roles')" :description="$t('roleDescription', {nodeName: nodeName})" disabled>
      <el-form-item :label="$t('roles')">
        <div class="roles">
          <NodeRoleTag :enabled="isKubeControlPlane" role="kube_control_plane" @clickTag="isKubeControlPlane = !isKubeControlPlane"></NodeRoleTag>
          <NodeRoleTag :enabled="isKubeNode" role="kube_node" @clickTag="isKubeNode = !isKubeNode"></NodeRoleTag>
          <NodeRoleTag :enabled="isEtcd" role="etcd" @clickTag="isEtcd = !isEtcd"></NodeRoleTag>
        </div>
      </el-form-item>
    </ConfigSection>
    <ConfigSection v-if="enabledEtcd" v-model:enabled="enabledEtcd" label="ETCD" :description="$t('etcd', {nodeName: nodeName})" disabled>
      <FieldString :holder="inventory.all.children.target.children.etcd.hosts[nodeName]" fieldName="etcd_member_name" :rules="etcd_member_name_rules"
        :prop="`all.children.target.children.etcd.hosts.${nodeName}`" required></FieldString>
    </ConfigSection>
  </el-form>
</template>

<script>
import SshParamsNode from './common/SshParamsNode.vue'
import NodeRoleTag from './common/NodeRoleTag.vue'
import PackageContentField from './common/PackageContentField.vue'

export default {
  props: {
    cluster: { type: Object, required: true },
    nodeName: { type: String, required: true },
  },
  data() {
    return {
      enabledRoles: true,
      fact: undefined,
      loadingFact: false,
      activeNames: ['1'],
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
  components: { SshParamsNode, NodeRoleTag, PackageContentField },
  mounted () {
    this.loadFacts(true)
  },
  watch: {
    nodeName: function() {
      this.fact = undefined
      setTimeout(() => {
        this.loadFacts(true)
      }, 200)
    }
  },
  methods: {
    async loadFacts(fromCache) {
      this.fact = undefined
      this.$refs.form.validate(async flag => {
        if (this.inventory.all.children.target.children.etcd.hosts[this.nodeName] && !this.inventory.all.children.target.children.etcd.hosts[this.nodeName].etcd_member_name && !fromCache) {
          this.$message.warning(this.$t('etcd_member_name_required'))
        }
        if (flag) {
          this.loadingFact = true
          let req = {
            from_cache: fromCache,
            ansible_host: this.inventory.all.hosts[this.nodeName].ansible_host || this.inventory.all.children.target.children.k8s_cluster.vars.ansible_host,
            ansible_port: this.inventory.all.hosts[this.nodeName].ansible_port || this.inventory.all.children.target.children.k8s_cluster.vars.ansible_port,
            ansible_user: this.inventory.all.hosts[this.nodeName].ansible_user || this.inventory.all.children.target.children.k8s_cluster.vars.ansible_user,
            ansible_password: this.inventory.all.hosts[this.nodeName].ansible_password || this.inventory.all.children.target.children.k8s_cluster.vars.ansible_password,
            ansible_ssh_private_key_file: this.inventory.all.hosts[this.nodeName].ansible_ssh_private_key_file || this.inventory.all.children.target.children.k8s_cluster.vars.ansible_ssh_private_key_file,
            ansible_become: this.inventory.all.hosts[this.nodeName].ansible_become || this.inventory.all.children.target.children.k8s_cluster.vars.ansible_become,
            ansible_become_user: this.inventory.all.hosts[this.nodeName].ansible_become_user || this.inventory.all.children.target.children.k8s_cluster.vars.ansible_become_user,
            ansible_become_password: this.inventory.all.hosts[this.nodeName].ansible_become_password || this.inventory.all.children.target.children.k8s_cluster.vars.ansible_host,
          }
          await this.kuboardSprayApi.post(`/facts/cluster/${this.cluster.name}/${this.nodeName}`, req).then(resp => {
            if (fromCache) {
              if (resp.data.ansible_facts !== undefined) {
                this.fact = resp.data
              }
            } else {
              this.fact = resp.data
            }
          }).catch(e => {
            if (e.response.status !== 417) {
              let msg = '' + e
              if (e.response && e.response.data && e.response.data.message){
                msg = e.response.data.message
              }
              this.fact = {
                changed: false,
                msg: msg,
              }
            } else {
              this.fact = {
                changed: false,
                msg: this.$t('no_cached_facts')
              }
            }
          })
          this.loadingFact = false
        }
      })
    }
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
.package_title {
  font-weight: bolder;
}
.package_info {
  margin-left: 20px;
}
.field_value {
  font-size: 13px;
}
</style>
