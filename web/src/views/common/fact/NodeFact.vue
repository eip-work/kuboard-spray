<i18n>
en:
  validate: Validate Connection
  facts: Host information
  baseInfo: Basic information
  network: Network information
  disk: Disk information
  ansible_machine: CPU Arch
  os: Operation System
  cpumem: CPU/Mem
  no_cached_facts: no cached facts found, please click the Validate Connection button above.
  python_executable: Python Executable
  python_version: Python Version
  choosen: Bind To
zh:
  validate: 验证连接
  facts: 主机信息
  baseInfo: 基本信息
  network: 网络信息
  disk: 磁盘信息
  ansible_machine: CPU 架构
  os: 操作系统
  cpumem: CPU/内存
  no_cached_facts: 未找到该节点的缓存信息，请点击上方的 "验证连接" 按钮
  python_executable: Python 路径
  python_version: Python 版本
  choosen: 绑定到此IP

</i18n>

<template>
  <div>
      <div style="text-align: right;">
        <el-button type="primary" :loading="loadingFact" icon="el-icon-refresh-left" @click="loadFacts(false)">{{$t('validate')}}</el-button>
      </div>
      <el-skeleton v-if="loadingFact" animated></el-skeleton>
      <div v-if="fact" class="app_form_mini app_margin_top">
        <el-form v-if="fact.ansible_facts" label-width="160px" label-position="left">
          <div style="text-align: center; margin-bottom: 10px; margin-top: -10px; font-weight: bold;">[ {{$t('facts')}} ]</div>
          <el-collapse v-model="activeNames">
            <el-collapse-item name="1">
              <template #title>
                <span class="package_title">{{$t('baseInfo')}}</span>
              </template>
              <div class="package_info">
                <NodeFactField :holder="fact.ansible_facts.ansible_python" fieldName="executable" :label="$t('python_executable')"></NodeFactField>
                <NodeFactField :holder="fact.ansible_facts" fieldName="ansible_python_version" :label="$t('python_version')"></NodeFactField>
                <el-divider style="margin: 10px 0;"></el-divider>
                <NodeFactField v-if="fact.ansible_facts.ansible_lsb && fact.ansible_facts.ansible_lsb.description" :holder="fact.ansible_facts.ansible_lsb" fieldName="description" :label="$t('os')"></NodeFactField>
                <el-form-item v-else :label="$t('os')">
                  <span class="field_value app_text_mono">{{fact.ansible_facts.ansible_distribution}} {{fact.ansible_facts.ansible_distribution_version}}</span>
                </el-form-item>
                <NodeFactField :holder="fact.ansible_facts" fieldName="ansible_machine" :label="$t('ansible_machine')"></NodeFactField>
                <el-form-item :label="$t('cpumem')">
                  <span class="field_value app_text_mono">{{fact.ansible_facts.ansible_processor_vcpus}}core / {{Math.round(fact.ansible_facts.ansible_memtotal_mb * 10 / 1024)/10}}Gi</span>
                </el-form-item>
              </div>
            </el-collapse-item>
            <el-collapse-item name="2">
              <template #title>
                <span class="package_title">{{$t('network')}}</span>
              </template>
              <div class="package_info" v-for="(ipv4, index) in ipv4s" :key="'ip'+index">
                <el-divider v-if="index !== 0" style="margin: 10px 0;"></el-divider>
                <NodeFactField :holder="ipv4.ipv4" fieldName="type"></NodeFactField>
                <NodeFactField :holder="ipv4.ipv4" fieldName="device"></NodeFactField>
                <NodeFactField :holder="ipv4.ipv4.ipv4" fieldName="address">
                  <el-tag type="success" effect="dark" v-if="ipv4.isPrefered">{{$t('choosen')}}</el-tag>
                </NodeFactField>
                <!-- <NodeFactField :holder="ipv4.ipv4.ipv4" fieldName="network"></NodeFactField> -->
                <!-- <NodeFactField :holder="ipv4" fieldName="gateway"></NodeFactField> -->
                <!-- <NodeFactField :holder="ipv4.ipv4.ipv4" fieldName="broadcast"></NodeFactField> -->
                <NodeFactField :holder="ipv4.ipv4.ipv4" fieldName="netmask"></NodeFactField>
                <NodeFactField :holder="ipv4.ipv4" fieldName="macaddress"></NodeFactField>
                <!-- <NodeFactField :holder="ipv4.ipv4" fieldName="mtu"></NodeFactField> -->
              </div>
            </el-collapse-item>
            <el-collapse-item name="3">
              <template #title>
                <span class="package_title">{{$t('disk')}}</span>
              </template>
              <div class="package_info">
                <template v-for="(item, key) in fact.ansible_facts.ansible_devices" :key="'disk' + key">
                  <el-form-item :label="key" v-if="item.model">
                    <el-form-item label="vendor" label-width="80px">{{item.vendor}}</el-form-item>
                    <el-form-item label="model" label-width="80px">{{item.model}}</el-form-item>
                    <el-form-item label="size" label-width="80px">{{item.size}}</el-form-item>
                  </el-form-item>
                </template>
                <el-form-item></el-form-item>
              </div>
            </el-collapse-item>
          </el-collapse>
        </el-form>
        <div v-else>
          <el-alert :closable="false" type="error">
            <pre class="app_text_mono">{{fact.msg}}</pre>
            <pre v-if="fact.module_stdout"><el-tag style="margin-right: 5px;">stdout</el-tag>{{fact.module_stdout}}</pre>
            <pre v-if="fact.module_stderr"><el-tag type="danger" effect="dark" style="margin-right: 5px;">stderr</el-tag>{{fact.module_stderr}}</pre>
          </el-alert>
        </div>
      </div>
  </div>
