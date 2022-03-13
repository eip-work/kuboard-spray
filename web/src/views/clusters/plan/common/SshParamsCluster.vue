<i18n>
en:
  addSshKey: Add Private Key
  ansible_host_placeholder: 'Should provide at node level.'
  password_and_bastion: You are using bastion/jumpserver to access the node, if you use password here, it takes longer time to create ssh connection, it would be better to clear the password and provide ssh_private_key.
  speedup: SSH Multiplexing
  rootuser: Must use root user
zh:
  addSshKey: 添加私钥
  ansible_host_placeholder: '必须在节点级别设置'
  password_and_bastion: 您将使用跳板机访问节点，如果在此处使用密码访问，创建 ssh 连接的时间较长，建议清除密码并提供 “私钥文件”。
  speedup: 加速执行
  rootuser: 必须使用 root 用户
</i18n>


<template>
  <ConfigSection v-model:enabled="enableSsh" label="SSH" :description="description" disabled anti-freeze>
    <FieldString disabled :holder="holder" fieldName="ansible_host" :prop="isNode ? `all.hosts.${nodeName}` : ''"
      :placeholder="$t('ansible_host_placeholder')"></FieldString>
    <FieldString :holder="holder" fieldName="ansible_port" anti-freeze></FieldString>
    <FieldCommon :holder="holder" fieldName="ansible_user" anti-freeze>
      <template #edit>
        <el-input v-model.trim="ansible_user"></el-input>
        <el-tag class="app_text_mono" style="display: block; line-height: 18px;" type="warning">{{ $t('rootuser') }}</el-tag>
      </template>
    </FieldCommon>
    <FieldSelect :holder="holder" fieldName="ansible_ssh_private_key_file" :loadOptions="loadSshKeyList" anti-freeze clearable>
      <template #edit>
        <el-button type="primary" plain style="margin-left: 10px;" icon="el-icon-plus" @click="$refs.addPrivateKey.show()">{{$t('addSshKey')}}</el-button>
      </template>
    </FieldSelect>
    <FieldString :holder="holder" fieldName="ansible_password" anti-freeze show-password clearable></FieldString>
    <el-alert type="warning" :closable="false" v-if="cluster && cluster.inventory.all.hosts.bastion && holder.ansible_password" style="margin-left: 120px; width: calc(100% - 120px);">
      {{ $t('password_and_bastion') }}
      <KuboardSprayLink href="https://kuboard-spray.cn/guide/extra/speedup.html" style="margin-left: 10px;" :size="12"></KuboardSprayLink>
    </el-alert>
    <!-- <FieldCommon :holder="holder" fieldName="ansible_become" anti-freeze>
      <template #view>
        <el-switch v-model="ansible_become" disabled></el-switch>
      </template>
      <template #edit>
        <el-switch v-model="ansible_become" disabled></el-switch>
      </template>
    </FieldCommon>
    <template v-if="holder.ansible_become">
      <FieldString :holder="holder" fieldName="ansible_become_user" anti-freeze disabled></FieldString>
      <FieldString :holder="holder" fieldName="ansible_become_password" show-password anti-freeze clearable></FieldString>
      <div v-if="editMode !== 'view'">{{ $t('become_password_desc', {ansible_user}) }}</div>
    </template> -->
    <FieldSelect :holder="holder" fieldName="ansible_python_interpreter" anti-freeze clearable :loadOptions="loadPythonInterpreter" allow-create filterable></FieldSelect>
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
  inject: ['editMode'],
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
      get () {return this.holder || {}},
      set () {}
    },
    ansible_user: {
      get () {
        return this.holder.ansible_user
      },
      set (v) {
        this.holderRef.ansible_user = 'root'
        if (v === 'root') {
          this.ansible_become = false
        } else {
          this.ansible_become = false
        }
      }
    },
    ansible_become: {
      get () {
        return this.holder.ansible_become
      },
      set (v) {
        this.holderRef.ansible_become = v
        if (v) {
          this.holderRef.ansible_become_user = 'root'
          if (!this.holder.ansible_become_password) {
            this.holderRef.ansible_become_password = this.holder.ansible_password
          }
        } else {
          this.holderRef.ansible_become = undefined
          this.holderRef.ansible_become_user = undefined
          this.holderRef.ansible_become_password = undefined
        }
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
