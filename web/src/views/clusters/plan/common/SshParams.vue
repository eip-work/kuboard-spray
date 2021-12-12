<i18n>
en:
  addSshKey: Add Private Key
zh:
  addSshKey: 添加私钥
</i18n>


<template>
  <ConfigSection v-model:enabled="enableSsh" label="SSH" :description="description">
    <FieldString v-if="isNode" :holder="holder" fieldName="ansible_host"></FieldString>
    <FieldString :holder="holder" fieldName="ansible_user"></FieldString>
    <FieldSelect :holder="holder" fieldName="ansible_ssh_private_key_file" :loadOptions="loadSshKeyList">
      <el-button type="primary" style="margin-left: 10px;" icon="el-icon-plus" @click="$refs.addPrivateKey.show()">{{$t('addSshKey')}}</el-button>
    </FieldSelect>
    <FieldString :holder="holder" fieldName="ansible_password" show-password></FieldString>
    <FieldBool :holder="holder" fieldName="ansible_become"></FieldBool>
    <template v-if="holder.ansible_become">
      <FieldString :holder="holder" fieldName="ansible_become_user"></FieldString>
      <FieldString :holder="holder" fieldName="ansible_become_password"></FieldString>
    </template>
    <SshAddPrivateKey ref="addPrivateKey" :clusterName="clusterName"></SshAddPrivateKey>
  </ConfigSection>
</template>

<script>
import ConfigSection from '../ConfigSection.vue'
import SshAddPrivateKey from './SshAddPrivateKey.vue'

export default {
  props: {
    clusterName: { type: String, required: true },
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
  components: { ConfigSection, SshAddPrivateKey },
  mounted () {
  },
  methods: {
    async loadSshKeyList () {
      let result = []
      await this.kuboardSprayApi.get(`/clusters/${this.clusterName}/private-keys`).then(resp => {
        for (let item of resp.data.data) {
          result.push({ label: item, value: '{{ kuboardspray_cluster_dir }}/private-key/' + item })
        }
      })
      return result
    }
  }
}
</script>

<style scoped lang="scss">

</style>
