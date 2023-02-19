<i18n>
en:
  add_insecure_registry: Add insecure registry
  insecure_registries_placeholder: e.g. http://192.168.30.56:5000
  domain_name: Domain name of the registry
  domain_name_required: Domain name is required
  address: Address of the registry
zh:
  add_insecure_registry: 添加 http 协议镜像仓库（或自签名 https 的镜像仓库）
  insecure_registries_placeholder: 例如：http://192.168.30.56:5000
  domain_name: 镜像仓库的域名
  domain_name_required: 镜像仓库的域名为必填字段
  address: 镜像仓库的地址
</i18n>

<template>
  <FieldCommon :holder="vars" :prop="prop" fieldName="containerd_insecure_registries" anti-freeze>
    <template #edit>
      <template v-for="(item, key) in vars.containerd_insecure_registries" :key="'cir' + key">
        <div class="insecure_registry">
          <div effect="dark" class="app_text_mono header" style="display: flex;">
            <div style="width: 100%;">{{key}}</div>
            <el-popconfirm :title="$t('msg.delete')" @confirm="delete vars.containerd_insecure_registries[key]" style="float: right;">
              <template #reference>
                <el-button type="primary" icon="el-icon-delete">{{ $t('msg.delete') }}</el-button>
              </template>
            </el-popconfirm>
          </div>
          <div style="padding: 0 0 0 10px; display: flex;" class="app_text_mono">
            <div style="width: 100%;">{{vars.containerd_insecure_registries[key]}}</div>
            <ContainerManagerRegistryEdit :cluster="cluster" :name="key"></ContainerManagerRegistryEdit>
          </div>
          <!-- <FieldArray :holder="vars.containerd_insecure_registries" newItemTemplate="" :fieldName="key" label-width="0px" label=" "
            :itemProp="`${prop}.containerd_insecure_registries['${key}']`" :itemRules="insecureRegistriesItemRules" anti-freeze>
            <template v-slot:editItem="scope">
              <el-input v-model.trim="item[scope.index]" :placeholder="$t('insecure_registries_placeholder')"></el-input>
            </template>
          </FieldArray> -->
        </div>
      </template>
      <ContainerManagerRegistryEdit :cluster="cluster"></ContainerManagerRegistryEdit>
      <KuboardSprayLink href="https://kuboard-spray.cn/guide/options/insecure-registry.html#containerd" target="_blank" style="margin-left: 20px;">帮 助</KuboardSprayLink>
    </template>
    <template #view>
      <template v-for="(item, key) in vars.containerd_insecure_registries" :key="'cir' + key">
        <div class="insecure_registry" style="padding-left: 10px;">
          <div effect="dark" class="app_text_mono header" style="margin-left: -10px;">
            {{key}}
          </div>
          <div style="padding-bottom: 5px;">
            {{ vars.containerd_insecure_registries[key] }}
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
import ContainerManagerRegistryEdit from './ContainerManagerRegistryEdit.vue'

export default {
  props: {
    cluster: { type: Object, required: true },
  },
  data() {
    return {
      insecureRegistriesItemRules: [
        { required: true, type: 'string', message: '字段不能为空', trigger: 'blur' },
      ],
    }
  },
  computed: {
    vars: {
      get () { return this.cluster.inventory.all.children.target.vars },
      set () {},
    },
    prop() { return 'all.children.target.vars' }
  },
  components: { ContainerManagerRegistryEdit },
  mounted () {
  },
  methods: {
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
