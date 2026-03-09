<i18n>
en:
  addSshKey: Add Private Key
  ansible_host_placeholder: 'PangeeCluster use this ip or hostname to connect to the node.'
  default_value: 'Default: {default_value} (inhirit from value configured in Global Config tab)'
  duplicateIP: "IP address conflict with {node}"
  ip: Bind to IP
  ip_placeholder: 'Default: {default_value} (inhirit from value ansible_host)'
  longTimeLoading: Loading... cost about 20-30s
  password_and_bastion: You are using bastion/jumpserver to access the node, if you use password here, it takes longer time to create ssh connection, it would be better to clear the password and provide ssh_private_key.
  password_in_global: You defined password in global configuration, it's suggested that you clear it.
  speedup: SSH Multiplexing
  cannotUseLocalhostAsTarget: Cannot use the machine that pangeecluster runs on as a target node.
  inhirit: inhirit from value configured in Global Config tab
  rootuser: Must use root user

  ansible_host: 'Host'
  ansible_port: 'SSH Port'
  ansible_port_placeholder: 'PangeeCluster use this SSH port to connect to the node'

  is_required_field: ' is required'

  ansible_user: 'Username'
  ansible_user_placeholder: 'PangeeCluster use this username to connect to the node'
  ansible_password: 'Password'
  ansible_password_placeholder: 'PangeeCluster use this password to connect to the node'
  ansible_ssh_private_key_file: 'Private Key File'
  ansible_ssh_private_key_file_placeholder: 'PangeeCluster use this private key file to connect to the node'
  ansible_become: 'Privilege Escalation'
  ansible_become_placeholder: 'Whether to use su command to switch identity after PangeeCluster log in to the node'
  ansible_become_user: 'Become User'
  ansible_become_user_placeholder: 'The username to switch to after PangeeCluster log in to the node and use su command to switch identity'
  ansible_become_password: 'Become Password'
  ansible_become_password_placeholder: 'The password for switching identity'

  ansible_python_interpreter: 'Python Interpreter Path'
  ansible_python_interpreter_placeholder: 'The Python interpreter path that PangeeCluster uses when executing Ansible tasks on the node'

zh:
  addSshKey: 添加私钥
  ansible_host_placeholder: 'PangeeCluster 连接该主机时所使用的主机名或 IP 地址'
  default_value: '默认值：{default_value} （继承自全局设置标签页中的配置）'
  duplicateIP: "IP 地址不能与其他节点相同：{node}"
  ip: 绑定到 IP
  ip_placeholder: '默认值：{default_value} （默认与主机 IP 相同）'
  longTimeLoading: 加载中... 可能需要 20-30 秒
  ipDescription: kubernetes 使用的地址有可能与 pangee-cluster 访问该机器时所使用的地址不同，此处指定 kubernetes 所使用的地址（必须为内网地址）
  password_and_bastion: 您将使用跳板机访问节点，如果在此处使用密码访问，创建 ssh 连接的时间较长，并且有很大的概率会发生ssh连接断开的情况，建议清除密码并提供 “私钥文件”。
  password_in_global: 您在全局参数中配置了密码，建议清除。
  speedup: 加速执行
  cannotUseLocalhostAsTarget: 不能使用 PangeeCluster 所在机器作为目标节点
  inhirit: 继承自全局设置标签页中的配置
  rootuser: 必须使用 root 用户

  ansible_host: '主机'
  ansible_port: 'SSH 端口'
  ansible_port_placeholder: 'PangeeCluster 连接该主机时所使用的 SHH 端口'

  is_required_field: " 为必填字段"
  
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
  ansible_python_interpreter_placeholder: 'PangeeCluster 在该主机上执行 Ansible 任务时所使用的 Python 解释器路径'
</i18n>


