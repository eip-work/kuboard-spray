<i18n>
en:
  addSshKey: Add Private Key
zh:
  addSshKey: 添加私钥
  ansible_host_placeholder: '必须在节点级别设置'
</i18n>


<template>
  <div>
    <FieldString :holder="holder" fieldName="ansible_host" :prop="prop"
      :placeholder="$t('ansible_host_placeholder')"></FieldString>
    <FieldString :holder="holder" fieldName="ansible_port" :prop="prop"></FieldString>
    <FieldString :holder="holder" fieldName="ansible_user" :prop="prop"></FieldString>
    <FieldSelect :holder="holder" fieldName="ansible_ssh_private_key_file" :loadOptions="loadSshKeyList">
      <el-button type="primary" plain style="margin-left: 10px;" icon="el-icon-plus" @click="$refs.addPrivateKey.show()">{{$t('addSshKey')}}</el-button>
    </FieldSelect>
    <FieldString :holder="holder" fieldName="ansible_password" show-password :prop="prop"></FieldString>
    <FieldBool :holder="holder" fieldName="ansible_become" :prop="prop"></FieldBool>
    <template v-if="holder.ansible_become">
      <FieldString :holder="holder" fieldName="ansible_become_user" :prop="prop"></FieldString>
      <FieldString :holder="holder" fieldName="ansible_become_password" :prop="prop"></FieldString>
    </template>
    <slot></slot>
    <SshAddPrivateKey ref="addPrivateKey" ownerType="mirror" :ownerName="mirrorName"></SshAddPrivateKey>
  </div>
</template>

<script>
import SshAddPrivateKey from '../../private_key/SshAddPrivateKey.vue'

export default {
  props: {
    holder: { type: Object, required: true },
    prop: { type: String, required: true },
    mirrorName: { type: String, required: true },
  },
  data() {
    return {

    }
  },
  computed: {
  },
  components: { SshAddPrivateKey },
  mounted () {
  },
  methods: {
    async loadSshKeyList () {
      let result = []
      await this.kuboardSprayApi.get(`/private-keys/mirror/${this.mirrorName}`).then(resp => {
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
