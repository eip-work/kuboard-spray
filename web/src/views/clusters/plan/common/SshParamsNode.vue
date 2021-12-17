<i18n>
en:
  addSshKey: Add Private Key
zh:
  addSshKey: 添加私钥
  ansible_host_placeholder: 'KuboardSpray 连接该主机时所使用的主机名或 IP 地址'
  default_value: '默认值：{default_value} （继承自 K8s 集群标签页中的设置）'
</i18n>


<template>
  <ConfigSection v-model:enabled="enableSsh" label="SSH" :description="description" disabled>
    <FieldString :holder="holder" fieldName="ansible_host" :prop="`all.hosts.${nodeName}`"
      :placeholder="$t('ansible_host_placeholder')" :required="isNode"></FieldString>
    <FieldString :holder="holder" fieldName="ansible_port" :prop="`all.hosts.${nodeName}`"
      :placeholder="placeholder('ansible_port')" 
      :required="!inventory.all.children.k8s_cluster.vars.ansible_port"></FieldString>
    <FieldString :holder="holder" fieldName="ansible_user" :placeholder="placeholder('ansible_user')"></FieldString>
    <FieldSelect :holder="holder" fieldName="ansible_ssh_private_key_file" :loadOptions="loadSshKeyList"
      :placeholder="placeholder('ansible_ssh_private_key_file')">
      <el-button type="primary" style="margin-left: 10px;" icon="el-icon-plus" @click="$refs.addPrivateKey.show()">{{$t('addSshKey')}}</el-button>
    </FieldSelect>
    <FieldString :holder="holder" fieldName="ansible_password" show-password
      :placeholder="placeholder('ansible_password')"></FieldString>
    <FieldBool :holder="holder" fieldName="ansible_become"></FieldBool>
    <template v-if="holder.ansible_become">
      <FieldString :holder="holder" fieldName="ansible_become_user" :placeholder="placeholder('ansible_become_user')"></FieldString>
      <FieldString :holder="holder" fieldName="ansible_become_password" :placeholder="placeholder('ansible_become_password')"></FieldString>
    </template>
    <slot></slot>
    <SshAddPrivateKey ref="addPrivateKey" :clusterName="clusterName"></SshAddPrivateKey>
  </ConfigSection>
</template>

<script>
import ConfigSection from '../ConfigSection.vue'
import SshAddPrivateKey from './SshAddPrivateKey.vue'

export default {
  props: {
    inventory: { type: Object, required: true },
    clusterName: { type: String, required: true },
    nodeName: { type: String, required: false },
    holder: { type: Object, required: true },
    prop: { type: String, required: true },
    description: { type: String, required: true },
    isNode: { type: Boolean, required: false, default: false }
  },
  data() {
    return {

    }
  },
  computed: {
    enableSsh: {
      get () {
        return true
      },
      set (v) {
        console.log(v)
      }
    },
  },
  components: { ConfigSection, SshAddPrivateKey },
  mounted () {
  },
  methods: {
    placeholder(fieldName) {
      let default_value = this.inventory.all.children.k8s_cluster.vars[fieldName]
      if (fieldName.indexOf('password') > 0) {
        default_value = '******'
      }
      return default_value ? this.$t('default_value', {default_value: default_value}) : this.$t('field.' + fieldName + '_placeholder')
    },
    async loadSshKeyList () {
      let result = []
      await this.kuboardSprayApi.get(`/clusters/${this.clusterName}/private-keys`).then(resp => {
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