<template>
  <ConfigSection ref="configSection" v-model:enabled="enableSsh" label="SSH" :description="description" disabled
    anti-freeze>
    <EditCommon v-model="holder.ansible_host" :prop="`all.hosts.${nodeName}.ansible_host`" :label="t('ansible_host')"
      :anti-freeze="onlineNodes[nodeName] === undefined || holder.pangeecluster_node_action === 'add_node'"
      :rules="hostRules">
      <template #edit>
        <el-input v-model.trim="ansible_host" :placeholder="t('ansible_host_placeholder')"></el-input>
      </template>
    </EditCommon>
    <EditString v-model="holder.ansible_port" :label="t('ansible_port')" :prop="`all.hosts.${nodeName}`"
      :placeholder="placeholder('ansible_port')" anti-freeze
      :required="!cluster.inventory.all.vars.ansible_port"></EditString>
    <EditCommon v-model="holder.ansible_user" :label="t('ansible_user')" :placeholder="placeholder('ansible_user')"
      anti-freeze>
      <template #edit>
        <el-input v-model.trim="ansible_user" :placeholder="placeholder('ansible_user')"></el-input>
        <el-tag class="app_text_mono" style="display: block; line-height: 18px;" type="warning">
          {{ t('rootuser') }} </el-tag>
      </template>
    </EditCommon>
    <EditSelect v-model="holder.ansible_ssh_private_key_file" :label="t('ansible_ssh_private_key_file')"
      :loadOptions="loadSshKeyList" anti-freeze clearable :placeholder="placeholder('ansible_ssh_private_key_file')">
      <template #edit>
        <el-button type="primary" plain style="margin-left: 10px;" icon="el-icon-plus"
          @click="$refs.addPrivateKey.show()">{{ t('addSshKey') }}</el-button>
      </template>
    </EditSelect>
    <EditString v-model="holder.ansible_password" :label="t('ansible_password')" show-password anti-freeze clearable
      :placeholder="placeholder('ansible_password')"></EditString>
    <el-alert type="warning" :closable="false"
      v-if="cluster && cluster.inventory.all.hosts.bastion && (holder.ansible_password || cluster.inventory.all.vars.ansible_password)"
      style="margin-left: 120px; width: calc(100% - 120px);">
      {{ t('password_and_bastion') }}
      <PangeeClusterLink href="https://github.com/opencmit/pangee-cluster/blob/main/docs/guide/extra/speedup.md" style="margin-left: 10px;"
        :size="12">
      </PangeeClusterLink>
      <li v-if="cluster && cluster.inventory.all.vars.ansible_password">{{ t('password_in_global') }}
      </li>
    </el-alert>
    <EditSelect v-model="holder.ansible_python_interpreter" :label="t('ansible_python_interpreter')" anti-freeze
      clearable :loadOptions="loadPythonInterpreter" :placeholder="placeholder('ansible_python_interpreter')"
      allow-create filterable>
    </EditSelect>
    <EditCommon v-model="holder.ip" :prop="`all.hosts.${nodeName}.ip`"
      :anti-freeze="onlineNodes[nodeName] === undefined" :label="t('ip')"
      :placeholder="t('ip_placeholder', { default_value: holder.ansible_host })" :rules="ipRules">
      <template #edit>
        <el-select v-model="holderRef.ip" style="width: 100%;" :loading="optionIpsLoading"
          @visible-change="loadOptionIps" :loading-text="t('longTimeLoading')">
          <el-option v-for="item in optionIps" :key="item" :value="item">
            <span class="app_text_mono">{{ item }}</span>
          </el-option>
        </el-select>
        <div style="font-size: 12px; color: #666; line-height: 20px; margin-top: 2px;">{{ t('ipDescription') }}</div>
      </template>
      <template #view>
        <span class="app_text_mono">{{ holderRef.ip || ansible_host }}</span>
      </template>
    </EditCommon>
    <slot></slot>
    <SshAddPrivateKey ref="addPrivateKey" ownerType="cluster" :ownerName="cluster.name"></SshAddPrivateKey>
  </ConfigSection>
