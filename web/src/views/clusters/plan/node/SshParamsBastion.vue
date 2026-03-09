<i18n>
en:
  bastionUsage: PangeeCluster can access Kubernetes Cluster Nodes through bastion.
  addSshKey: Add Private Key
  ansible_host_placeholder: 'PangeeCluster use this ip or hostname to connect to the node.'
  default_value: 'Default: {default_value} (inhirit from value configured in Global Config tab)'
  duplicateIP: "IP address conflict with {node}"
  description: Access target nodes throguh bastion/jumpserver

  AllowTcpForwarding: 跳板机 /etc/ssh/sshd_config 文件中的 AllowTcpForwarding 属性必须设置为 true
zh:
  bastionUsage: PangeeCluster 可以通过跳板机或堡垒机访问将要安装 K8S 集群的目标节点。
  addSshKey: 添加私钥
  ansible_host_placeholder: 'PangeeCluster 连接该主机时所使用的主机名或 IP 地址'
  default_value: '默认值：{default_value} （继承自全局设置标签页中的配置）'
  duplicateIP: "IP 地址不能与其他节点相同：{node}"
  description: 通过跳板机或者堡垒机访问目标节点
  protocol: 跳板机协议

  ansible_host: '主机'
  ansible_port: 'SSH 端口'
  ansible_port_placeholder: 'PangeeCluster 连接该主机时所使用的 SHH 端口'
  
  ansible_user: '用户名'
  ansible_user_placeholder: 'PangeeCluster 连接该主机时所使用的用户名'
  ansible_password: '密码'
  ansible_password_placeholder: 'PangeeCluster 连接该主机时所使用的密码'
  ansible_ssh_private_key_file: '私钥文件'
  ansible_ssh_private_key_file_placeholder: 'PangeeCluster 连接该主机时所使用的私钥文件'
  ansible_become: '切换身份'
  ansible_become_placeholder: 'PangeeCluster 登录到该主机后，是否使用 su 命令切换身份'
  ansible_become_user: '切换到用户'
  ansible_become_user_placeholder: 'PangeeCluster 登录该主机后，使用 su 命令切换到用户名'
  ansible_become_password: '切换密码'
  ansible_become_password_placeholder: '切换密码'

  ansible_python_interpreter: 'Python 路径'

  AllowTcpForwarding: 跳板机 /etc/ssh/sshd_config 文件中的 AllowTcpForwarding 属性必须设置为 true
</i18n>


<template>
  <ConfigSection ref="configSection" v-model:enabled="bastionEnabled" :label="$t('obj.bastion')"
    :description="t('description')" anti-freeze>
    <template #more>
      {{ t('bastionUsage') }}
    </template>
    <el-form-item :label="t('protocol')" required>
      <el-radio-group v-model="computedBastionType" :disabled="editMode != 'edit'">
        <el-radio value="ssh">SSH</el-radio>
        <el-radio value="socks5">SOCKS5</el-radio>
      </el-radio-group>
    </el-form-item>
    <EditString v-model="holder.ansible_host" :label="t('ansible_host')" :prop="`all.hosts.${nodeName}.ansible_host`"
      anti-freeze :placeholder="t('ansible_host_placeholder')" :rules="hostRules"></EditString>
    <div
      style="margin-left: 120px;margin-bottom: 10px; margin-top: -10px; color: var(--el-color-white); background-color: var(--el-color-warning); padding: 5px 15px; font-weight: bolder"
      class="app_text_mono">
      {{ t("AllowTcpForwarding") }}
    </div>
    <EditString v-model="holder.ansible_port" :label="t('ansible_port')" :prop="`all.hosts.${nodeName}.ansible_port`"
      :placeholder="placeholder('ansible_port')" anti-freeze required></EditString>
    <template v-if="computedBastionType == 'ssh'">
      <EditString v-model="holder.ansible_user" :label="t('ansible_user')" :prop="`all.hosts.${nodeName}.ansible_user`"
        :placeholder="placeholder('ansible_user')" anti-freeze required></EditString>
      <EditSelect v-model="holder.ansible_ssh_private_key_file" :label="t('ansible_ssh_private_key_file')"
        :loadOptions="loadSshKeyList" anti-freeze clearable :placeholder="placeholder('ansible_ssh_private_key_file')">
        <template #edit>
          <el-button type="primary" plain style="margin-left: 10px;" icon="el-icon-plus"
            @click="$refs.addPrivateKey.show()">{{ t('addSshKey') }}</el-button>
        </template>
      </EditSelect>
      <EditString v-model="holder.ansible_password" :label="t('ansible_password')" show-password anti-freeze clearable
        :placeholder="placeholder('ansible_password')"></EditString>
      <SshAddPrivateKey ref="addPrivateKey" ownerType="cluster" :ownerName="cluster.name"></SshAddPrivateKey>
    </template>
    <template v-else-if="computedBastionType == 'socks5'">
      <EditString v-model="holder.ansible_user" :label="t('ansible_user')" :prop="`all.hosts.${nodeName}.ansible_user`"
        :placeholder="placeholder('ansible_user')" anti-freeze></EditString>
      <EditString v-model="holder.ansible_password" :label="t('ansible_password')" show-password anti-freeze clearable
        :placeholder="placeholder('ansible_password')"></EditString>
    </template>
  </ConfigSection>
