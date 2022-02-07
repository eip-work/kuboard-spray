<i18n>
en:
  addSshKey: Add Private Key
  ansible_host_placeholder: 'KuboardSpray use this ip or hostname to connect to the node.'
  default_value: 'Default: {default_value} (inhirit from value configured in Global Config tab)'
  duplicateIP: "IP address conflict with {node}"
  ip: Bind to IP
  ip_placeholder: 'Default: {default_value} (inhirit from value ansible_host)'
  longTimeLoading: Loading... cost about 20-30s
  password_and_bastion: You are using bastion/jumpserver to access the node, if you use password here, it takes longer time to create ssh connection, it would be better to clear the password and provide ssh_private_key.
  password_in_global: You defined password in global configuration, it's suggested that you clear it.
  speedup: SSH Multiplexing
zh:
  addSshKey: 添加私钥
  ansible_host_placeholder: 'KuboardSpray 连接该主机时所使用的主机名或 IP 地址'
  default_value: '默认值：{default_value} （继承自全局设置标签页中的配置）'
  duplicateIP: "IP 地址不能与其他节点相同：{node}"
  ip: 绑定到 IP
  ip_placeholder: '默认值：{default_value} （默认与主机 IP 相同）'
  longTimeLoading: 加载中... 可能需要 20-30 秒
  ipDescription: kubernetes 使用的地址有可能与 kuboard-spray 访问该机器时所使用的地址不同，此处指定 kubernetes 所使用的地址（必须为内网地址）
  password_and_bastion: 您将使用跳板机访问节点，如果在此处使用密码访问，创建 ssh 连接的时间较长，建议清除密码并提供 “私钥文件”。
  password_in_global: 您在全局参数中配置了密码，建议清除。
  speedup: 加速执行
</i18n>


