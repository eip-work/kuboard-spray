<i18n>
en:
  label: Kuboard
  description: kuboard
  addon_function: A cool kubernetes dashboard.
zh:
  label: Kuboard
  description: kuboard
  addon_function: Kubernetes 管理界面。
</i18n>

<template>
  <AddonSection v-model:enabled="enabled" :label="$t('label')" :description="$t('obj.addon', {name: this.$t('description')})"
    :cluster="cluster" addonName="kuboard" @refresh="$emit('refresh')">
    <template #more>
      {{$t('addon_function')}}
    </template>
    <FieldCommon :holder="holder" :prop="prop" fieldName="kuboard_version" :placeholder="'默认值 ' + addon.params_default.kuboard_version">
      <template #edit>
        <el-input v-model.trim="kuboard_version" :placeholder="'默认值 ' + addon.params_default.kuboard_version"></el-input>
        <el-alert v-if="kuboard_version && kuboard_version != addon.params_default.kuboard_version" :closable="false" type="warning">
          您指定的版本与资源包中的 kuboard 版本不匹配，安装时将自动从 docker hub 拉取 eipwork/kuboard:{{kuboard_version}} 镜像，如果您的服务器不能访问外网，安装将会失败。
        </el-alert>
      </template>
    </FieldCommon>
    <FieldNumber :holder="vars" :prop="prop" fieldName="kuboard_port" :placeholder="'默认值 ' + addon.params_default.kuboard_port"></FieldNumber>
    <FieldString :holder="vars" :prop="prop" fieldName="kuboard_data_dir" :placeholder="'默认值 ' + addon.params_default.kuboard_data_dir">
    </FieldString>
    <FieldString :holder="vars" :prop="prop" fieldName="kuboard_cluster_name" :rules="cluster_name_rules" :placeholder="'默认值 ' + addon.params_default.kuboard_cluster_name + '，显示在 Kuboard 中的集群名称'">
    </FieldString>
    <el-form-item label-position="left" label="Kuboard 界面" v-if="cluster.state && cluster.state.addons && cluster.state.addons.kuboard.is_installed" class="app_text_mono">
      <div style="background-color: var(--el-color-success-light-9); padding: 10px; border-radius: 5px;">
        <KuboardSprayLink :href="kuboard_url">{{ kuboard_url }}</KuboardSprayLink>
        <div><div style="display: inline-block; width: 72px;">默认用户名</div>:  admin</div>
        <div><div style="display: inline-block; width: 72px;">默认密 码</div>: Kuboard123</div>
      </div>
    </el-form-item>
  </AddonSection>
</template>

<script>
import AddonSection from '../AddonSection.vue'

export default {
  props: {
    cluster: { type: Object, required: true },
    addon: { type: Object, required: true },
  },
  data() {
    return {
      cluster_name_rules: [
          { required: false, message: '请输入集群名称', trigger: 'blur' },
          { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' },
          {
            validator: (rule, value, callback) => {
              if (!value) {
                return callback()
              }
              if (value === 'GLOBAL') {
                return callback(new Error('GLOBAL 为保留字'))
              }
              return callback()
            },
            trigger: 'blur'
          },
          { validator: this.$validators.dnsName, trigger: 'blur' },
        ],
    }
  },
  computed: {
    vars: {
      get () { return this.cluster.inventory.all.children.target.children.k8s_cluster.vars },
      set () {},
    },
    prop () {
      return 'all.children.target.children.k8s_cluster.vars'
    },
    holder: {
      get () { return this },
      set () {},
    },
    kuboard_version: {
      get () {
        return this.vars.kuboard_version
      },
      set (v) {
        if (v) {
          this.vars.kuboard_version = v
        } else {
          delete(this.vars.kuboard_version)
        }
      }
    },
    kuboard_url () {
      let result = 'http://'
      for (let key in this.cluster.inventory.all.children.target.children.k8s_cluster.children.kube_control_plane.hosts) {
        result += this.cluster.inventory.all.hosts[key].ip || this.cluster.inventory.all.hosts[key].ansible_host
        break
      }
      let addon = undefined
      for (let key in this.cluster.resourcePackage.data.addon) {
        addon = this.cluster.resourcePackage.data.addon[key]
        if (addon.name === 'kuboard')
          break
      }
      let kuboard_port = this.cluster.inventory.all.children.target.children.k8s_cluster.vars.kuboard_port || addon.params_default.kuboard_port
      if (kuboard_port != 80) {
        result += ':' + kuboard_port
      }
      return result
    },
    enabled: {
      get () {
        return this.vars.kuboard_enabled || false
      },
      set (v) {
        this.vars.kuboard_enabled = v
      }
    }
  },
  components: { AddonSection },
  mounted () {
  },
  methods: {

  }
}
</script>

<style scoped lang="css">

</style>