</template>

<script>
import NodeFactField from './NodeFactField.vue'

export default {
  props: [
    'node_owner_type', 'node_owner', 'node_name',
    'ansible_host', 'ansible_port', 'ansible_user', 'ansible_password', 'ansible_ssh_private_key_file',
    'ansible_become', 'ansible_become_user', 'ansible_become_password', 'form', 'ip', 'ansible_ssh_common_args', 'ansible_python_interpreter'
  ],
  data() {
    return {
      loadingFact: false,
      fact: undefined,
      activeNames: ['1'],
    }
  },
  computed: {
    ipv4s () {
      let result = []
      if (this.fact) {
        let defaultIp = this.ip || this.ansible_host
        for (let ip of this.fact.ansible_facts.ansible_all_ipv4_addresses) {
          for (let key in this.fact.ansible_facts) {
            let temp = this.fact.ansible_facts[key]
            if (temp.ipv4 && temp.ipv4.address == ip) {
              result.push({ ipv4: temp, isPrefered: defaultIp === temp.ipv4.address })
            }
          }
        }
      }
      return result
    }
  },
  components: { NodeFactField },
  mounted () {
  },
  methods: {
    clear () {
      this.fact = undefined
    },
    loadFacts(fromCache = true) {
      this.fact = undefined
      if (this.form) {
        this.form.validate(flag => {
          if (flag) {
            this.doLoad(fromCache)
          }
        })
      } else {
        this.doLoad(fromCache)
      }
    },
    async doLoad (fromCache) {
      this.loadingFact = true
      let _this = this
      let req = {
        from_cache: fromCache,
        ansible_host: this.ansible_host,
        ansible_port: this.ansible_port,
        ansible_user: this.ansible_user,
        ansible_password: this.ansible_password,
        ansible_ssh_private_key_file: this.ansible_ssh_private_key_file,
        ansible_become: this.ansible_become,
        ansible_become_user: this.ansible_become_user,
        ansible_become_password: this.ansible_become_password,
        ansible_ssh_common_args: this.ansible_ssh_common_args,
        ansible_python_interpreter: this.ansible_python_interpreter,
        gather_subset: '!all,network,hardware',
      }
      await this.kuboardSprayApi.post(`/facts/${this.node_owner_type}/${this.node_owner}/${this.node_name}`, req).then(resp => {
        if (fromCache) {
          if (resp.data.ansible_facts !== undefined) {
            this.fact = resp.data
          }
        } else {
          this.fact = resp.data
        }
      }).catch(e => {
        if (e.response.status !== 417) {
          let msg = '' + e
          if (e.response && e.response.data && e.response.data.message){
            msg = e.response.data.message
          }
          this.fact = {
            changed: false,
            msg: msg,
          }
        } else {
          this.fact = {
            changed: false,
            msg: _this.$t('no_cached_facts')
          }
        }
      })
      this.loadingFact = false
    }
  }
}
</script>

<style scoped lang="css">
.package_title {
  font-weight: bolder;
}
.package_info {
  margin-left: 20px;
}
.field_value {
  font-size: 13px;
}
</style>