</template>

<script>
import SshAddPrivateKey from '../../../private_key/SshAddPrivateKey.vue'

export default {
  props: {
    cluster: { type: Object, required: true },
    nodeName: { type: String, required: false },
    holder: { type: Object, required: true },
    prop: { type: String, required: true },
  },
  data() {
    return {
      hostRules: [
        { required: true, message: this.$t('field.ansible_host') + this.$t('field.is_required_field'), trigger: 'blur' },
      ]
    }
  },
  inject: ['editMode'],
  computed: {
    computedBastionType: {
      get() {
        if (this.inventory.all.hosts.bastion && this.inventory.all.hosts.bastion.bastion_type) {
          return this.inventory.all.hosts.bastion.bastion_type
        }
        return "ssh"
      },
      set(type) {
        this.inventory.all.hosts.bastion.bastion_type = type;
      }
    },
    inventory: {
      get() { return this.cluster.inventory },
      set() { }
    },
    bastionEnabled: {
      get() {
        return this.inventory.all.hosts.bastion !== undefined
      },
      set(v) {
        if (v) {
          this.inventory.all.hosts.bastion = this.inventory.all.hosts.bastion || { ansible_host: '', ansible_user: '' }
        } else {
          delete this.inventory.all.hosts.bastion
          delete this.inventory.all.vars.ansible_ssh_common_args
        }
      }
    },
    holderRef: {
      get() { return this.holder },
      set() { }
    },
    ansible_become: {
      get() {
        if (this.holder.ansible_become !== undefined) {
          return this.holder.ansible_become
        }
        return this.cluster.inventory.all.vars.ansible_become
      },
      set(v) {
        this.holderRef.ansible_become = v
      }
    }
  },
  components: { SshAddPrivateKey },
  mounted() {
  },
  watch: {
    holder: {
      handler: function (bastion) {
        if (bastion && bastion.ansible_host) {
          if (this.computedBastionType == 'ssh') {
            let sshPass = ''
            if (bastion['ansible_password']) {
              sshPass = `sshpass -p '${bastion['ansible_password']}' `
            }
            let temp = `-o ProxyCommand="${sshPass}ssh -F /dev/null -o ConnectTimeout=10 -o ConnectionAttempts=100 -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -W %h:%p -p`
            temp += bastion["ansible_port"] + " " + bastion["ansible_user"] + "@" + bastion["ansible_host"]
            if (bastion["ansible_ssh_private_key_file"]) {
              temp += " -i " + bastion["ansible_ssh_private_key_file"]
            }
            temp += '"'
            this.inventory.all.vars.ansible_ssh_common_args = temp
          } else if (this.computedBastionType == 'socks5') {
            let temp = `-o ProxyCommand="nc -X 5 -x ${bastion.ansible_host}:${bastion.ansible_port} %h %p"`
            if (bastion['ansible_user']) {
              temp = `-o ProxyCommand="nc -X 5 -x ${bastion['ansible_user']}@${bastion.ansible_host}:${bastion.ansible_port} %h %p"`
            }
            if (bastion['ansible_password'] && bastion['ansible_user']) {
              temp = `-o ProxyCommand="nc -X 5 -x ${bastion['ansible_user']}:${bastion['ansible_password']}@${bastion.ansible_host}:${bastion.ansible_port} %h %p"`
            }
            this.inventory.all.vars.ansible_ssh_common_args = temp
          } else {
            delete this.inventory.all.vars.ansible_ssh_common_args
          }
        } else {
          delete this.inventory.all.vars.ansible_ssh_common_args
        }
      },
      deep: true,
    }
  },
  methods: {
    placeholder(fieldName) {
      return this.$t('field.' + fieldName + '_placeholder')
    },
    async loadSshKeyList() {
      let result = []
      await this.pangeeClusterApi.get(`/private-keys/cluster/${this.cluster.name}`).then(resp => {
        for (let item of resp.data.data) {
          result.push({ label: item, value: '{{ pangeecluster_cluster_dir }}/private-key/' + item })
        }
      })
      return result
    },
  }
}
</script>

<style scoped lang="css"></style>
