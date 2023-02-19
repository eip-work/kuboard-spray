<i18n>
en:
  add_insecure_registry: Add insecure registry
  edit_insecure_registry: Edit insecure registry
  insecure_registries_placeholder: e.g. http://192.168.30.56:5000
  domain_name: Domain name of the registry
  domain_name_required: Domain name is required
  address: Address of the registry
zh:
  add_insecure_registry: 添加 http 协议镜像仓库（或自签名 https 的镜像仓库）
  edit_insecure_registry: 编辑 http 协议镜像仓库（或自签名 https 的镜像仓库）
  insecure_registries_placeholder: 例如：http://192.168.30.56:5000
  domain_name: 镜像仓库的域名
  domain_name_required: 镜像仓库的域名为必填字段
  address: 镜像仓库的地址
</i18n>

<template>
  <el-popover :visible="visible" placement="top-end" :width="630" trigger="manual">
    <template #reference>
      <el-button v-if="isAdd" type="primary" plain icon="el-icon-plus" @click="visible = true">{{$t('add_insecure_registry')}}</el-button>
      <el-button v-else type="primary" plain icon="el-icon-edit" @click="visible = true" style="margin-top: -7px;">{{$t('msg.edit')}}</el-button>
    </template>
    <div v-if="isAdd" class="app_block_title">{{ $t('add_insecure_registry') }}</div>
    <div v-else class="app_block_title">{{ $t('edit_insecure_registry') }}</div>
    <div style="margin-top: 20px;">
      <el-form ref="form" :model="insecure_registry_form" @submit.prevent.stop>
        <div v-if="isAdd" class="app_margin_bottom description" style="line-height: 28px; margin-bottom: 20px;">
          <li>假设镜像仓库的访问地址为 <el-tag class="app_text_mono">192.168.30.56:5000</el-tag> </li>
          <li>假设引用镜像的格式为 <el-tag class="app_text_mono">my-registry.com/repo-name/image-name:tag</el-tag> </li>
        </div>
        <el-form-item :label="$t('domain_name')" prop="domain_name" :rules="domainNameRules" label-width="120px">
          <div class="description">
            <el-input v-model.trim="insecure_registry_form.domain_name" placeholder="基于上述假设，此处应填写 my-registry.com"></el-input>
            <div v-if="isAdd">
              请将 my-registry.com 替换成你镜像仓库的域名
            </div>
          </div>
        </el-form-item>
        <el-form-item :label="$t('address')" prop="address" :rules="addressRules" label-width="120px">
          <div class="description">
            <el-input v-model.trim="insecure_registry_form.address" placeholder="基于上述假设，此处应填写 http(s)://192.168.30.56:5000"></el-input>
            <div v-if="isAdd">
              请将 192.168.30.56:5000 替换成你镜像仓库的地址及端口
            </div>
          </div>
        </el-form-item>
        <div v-if="isAdd" class="description">
          <li>
            完成配置后，在任意一个 K8S 节点执行以下命令，以验证配置是否生效：
            <pre style="margin-left: 18px; background-color: #333; color: white; padding: 5px 10px; border-radius: 3px;">crictl pull my-registry.com/repo-name/image-name:tag</pre>
          </li>
        </div>
        <div style="text-align: right; margin-top: 10px;">
          <el-button icon="el-icon-close" @click="visible = false">{{ $t('msg.cancel') }}</el-button>
          <el-button type="primary" icon="el-icon-check" @click="add">{{ $t('msg.ok') }}</el-button>
        </div>
      </el-form>
    </div>
  </el-popover>
</template>

<script>
export default {
  props: {
    cluster: { type: Object, required: true },
    name: { type: String, required: false, default: undefined },
  },
  data () {
    return {
      visible: false,
      insecure_registry_form: {
        domain_name: undefined,
        address: undefined,
      },
      domainNameRules: [
        { required: true, message: this.$t('domain_name_required'), trigger: 'blur' }
      ],
      addressRules: [
        { required: true, type: 'string', message: '字段不能为空', trigger: 'blur' },
      ]
    }
  },
  computed: {
    isAdd () {
      return this.name == undefined
    },
    vars: {
      get () { return this.cluster.inventory.all.children.target.vars },
      set () {},
    },
    prop() { return 'all.children.target.vars' }
  },
  methods: {
    add () {
      this.$refs.form.validate(flag => {
        if (flag) {
          this.vars.containerd_insecure_registries = this.vars.containerd_insecure_registries || {}
          this.vars.containerd_insecure_registries[this.insecure_registry_form.domain_name] = this.insecure_registry_form.address
          this.visible = false
        }
      })
    }
  }
}
</script>

<style>

</style>