</template>

<script>
import SshAddPrivateKey from '../../../private_key/SshAddPrivateKey.vue'
import { Netmask } from 'netmask'
import { useI18n } from 'vue-i18n';

export default {
  props: {
    cluster: { type: Object, required: true },
    nodeName: { type: String, required: false },
    holder: { type: Object, required: true },
    prop: { type: String, required: true },
    description: { type: String, required: true },
  },
  data() {
    return {
      hostRules: [
        { required: true, message: 'ansible_host' + this.i18n('is_required_field'), trigger: 'blur' },
        { validator: this.$validators.ipv4, trigger: 'change' },
        {
          validator: (rule, value, callback) => {
            if (value === location.hostname) {
              console.log(location.hostname)
              return callback(this.i18n('cannotUseLocalhostAsTarget'))
            }
            for (let key in this.cluster.inventory.all.hosts) {
              if (key !== this.nodeName && this.cluster.inventory.all.hosts[key].ansible_host === value) {
                return callback(this.i18n('duplicateIP', { node: key }))
              }
            }
            return callback()
          },
          trigger: 'blur'
        },
      ],
      ipRules: [
        { required: true, message: this.i18n('ip') + this.i18n('is_required_field'), trigger: 'change' },
        { validator: this.$validators.ipv4, trigger: 'change' },
        {
          validator: (rule, value, callback) => {
            for (let key in this.cluster.inventory.all.hosts) {
              if (key !== this.nodeName && this.cluster.inventory.all.hosts[key].ip === value) {
                return callback(this.i18n('duplicateIP', { node: key }))
              }
            }

            // let block = new Netmask(this.cluster.inventory.all.children.target.children.k8s_cluster.vars.calico_ippool_cidr)
            // if (block.contains(value)) {
            //   return callback('IP 不能被包含在容器组子网中 ' + this.cluster.inventory.all.children.target.children.k8s_cluster.vars.calico_ippool_cidr)
            // }

            // block = new Netmask(this.cluster.inventory.all.children.target.children.k8s_cluster.vars.kube_service_cidr)
            // if (block.contains(value)) {
            //   return callback('IP 不能被包含在服务子网中 ' + this.cluster.inventory.all.children.target.children.k8s_cluster.vars.kube_service_cidr)
            // }

            // A类地址：10.0.0.0--10.255.255.255
            // B类地址：172.16.0.0--172.31.255.255 
            // C类地址：192.168.0.0--192.168.255.255
            let block = new Netmask('10.0.0.0/8')
            if (block.contains(value)) {
              return callback()
            }
            block = new Netmask('172.16.0.0/12')
            if (block.contains(value)) {
              return callback()
            }
            block = new Netmask('192.168.0.0/16')
            if (block.contains(value)) {
              return callback()
            }
            return callback('必须绑定到内网 IP')
          },
          trigger: 'change'
        },
      ],
      optionIps: [],
      optionIpsLoading: false,
    }
  },
  setup() {
    const { t } = useI18n({
      inheritLocale: true,
      useScope: 'local'
    })
    return { i18n: t };
  },
  inject: ['onlineNodes', 'isClusterInstalled'],
  computed: {
    ansible_host: {
      get() {
        return this.holderRef.ansible_host
      },
      set(v) {
        this.holderRef.ansible_host = v
        this.holderRef.ip = v
      }
    },
    enableSsh: {
      get() {
        return true
      },
      set(v) {
        console.log(v)
      }
    },
    holderRef: {
      get() { return this.holder || {} },
      set() { }
    },
    ansible_user: {
      get() {
        return this.holder.ansible_user
      },
      set(v) {
        this.holderRef.ansible_user = 'root'
        if (v === 'root') {
          this.ansible_become = false
        } else {
          this.ansible_become = false
        }
      }
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
        if (v) {
          this.holderRef.ansible_become_user = 'root'
          if (!this.holder.ansible_become_password) {
            this.holderRef.ansible_become_password = this.cluster.inventory.all.vars.ansible_become_password || this.holder.ansible_password || this.cluster.inventory.all.children.target.vars.ansible_password
          }
        } else {
          this.holderRef.ansible_become = undefined
          this.holderRef.ansible_become_user = undefined
          this.holderRef.ansible_become_password = undefined
        }
      }
    }
  },
  watch: {
    nodeName(newValue) {
      // 如果节点已存在，则折叠 ssh 配置参数
      this.$refs.configSection.expand(this.onlineNodes[newValue] == undefined)
    },
  },
  components: { SshAddPrivateKey },
  mounted() {
    setTimeout(() => {
      this.$refs.configSection.expand(this.onlineNodes[this.nodeName] == undefined)
    }, 200)
  },
  methods: {
    placeholder(fieldName) {
      let default_value = this.cluster.inventory.all.vars[fieldName]
      if (fieldName.indexOf('password') > 0 && default_value) {
        default_value = '******'
      }
      if (fieldName === 'ansible_ssh_private_key_file') {
        if (default_value && default_value.length > 43) {
          return default_value.slice(43)
        }
      }
      return default_value ? this.i18n('default_value', { default_value: default_value }) : this.t(fieldName + '_placeholder')
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
    async loadOptionIps(flag) {
      if (flag) {
        this.optionIpsLoading = true
        let vars = this.cluster.inventory.all.vars
        let req = {
          from_cache: false,
          ansible_host: this.ansible_host,
          ansible_port: this.holder.ansible_port || vars.ansible_port,
          ansible_user: this.holder.ansible_user || vars.ansible_user,
          ansible_password: this.holder.ansible_password || vars.ansible_password,
          ansible_ssh_private_key_file: this.holder.ansible_ssh_private_key_file || vars.ansible_ssh_private_key_file,
          ansible_become: this.holder.ansible_become !== undefined ? this.holder.ansible_become : vars.ansible_become,
          ansible_become_user: this.holder.ansible_become_user || vars.ansible_become_user,
          ansible_become_password: this.holder.ansible_become_password || vars.ansible_become_password,
          ansible_ssh_common_args: vars.ansible_ssh_common_args,
          gather_subset: '!all,!min,network',
          filter: 'ansible_all_ipv4_addresses',
        }
        await this.pangeeClusterApi.post(`/facts/cluster/${this.cluster.name}/${this.nodeName}`, req).then(resp => {
          if (resp.data.ansible_facts) {
            this.optionIps = resp.data.ansible_facts.ansible_all_ipv4_addresses
          } else {
            this.$message.error(resp.data.msg)
          }
        }).catch(e => {
          let msg = '' + e
          if (e.response) {
            if (e.response.status !== 417) {
              if (e.response && e.response.data && e.response.data.message) {
                msg = e.response.data.message
              }
            }
          }
          this.$message.error(msg)
        })
        this.optionIpsLoading = false
      }
    },
    async loadPythonInterpreter() {
      let result = [
        { label: 'auto', value: 'auto' },
        { label: 'python3.10', value: 'python3.10' },
        { label: 'python3.9', value: 'python3.9' },
        { label: 'python3.8', value: 'python3.8' },
        { label: 'python3.7', value: 'python3.7' },
        { label: 'python3.6', value: 'python3.6' },
        { label: 'python3.5', value: 'python3.5' },
        { label: '/usr/bin/python3', value: '/usr/bin/python3' },
        { label: '/usr/libexec/platform-python', value: '/usr/libexec/platform-python' },
        { label: 'python2.7', value: 'python2.7' },
        { label: 'python2.6', value: 'python2.6' },
        { label: '/usr/bin/python', value: '/usr/bin/python' },
        { label: 'python', value: 'python' },
      ]
      return result
    }
  }
}
</script>

<style scoped lang="css"></style>
