<i18n>
en:
  proxy: Proxy
  proxyDescription: Proxy (Generally speaking, you don't need this.)
  proxyUsage1: KuboardSpray can use the following proxy params to fetch content from internet that are not included in the resource package.
  proxyUsage2: "Usually speaking, all the resources (exclude docker-ce / OS mirror) are already included in the resource package, you don't need the proxy params here."
zh:
  proxy: 代理
  proxyDescription: 代理（通常不需要设置）
  proxyUsage1: KuboardSpray 可以使用如下代理参数从外网获取资源包中未找到的资源；
  proxyUsage2: 通常资源包中包含所需资源，您无需设置此处的代理参数；
</i18n>

<template>
  <ConfigSection v-model:enabled="proxyEnabled" :label="$t('proxy')" :description="$t('proxyDescription')" anti-freeze>
    <el-alert class="app_margin_bottom" :closable="false">
      <li>{{$t('proxyUsage1')}}</li>
      <li>{{$t('proxyUsage2')}}</li>
    </el-alert>
    <FieldString :holder="vars" fieldName="http_proxy" prop="all.children.target.vars" anti-freeze></FieldString>
    <FieldString :holder="vars" fieldName="https_proxy" prop="all.children.target.vars" anti-freeze></FieldString>
    <FieldString :holder="vars" fieldName="additional_no_proxy" prop="all.children.target.vars" anti-freeze></FieldString>
    <FieldCommon :holder="cluster.inventory.all.vars" fieldName="container_manager_on_localhost" prop="all.children.target.vars" anti-freeze>
      <template #edit>
        <el-select v-model="container_manager_on_localhost" style="width: 100%;">
          <el-option label="docker" value="docker">docker</el-option>
          <el-option label="containerd" value="containerd">containerd</el-option>
        </el-select>
      </template>
      <template #view>
        {{ container_manager_on_localhost }}
      </template>
    </FieldCommon>
    <FieldSelect v-if="cluster.inventory.all.vars.container_manager_on_localhost === 'docker'"
      :holder="cluster.inventory.all.vars" fieldName="docker_bin_dir" prop="all.children.target.vars" :loadOptions="loadDockerDirOptions"></FieldSelect>
  </ConfigSection>
</template>

<script>
export default {
  props: {
    cluster: { type: Object, required: true },
  },
  data() {
    return {

    }
  },
  computed: {
    inventory: {
      get () { return this.cluster.inventory },
      set () {},
    },
    container_manager_on_localhost: {
      get () { return this.inventory.all.vars.container_manager_on_localhost },
      set (v) {
        if (v === 'docker') {
          this.inventory.all.vars.image_command_tool_on_localhost = 'docker'
          this.inventory.all.vars.docker_bin_dir = '/Applications/Docker.app/Contents/Resources/bin'
        } else {
          this.inventory.all.vars.image_command_tool_on_localhost = 'nerdctl'
        }
        this.inventory.all.vars.container_manager_on_localhost = v
      }
    },
    vars: {
      get () {
        return this.cluster.inventory.all.children.target.vars
      },
      set () {}
    },
    proxyEnabled: {
      get () {
        if (this.vars) {
          return this.vars.http_proxy !== undefined
        }
        return false
      },
      set (v) {
        if (v) {
          this.vars = this.vars || {}
          this.vars.http_proxy = ''
        } else {
          delete this.vars.http_proxy
          delete this.vars.https_proxy
          delete this.vars.additional_no_proxy
        }
      }
    },
  },
  components: { },
  mounted () {
  },
  methods: {
    async loadDockerDirOptions () {
      return [
        { label: '/usr/bin', value: '/usr/bin' },
        { label: '/usr/local/bin', value: '/usr/local/bin' },
        { label: '/Applications/Docker.app/Contents/Resources/bin', value: '/Applications/Docker.app/Contents/Resources/bin' },
      ]
    }
  }
}
</script>

<style scoped lang="css">

</style>
