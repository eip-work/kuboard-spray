<i18n>
en:
  addSshKey: Add Private Key
zh:
  addSshKey: 添加私钥
  ansible_host_placeholder: '必须在节点级别设置'
</i18n>


<template>
  <ConfigSection v-model:enabled="enableSsh" label="SSH" :description="description" disabled>
    <FieldString disabled :holder="holder" fieldName="ansible_host" :prop="isNode ? `all.hosts.${nodeName}` : ''"
      :placeholder="$t('ansible_host_placeholder')"></FieldString>
    <FieldString :holder="holder" fieldName="ansible_port"></FieldString>
    <FieldString :holder="holder" fieldName="ansible_user"></FieldString>
    <FieldSelect :holder="holder" fieldName="ansible_ssh_private_key_file" :loadOptions="loadSshKeyList">
      <el-button type="primary" plain style="margin-left: 10px;" icon="el-icon-plus" @click="$refs.addPrivateKey.show()">{{$t('addSshKey')}}</el-button>
    </FieldSelect>
    <FieldString :holder="holder" fieldName="ansible_password" show-password></FieldString>
    <FieldBool :holder="holder" fieldName="ansible_become"></FieldBool>
    <template v-if="holder.ansible_become">
      <FieldString :holder="holder" fieldName="ansible_become_user"></FieldString>
      <FieldString :holder="holder" fieldName="ansible_become_password"></FieldString>
    </template>
    <slot></slot>
    <SshAddPrivateKey ref="addPrivateKey" ownerType="cluster" :ownerName="cluster.name"></SshAddPrivateKey>
  </ConfigSection>
</template>

<script>
import SshAddPrivateKey from './SshAddPrivateKey.vue'

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
    }
  },
  components: { SshAddPrivateKey },
  mounted () {
  },
  methods: {
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
