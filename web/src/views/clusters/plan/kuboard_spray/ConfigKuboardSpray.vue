<i18n>
en:
  proxy: Proxy
  proxyDescription: Proxy (Generally speaking, you don't need this.)
  createResource: "Add Resource Package"
  proxyUsage1: KuboardSpray can use the following proxy params to fetch content from internet that are not included in the resource package.
  proxyUsage2: "Usually speaking, all the resources (exclude docker-ce / OS mirror) are already included in the resource package, you don't need the proxy params here."
  proxyUsage3: "Proxy params here only apply to the KuboardSpray host."
  bastionUsage: KuboardSpray can access Kubernetes Cluster Nodes through bastion. 
  setSshParam: Bastion is not enabled, please set SSH params in {tabName} tab.
zh:
  proxy: 代理
  proxyDescription: 代理（通常不需要设置）
  createResource: '添加资源包'
  proxyUsage1: KuboardSpray 可以使用如下代理参数从外网获取资源包中未找到的资源；
  proxyUsage2: 通常资源包中包含所需资源，您无需设置此处的代理参数；
  proxyUsage3: 此处的代理参数只作用于 KuboardSpray 这台机器；
  bastionUsage: KuboardSpray 可以通过堡垒机访问集群的节点。
  setSshParam: 未使用堡垒机时，请在 {tabName} 标签页设置 SSH 连接参数。
</i18n>

<template>
  <div>
    <ConfigSection v-model:enabled="useResourcePackage" disabled
      :label="$t('obj.resource')"
      :description="$t('obj.resource') + ' ' + (inventory.all.hosts.localhost.kuboardspray_resource_package ? inventory.all.hosts.localhost.kuboardspray_resource_package : '')">
      <FieldSelect :holder="inventory.all.hosts.localhost" fieldName="kuboardspray_resource_package" :loadOptions="loadResourceList" prop="all.hosts.localhost" required :disabled="isInstalled">
        <el-button style="margin-left: 10px;" type="primary" icon="el-icon-plus">{{$t('createResource')}}</el-button>
      </FieldSelect>
      <div v-if="resourcePackage">
        <ResourceDetails :resourcePackage="resourcePackage"></ResourceDetails>
      </div>
    </ConfigSection>
    <ConfigSection v-model:enabled="proxyEnabled" :label="$t('proxy')" :description="$t('proxyDescription')">
      <el-alert class="app_margin_bottom" :closable="false">
        <li>{{$t('proxyUsage1')}}</li>
        <li>{{$t('proxyUsage2')}}</li>
        <li>{{$t('proxyUsage3')}}</li>
      </el-alert>
      <FieldString :holder="inventory.all.hosts.localhost" fieldName="http_proxy" prop="all.hosts.localhost" required></FieldString>
      <FieldString :holder="inventory.all.hosts.localhost" fieldName="https_proxy" prop="all.hosts.localhost" required></FieldString>
      <FieldString :holder="inventory.all.hosts.localhost" fieldName="no_proxy" prop="all.hosts.localhost" required></FieldString>
    </ConfigSection>
    <ConfigSection v-model:enabled="bastionEnabled" :label="$t('obj.bastion')">
      <el-alert class="app_margin_bottom" :closable="false">{{$t('bastionUsage')}}</el-alert>
      <FieldString :holder="inventory.all.hosts.bastion" fieldName="ansible_host" prop="all.hosts.bastion" required></FieldString>
      <FieldString :holder="inventory.all.hosts.bastion" fieldName="ansible_user" prop="all.hosts.bastion" required></FieldString>
    </ConfigSection>
    <div v-if="!bastionEnabled" style="margin-left: 80px;">
      <el-alert :closable="false" type="warning">{{$t('setSshParam', {tabName: $t('node.k8s_cluster')})}}</el-alert>
    </div>
  </div>
</template>

<script>
import ResourceDetails from '../../../resources/details/ResourceDetails.vue'

export default {
  props: {
    cluster: { type: Object, required: true },
  },
  data () {
    return {
      localhostRules: [
        {
          validator: (rule, value, callback) => {
            callback('ErrorMessage')
          },
          trigger: 'blur'
        }
      ],
      useResourcePackage: true,
    }
  },
  inject: ['isInstalled'],
  computed: {
    inventory: {
      get () {
        return this.cluster.inventory
      },
      set () {}
    },
    resourcePackage: {
      get () { return this.cluster.resourcePackage},
      set () {}
    },
    proxyEnabled: {
      get () {
        if (this.inventory.all.hosts.localhost.vars) {
          return this.inventory.all.hosts.localhost.vars.http_proxy !== undefined
        }
        return false
      },
      set (v) {
        if (v) {
          this.inventory.all.hosts.localhost.vars = this.inventory.all.hosts.localhost.vars || {}
          this.inventory.all.hosts.localhost.vars.http_proxy = ''
        } else {
          delete this.inventory.all.hosts.localhost.vars.http_proxy
        }
      }
    },
    bastionEnabled: {
      get () {
        return this.inventory.all.children.target.children.bastion !== undefined
      },
      set (v) {
        if (v) {
          this.inventory.all.hosts.bastion = this.inventory.all.hosts.bastion || {ansible_host: '', ansible_user: ''}
          this.inventory.all.children.target.children.bastion = {hosts: {bastion: {}}}
        } else {
          delete this.inventory.all.children.target.children.bastion
        }
      }
    }
  },
  components: { ResourceDetails },
  mounted () {
  },
  methods: {
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
    }
  }
}
</script>

<style scoped lang="scss">

</style>