<template>
  <ConfigSection ref="configSection" v-model:enabled="enableSsh" label="SSH" :description="description" disabled anti-freeze>
    <FieldCommon :holder="holder" fieldName="ansible_host" :prop="`all.hosts.${nodeName}`" :anti-freeze="onlineNodes[nodeName] === undefined || holder.kuboardspray_node_action === 'add_node'" :rules="hostRules">
      <template #edit>
        <el-input v-model.trim="ansible_host" :placeholder="$t('ansible_host_placeholder')"></el-input>
      </template>
    </FieldCommon>
    <FieldString :holder="holder" fieldName="ansible_port" :prop="`all.hosts.${nodeName}`"
      :placeholder="placeholder('ansible_port')" anti-freeze
      :required="!cluster.inventory.all.children.target.vars.ansible_port"></FieldString>
    <FieldString :holder="holder" fieldName="ansible_user" :placeholder="placeholder('ansible_user')" anti-freeze></FieldString>
    <FieldSelect :holder="holder" fieldName="ansible_ssh_private_key_file" :loadOptions="loadSshKeyList" anti-freeze clearable
      :placeholder="placeholder('ansible_ssh_private_key_file')">
      <template #edit>
        <el-button type="primary" plain style="margin-left: 10px;" icon="el-icon-plus" @click="$refs.addPrivateKey.show()">{{$t('addSshKey')}}</el-button>
      </template>
    </FieldSelect>
    <FieldString :holder="holder" fieldName="ansible_password" show-password anti-freeze clearable
      :placeholder="placeholder('ansible_password')"></FieldString>
    <el-alert type="warning" :closable="false" v-if="cluster && cluster.inventory.all.hosts.bastion && (holder.ansible_password || cluster.inventory.all.children.target.vars.ansible_password)" style="margin-left: 120px; width: calc(100% - 120px);">
      {{ $t('password_and_bastion') }}
      <KuboardSprayLink href="https://kuboard-spray.cn/guide/extra/speedup.html" style="margin-left: 10px;" :size="12"></KuboardSprayLink>
      <li v-if="cluster && cluster.inventory.all.children.target.vars.ansible_password">{{ $t('password_in_global') }}</li>
    </el-alert>
    <FieldCommon :holder="holder" fieldName="ansible_become" anti-freeze>
      <template #view>
        <el-switch v-model="ansible_become" disabled></el-switch>
      </template>
      <template #edit>
        <el-switch v-model="ansible_become"></el-switch>
      </template>
    </FieldCommon>
    <template v-if="ansible_become">
      <FieldString :holder="holder" fieldName="ansible_become_user" :placeholder="placeholder('ansible_become_user')" anti-freeze></FieldString>
      <FieldString :holder="holder" fieldName="ansible_become_password" :placeholder="placeholder('ansible_become_password')" anti-freeze clearable></FieldString>
    </template>
    <FieldSelect :holder="holder" fieldName="ansible_python_interpreter" anti-freeze clearable :loadOptions="loadPythonInterpreter" :placeholder="placeholder('ansible_python_interpreter')" allow-create filterable></FieldSelect>
    <FieldCommon :holder="holder" fieldName="ip" :prop="`all.hosts.${nodeName}`" :anti-freeze="onlineNodes[nodeName] === undefined" :label="$t('ip')"
      :placeholder="$t('ip_placeholder', { default_value: holder.ansible_host})">
      <template #edit>
        <el-select v-model="holderRef.ip" style="width: 100%;" :loading="optionIpsLoading" @visible-change="loadOptionIps" :loading-text="$t('longTimeLoading')">
          <el-option v-for="item in optionIps" :key="item" :value="item">
            <span class="app_text_mono">{{ item }}</span>
          </el-option>
        </el-select>
        <div style="font-size: 12px; color: #666; line-height: 20px; margin-top: 2px;">{{ $t('ipDescription') }}</div>
      </template>
      <template #view>
        <span class="app_text_mono">{{ holderRef.ip || ansible_host }}</span>
      </template>
    </FieldCommon>
    <slot></slot>
    <SshAddPrivateKey ref="addPrivateKey" ownerType="cluster" :ownerName="cluster.name"></SshAddPrivateKey>
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
    description: { type: String, required: true },
  },
  data() {
    return {
      hostRules: [
        { required: true, message: this.$t('field.ansible_host') + this.$t('field.is_required_field'), trigger: 'blur' },
        {
          validator: (rule, value, callback) => {
            for (let key in this.cluster.inventory.all.hosts) {
              if (key !== this.nodeName && this.cluster.inventory.all.hosts[key].ansible_host === value) {
                return callback(this.$t('duplicateIP', {node: key}))
              }
            }
            return callback()
          },
          trigger: 'blur'
        },
      ],
      optionIps: [],
      optionIpsLoading: false,
    }
  },
  inject: ['onlineNodes', 'isClusterInstalled'],
  computed: {
    ansible_host: {
      get () {
        return this.holderRef.ansible_host
      },
      set (v) {
        this.holderRef.ansible_host = v
        this.holderRef.ip = v
      }
    },
    enableSsh: {
      get () {
        return true
      },
      set (v) {
        console.log(v)
      }
    },
    holderRef: {
      get () {return this.holder || {}},
      set () {}
    },
    ansible_become: {
      get () {
        if (this.holder.ansible_become !== undefined) {
          return this.holder.ansible_become
        }
        return this.cluster.inventory.all.children.target.vars.ansible_become 
      },
      set (v) {
        this.holderRef.ansible_become = v
      }
    }
  },
  watch: {
    nodeName (newValue) {
      // 如果节点已存在，则折叠 ssh 配置参数
      this.$refs.configSection.expand(this.onlineNodes[newValue] == undefined)
    },
  },
  components: { SshAddPrivateKey },
  mounted () {
    setTimeout(() => {
      this.$refs.configSection.expand(this.onlineNodes[this.nodeName] == undefined)
    }, 200)
  },
  methods: {
    placeholder(fieldName) {
      let default_value = this.cluster.inventory.all.children.target.vars[fieldName]
      if (fieldName.indexOf('password') > 0 && default_value) {
        default_value = '******'
      }
      if (fieldName === 'ansible_ssh_private_key_file') {
        if (default_value && default_value.length > 43) {
          return default_value.slice(43)
        }
      }
      return default_value ? this.$t('default_value', {default_value: default_value}) : this.$t('field.' + fieldName + '_placeholder')
    },
    async loadSshKeyList () {
      let result = []
      await this.kuboardSprayApi.get(`/private-keys/cluster/${this.cluster.name}`).then(resp => {
        for (let item of resp.data.data) {
          result.push({ label: item, value: '{{ kuboardspray_cluster_dir }}/private-key/' + item })
        }
      })
      return result
    },
    async loadOptionIps (flag) {
      if (flag) {
        this.optionIpsLoading = true
        let vars = this.cluster.inventory.all.children.target.vars
        let req = {
          from_cache: false,
          ansible_host: this.ansible_host,
          ansible_port: this.holder.ansible_port || vars.ansible_port,
          ansible_user: this.holder.ansible_user || vars.ansible_user,
          ansible_password: this.holder.ansible_password || vars.ansible_password,
          ansible_ssh_private_key_file: this.holder.ansible_ssh_private_key_file || vars.ansible_ssh_private_key_file,
          ansible_become: this.holder.ansible_become || vars.ansible_become,
          ansible_become_user: this.holder.ansible_become_user || vars.ansible_become_user,
          ansible_become_password: this.holder.ansible_become_password || vars.ansible_become_password,
          ansible_ssh_common_args: vars.ansible_ssh_common_args,
          gather_subset: '!all,!min,network',
          filter: 'ansible_all_ipv4_addresses',
        }
        await this.kuboardSprayApi.post(`/facts/cluster/${this.cluster.name}/${this.nodeName}`, req).then(resp => {
          if (resp.data.ansible_facts) {
            this.optionIps = resp.data.ansible_facts.ansible_all_ipv4_addresses
          } else {
            this.$message.error(resp.data.msg)
          }
        }).catch(e => {
          let msg = '' + e
          if (e.response) {
            if (e.response.status !== 417) {
              if (e.response && e.response.data && e.response.data.message){
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

<style scoped lang="css">

</style>
