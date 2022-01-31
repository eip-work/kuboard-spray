<i18n>
en:
  bastionUsage: KuboardSpray can access Kubernetes Cluster Nodes through bastion.
  addSshKey: Add Private Key
  ansible_host_placeholder: 'KuboardSpray use this ip or hostname to connect to the node.'
  default_value: 'Default: {default_value} (inhirit from value configured in Global Config tab)'
  duplicateIP: "IP address conflict with {node}"
  description: Access target nodes throguh bastion/jumpserver【Cannot work now.】
zh:
  bastionUsage: KuboardSpray 可以通过跳板机或堡垒机访问将要安装 K8S 集群的目标节点。
  addSshKey: 添加私钥
  ansible_host_placeholder: 'KuboardSpray 连接该主机时所使用的主机名或 IP 地址'
  default_value: '默认值：{default_value} （继承自全局设置标签页中的配置）'
  duplicateIP: "IP 地址不能与其他节点相同：{node}"
  description: 通过跳板机或者堡垒机访问目标节点【待完善】
</i18n>


<template>
  <ConfigSection ref="configSection" v-model:enabled="bastionEnabled" :label="$t('obj.bastion')" :description="$t('description')" anti-freeze>
    <template #more>
      {{ $t('bastionUsage') }}
    </template>
    <FieldString :holder="holder" fieldName="ansible_host" :prop="`all.hosts.${nodeName}`" anti-freeze
      :placeholder="$t('ansible_host_placeholder')" :rules="hostRules"></FieldString>
    <FieldString :holder="holder" fieldName="ansible_port" :prop="`all.hosts.${nodeName}`"
      :placeholder="placeholder('ansible_port')" anti-freeze required></FieldString>
    <FieldString :holder="holder" fieldName="ansible_user" :prop="`all.hosts.${nodeName}`"
      :placeholder="placeholder('ansible_user')" anti-freeze required></FieldString>
    <FieldSelect :holder="holder" fieldName="ansible_ssh_private_key_file" :loadOptions="loadSshKeyList" anti-freeze clearable
      :placeholder="placeholder('ansible_ssh_private_key_file')">
      <template #edit>
        <el-button type="primary" plain style="margin-left: 10px;" icon="el-icon-plus" @click="$refs.addPrivateKey.show()">{{$t('addSshKey')}}</el-button>
      </template>
    </FieldSelect>
    <FieldString :holder="holder" fieldName="ansible_password" show-password anti-freeze
      :placeholder="placeholder('ansible_password')"></FieldString>
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
  },
  data() {
    return {
      hostRules: [
        { required: true, message: this.$t('field.ansible_host') + this.$t('field.is_required_field'), trigger: 'blur' },
      ]
    }
  },
  computed: {
    inventory: {
      get () { return this.cluster.inventory },
      set () {}
    },
    bastionEnabled: {
      get () {
        return this.inventory.all.hosts.bastion !== undefined
      },
      set (v) {
        if (v) {
          this.inventory.all.hosts.bastion = this.inventory.all.hosts.bastion || {ansible_host: '', ansible_user: ''}
        } else {
          delete this.inventory.all.hosts.bastion
          this.inventory.all.children.target.vars.ansible_ssh_common_args = `-o ControlPath={{ kuboardspray_cluster_dir }}/ansible-{{ansible_user}}@{{ansible_host}}:{{ansible_port}}`
        }
      }
    },
    holderRef: {
      get () {return this.holder},
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
  components: { SshAddPrivateKey },
  mounted () {
  },
  watch: {
    holder: {
      handler: function (bastion) {
        if (bastion && bastion.ansible_host) {
          let sshPass = ''
          if (bastion['ansible_password']) {
            sshPass = `sshpass -p '${bastion['ansible_password']}' `
          }
          let temp = `-o ControlPath={{ kuboardspray_cluster_dir }}/ansible-{{ansible_user}}@{{ansible_host}}:{{ansible_port}} -o ProxyCommand="${sshPass} ssh -F /dev/null -o ControlPath={{ kuboardspray_cluster_dir }}/ansible-%%r@%%h:%%p -o ControlMaster=auto -o ControlPersist=30m -o ConnectTimeout=6 -o ConnectionAttempts=10 -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -W %h:%p -p`
          temp += bastion["ansible_port"] + " " + bastion["ansible_user"] + "@" + bastion["ansible_host"]
          if (bastion["ansible_ssh_private_key_file"] != undefined) {
            temp += " -i " + bastion["ansible_ssh_private_key_file"]
          }
          temp += '"'
          this.inventory.all.children.target.vars.ansible_ssh_common_args = temp
        } else {
          this.inventory.all.children.target.vars.ansible_ssh_common_args = `-o ControlPath={{ kuboardspray_cluster_dir }}/ansible-{{ansible_user}}@{{ansible_host}}:{{ansible_port}}`
        }
      },
      deep: true,
    }
  },
  methods: {
    placeholder(fieldName) {
      return this.$t('field.' + fieldName + '_placeholder')
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
  }
}
</script>

<style scoped lang="scss">

</style>
