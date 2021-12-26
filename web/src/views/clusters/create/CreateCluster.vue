<i18n>
en:
  addCluster: Add Cluster Installation Plan
  createResource: "Add Resource Package"
  name: Name
  requiresName: Please input name.
  conflict: conflict with existing one.
zh:
  addCluster: 添加集群安装计划
  createResource: '添加资源包'
  name: 名称
  requiresName: 请填写名称
  conflict: 与已有的重复 {name}
</i18n>


<template>
  <div>
    <el-dialog v-model="dialogVisible" :close-on-click-modal="false" :modal="true" top="20vh"
      :title="$t('addCluster')" width="45%">
      <el-form :model="form" label-position="left" label-width="120px" v-if="dialogVisible" ref="form">
        <FieldString :holder="form" fieldName="cluster_name" required :placeholder="$t('requiresName')" :rules="nameRules"></FieldString>
        <FieldSelect :holder="form" fieldName="kuboardspray_resource_package" :loadOptions="loadResourceList" required>
          <el-button style="margin-left: 10px;" type="primary" icon="el-icon-plus">{{$t('createResource')}}</el-button>
        </FieldSelect>
        <template v-if="form.kuboardspray_resource_package">
          <FieldSelect :holder="form" fieldName="container_manager" 
            required :loadOptions="loadContainerEngines"></FieldSelect>
          <FieldSelect :holder="form" fieldName="kube_network_plugin" 
            required :loadOptions="loadKubeNetworkPlugin"></FieldSelect>
          <FieldSelect :holder="form" fieldName="etcd_deployment_type" 
            required :loadOptions="loadEtcdDeploymentOptions"></FieldSelect>
        </template>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false" icon="el-icon-close" type="info" plain>{{$t('msg.cancel')}}</el-button>
          <el-button @click="save" icon="el-icon-check" type="primary" :loading="saving">{{$t('msg.ok')}}</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import {ref} from 'vue'

export default {
  setup () {
    return {
      dialogVisible: ref(false),
    }
  },
  props: {
  },
  provide () {
    return {
      editMode: 'edit',
    }
  },
  data() {
    return {
      saving: false,
      form: {
        create: {
          cluster_name: '',
          kuboardspray_resource_package: '',
          container_manager: '',
          kube_network_plugin: '',
          etcd_deployment_type: '',
        }
      },
      nameRules: [
        {
          validator: (rule, value, callback) => {
            if (value.length < 4) {
              return callback('min length ' + 4)
            }
            if (value.length > 20) {
              return callback('max length ' + 20)
            }
            if (!/^[a-zA-Z][a-zA-Z0-9_]{3,21}$/.test(value)) {
              return callback('必须以字母开头，可以包含数字和字母')
            }
            this.kuboardSprayApi.get(`/clusters/${value}`).then(() => {
              callback(this.$t('conflict', {name: value}))
            }).catch(e => {
              console.log(e.response)
              if (e.response && e.response.data.code === 500) {
                callback()
              }
            })
          },
          trigger: 'blur',
        }
      ]
    }
  },
  computed: {
  },
  components: { },
  mounted () {
  },
  watch: {
    'form.kuboardspray_resource_package': function() {
      this.form.container_manager = ''
      this.form.etcd_deployment_type = ''
    },
    'form.container_manager': function () {
      this.form.etcd_deployment_type = ''
    }
  },
  methods: {
    async loadEtcdDeploymentOptions () {
      return [
        {
          label: this.$t('field.etcd_deployment_type-host'),
          value: 'host'
        },
        {
          label: this.$t('field.etcd_deployment_type-docker'),
          value: 'docker',
          disabled: this.form.container_manager !== 'docker',
        }
      ]
    },
    show () {
      this.dialogVisible = true
    },
    async loadResourceList () {
      let result = []
      await this.kuboardSprayApi.get('/resources').then(resp => {
        for (let res of resp.data.data) {
          result.push({ label: res, value: res })
        }
      }).catch(e => {
        console.log(e)
      })
      return result
    },
    async loadContainerEngines () {
      let result = []
      await this.kuboardSprayApi.get(`/resources/${this.form.kuboardspray_resource_package}`).then(resp => {
        let engines = resp.data.data.package.container_engine
        for (let i in engines) {
          let engine = engines[i]
          result.push({
            label: engine.container_manager + (engine.version ? '_' + engine.version : ''),
            value: engine.container_manager,
          })
        }
      }).catch(e => {
        console.log(e)
      })
      return result
    },
    async loadKubeNetworkPlugin () {
      let result = []
      await this.kuboardSprayApi.get(`/resources/${this.form.kuboardspray_resource_package}`).then(resp => {
        let cnies = resp.data.data.package.cni
        for (let i in cnies) {
          let cni = cnies[i]
          result.push({
            label: cni.name + (cni.version ? '_' + cni.version : ''),
            value: cni.name,
          })
        }
      }).catch(e => {
        console.log(e)
      })
      return result
    },
    save () {
      this.$refs.form.validate(async flag => {
        if (flag) {
          this.saving = true
          let req = {
            name: this.form.cluster_name,
            resource_package: this.form.kuboardspray_resource_package,
            container_manager: this.form.container_manager,
            kube_network_plugin: this.form.kube_network_plugin,
            etcd_deployment_type: this.form.etcd_deployment_type,
          }
          await this.kuboardSprayApi.post('/clusters', req).then(resp => {
            console.log(resp.data.data)
            this.$message.success(this.$t('msg.save_succeeded'))
            this.$router.push(`/clusters/${this.form.cluster_name}?mode=edit`)
          }).catch(e => {
            this.$message.error(this.$t('msg.save_failed', {msg: e.response.data.message}))
          })
          this.saving = false
        }
      })
    }
  }
}
</script>

<style scoped lang="scss">

</style>
