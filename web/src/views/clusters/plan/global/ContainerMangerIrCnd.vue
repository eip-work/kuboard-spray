<i18n>
en:
  add_insecure_registry: Add insecure registry
  insecure_registries_placeholder: e.g. http://192.168.30.56:5000
  domain_name: Domain name of the registry
  domain_name_required: Domain name is required
zh:
  add_insecure_registry: 添加 http 协议镜像仓库（或自签名 https 的镜像仓库）
  insecure_registries_placeholder: 例如：http://192.168.30.56:5000
  domain_name: 镜像仓库的域名
  domain_name_required: 镜像仓库的域名为必填字段
</i18n>

<template>
  <FieldCommon :holder="vars" :prop="prop" fieldName="containerd_insecure_registries" anti-freeze>
    <template #edit>
      <template v-for="(item, key) in vars.containerd_insecure_registries" :key="'cir' + key">
        <div class="insecure_registry">
          <div effect="dark" class="app_text_mono header">
            {{key}}
            <el-button type="text" style="float: right; color: white;" icon="el-icon-delete" @click="delete vars.containerd_insecure_registries[key]">{{ $t('msg.delete') }}</el-button>
          </div>
          <FieldArray :holder="vars.containerd_insecure_registries" :prop="null" newItemTemplate="" :fieldName="key" label-width="0px"
            label=" " :itemRules="insecureRegistriesItemRules" anti-freeze>
            <template v-slot:editItem="scope">
              <el-input v-model.trim="item[scope.index]" :placeholder="$t('insecure_registries_placeholder')"></el-input>
            </template>
          </FieldArray>
        </div>
      </template>
      <el-popover :visible="visible" placement="top-end" :width="600" trigger="manual">
        <template #reference>
          <el-button type="primary" plain icon="el-icon-plus" @click="visible = true">{{$t('add_insecure_registry')}}</el-button>
        </template>
        <div class="app_block_title">{{ $t('add_insecure_registry') }}</div>
        <div>
          <el-form ref="form" :model="insecure_registry_form" @submit.prevent.stop>
            <el-form-item :label="$t('domain_name')" prop="domain_name" :rules="domainNameRules">
              <div class="description">
                <li>
                  如果您的镜像仓库地址为 <el-tag class="app_text_mono">192.168.30.56:5000</el-tag> 且没有配置 https 证书（或者配置了自签名证书），同时，您在定义容器组时使用 <el-tag class="app_text_mono">my-registry.com/registry-name/image-name:tag</el-tag> 格式引用镜像，
                  则请在此处填写 <el-tag type="success" class="app_text_mono">my-registry.com</el-tag>
                  并在点击确定后添加 <el-tag class="app_text_mono">http://192.168.30.56:5000</el-tag> 或 <el-tag class="app_text_mono">https://192.168.30.56:5000</el-tag>
                </li>
                <li>
                  完成配置后，在任意一个 K8S 节点执行以下命令，以验证配置是否生效：
                  <pre style="margin-left: 18px; background-color: #333; color: white; padding: 5px 10px; border-radius: 3px;">crictl pull my-registry.com/registry-name/image-name:tag</pre>
                </li>
                <li>
                  请将上面例子中的 my-registry.com 替换成你镜像仓库的域名
                </li>
                <li>
                  请将上面例子中的 192.168.30.56:5000 替换成你镜像仓库的地址及端口
                </li>
                <el-input style="margin-top: 10px;" v-model.trim="insecure_registry_form.domain_name" :placeholder="$t('domain_name')"></el-input>
              </div>
            </el-form-item>
            <div style="text-align: right; margin-top: 10px;">
              <el-button icon="el-icon-close" @click="visible = false">{{ $t('msg.cancel') }}</el-button>
              <el-button type="primary" icon="el-icon-check" @click="add">{{ $t('msg.ok') }}</el-button>
            </div>
          </el-form>
        </div>
      </el-popover>
      <KuboardSprayLink href="https://kuboard-spray.cn/guide/options/insecure-registry.html#containerd" target="_blank" style="margin-left: 20px;">帮 助</KuboardSprayLink>
    </template>
    <template #view>
      <template v-for="(item, key) in vars.containerd_insecure_registries" :key="'cir' + key">
        <div class="insecure_registry" style="padding-left: 10px;">
          <div effect="dark" class="app_text_mono header" style="margin-left: -10px;">
            {{key}}
          </div>
          <div v-for="(item, index) in vars.containerd_insecure_registries[key]" :key="key + index" style="padding-bottom: 5px;">
            <el-tag>{{ item }}</el-tag>
          </div>
        </div>
      </template>
      <div style="height: 10px; margin-top: 0px;">
        <KuboardSprayLink href="https://kuboard-spray.cn/guide/options/insecure-registry.html#containerd" target="_blank" style="font-size: 12px;">帮 助</KuboardSprayLink>
      </div>
    </template>
  </FieldCommon>
</template>

<script>
export default {
  props: {
    cluster: { type: Object, required: true },
  },
  data() {
    return {
      visible: false,
      insecure_registry_form: {
        domain_name: undefined,
      },
      insecureRegistriesItemRules: [
        { required: true, type: 'string', message: 'cannot be empty', trigger: 'blur' },
      ],
      domainNameRules: [
        { required: true, message: this.$t('domain_name_required'), trigger: 'blur' }
      ]
    }
  },
  computed: {
    vars: {
      get () { return this.cluster.inventory.all.children.target.vars },
      set () {},
    },
    prop() { return 'all.children.target.vars' }
  },
  components: { },
  mounted () {
  },
  methods: {
    add () {
      this.$refs.form.validate(flag => {
        if (flag) {
          this.vars.containerd_insecure_registries = this.vars.containerd_insecure_registries || {}
          this.vars.containerd_insecure_registries[this.insecure_registry_form.domain_name] = [undefined]
          this.visible = false
        }
      })
    }
  }
}
</script>

<style scoped lang="scss">
.insecure_registry {
  background-color: var(--el-color-primary-light-9);
  margin-bottom: 10px;
  padding: 0 10px 0px 0;
  .header {
    background-color: var(--el-color-primary-light-3);
    color: white;
    padding: 0 10px;
    margin-bottom: 10px;
    margin-right: -10px;
  }
}
.description {
  color: var(--el-color-text-secondary);
  line-height: 24px;
  padding: 10px;
  background-color: var(--el-color-info-light-9);
}
</style>
