<i18n>
en:
  addSshKey: Add Private Key
zh:
  addSshKey: 添加私钥
  ansible_host_placeholder: 'KuboardSpray 连接该主机时所使用的主机名或 IP 地址'
  default_value: '默认值：{default_value} （继承自 K8s 集群标签页中的设置）'
</i18n>


<template>
  <ConfigSection v-model:enabled="enableSsh" label="SSH" :description="description" disabled anti-freeze>
    <FieldString :holder="holder" fieldName="ansible_host" :prop="`all.hosts.${nodeName}`" anti-freeze
      :placeholder="$t('ansible_host_placeholder')" :required="isNode"></FieldString>
    <FieldString :holder="holder" fieldName="ansible_port" :prop="`all.hosts.${nodeName}`"
      :placeholder="placeholder('ansible_port')" anti-freeze
      :required="!cluster.inventory.all.children.target.vars.ansible_port"></FieldString>
    <FieldString :holder="holder" fieldName="ansible_user" :placeholder="placeholder('ansible_user')" anti-freeze></FieldString>
    <FieldSelect :holder="holder" fieldName="ansible_ssh_private_key_file" :loadOptions="loadSshKeyList" anti-freeze
      :placeholder="placeholder('ansible_ssh_private_key_file')">
      <template #edit>
        <el-button type="primary" plain style="margin-left: 10px;" icon="el-icon-plus" @click="$refs.addPrivateKey.show()">{{$t('addSshKey')}}</el-button>
      </template>
    </FieldSelect>
    <FieldString :holder="holder" fieldName="ansible_password" show-password anti-freeze
      :placeholder="placeholder('ansible_password')"></FieldString>
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
      <FieldString :holder="holder" fieldName="ansible_become_password" :placeholder="placeholder('ansible_become_password')" anti-freeze></FieldString>
    </template>
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
  methods: {
    placeholder(fieldName) {
      let default_value = this.cluster.inventory.all.children.target.vars[fieldName]
      if (fieldName.indexOf('password') > 0 && default_value) {
        default_value = '******'
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
  }
}
</script>

<style scoped lang="scss">

</style>
