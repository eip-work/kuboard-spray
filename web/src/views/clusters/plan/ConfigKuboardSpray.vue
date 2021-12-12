<i18n>
en:
  proxy: Proxy
  createResource: "Add Resource Package"
  included: Included
  not_included: Not included
  addons: Addons
  cni: Cni plugins
  dependency: Dependencies
  package_content: Package Content
  proxyUsage1: KuboardSpray can use the following proxy params to fetch content from internet that are not included in the resource package.
  proxyUsage2: "Usually speaking, all the resources (exclude docker-ce / OS mirror) are already included in the resource package, you don't need the proxy params here."
  proxyUsage3: "Proxy params here only apply to the KuboardSpray host."
  bastionUsage: KuboardSpray can access Kubernetes Cluster Nodes through bastion. 
  setSshParam: Bastion is not enabled, please set SSH params in {tabName} tab.
zh:
  proxy: 代理
  createResource: '添加资源包'
  included: 已包含
  not_included: 未包含
  addons: 可选组件
  cni: 网络插件
  dependency: 依赖组件
  package_content: 资源包内容
  proxyUsage1: KuboardSpray 可以使用如下代理参数从外网获取资源包中未找到的资源；
  proxyUsage2: 通常资源包中包含所需资源，您无需设置此处的代理参数；
  proxyUsage3: 此处的代理参数只作用于 KuboardSpray 这台机器；
  bastionUsage: KuboardSpray 可以通过堡垒机访问集群的节点。
  setSshParam: 未使用堡垒机时，请在 {tabName} 标签页设置 SSH 连接参数。
</i18n>

<template>
  <div>
    <ConfigSection v-model:enabled="useResourcePackage" disabled
      :label="$t('obj.resources')"
      :description="$t('obj.resources') + ' ' + (inventory.all.hosts.localhost.kuboardspray_resource_package ? inventory.all.hosts.localhost.kuboardspray_resource_package : '')">
      <FieldSelect :holder="inventory.all.hosts.localhost" fieldName="kuboardspray_resource_package" :loadOptions="loadResourceList" prop="all.hosts.localhost" required>
        <el-button style="margin-left: 10px;" type="primary" icon="el-icon-plus">{{$t('createResource')}}</el-button>
      </FieldSelect>
      <div v-if="resourcePackage" class="app_form_mini">
        <div style="text-align: center; margin-bottom: 10px; margin-top: -10px; font-weight: bold;">[ {{$t('package_content')}} ]</div>
        <el-collapse v-model="activeNames">
          <el-collapse-item name="1">
            <template #title>
              <span class="package_title">kubespray</span>
            </template>
            <div class="package_info">
              <PackageContentField :holder="resourcePackage.kuboardspray" fieldName="kubespray_version"></PackageContentField>
            </div>
          </el-collapse-item>
          <el-collapse-item name="2">
            <template #title>
              <span class="package_title">kubernetes</span>
            </template>
            <div class="package_info">
              <PackageContentField :holder="resourcePackage.kubernetes" fieldName="image_arch"></PackageContentField>
              <PackageContentField :holder="resourcePackage.kubernetes" fieldName="gcr_image_repo"></PackageContentField>
              <PackageContentField :holder="resourcePackage.kubernetes" fieldName="kube_image_repo"></PackageContentField>
              <PackageContentField :holder="resourcePackage.kubernetes" fieldName="kube_version"></PackageContentField>
              <el-form-item label="container_manager" label-width="160px">
                <span class="app_text_mono">{{resourcePackage.docker_engine.container_manager}}_{{resourcePackage.docker_engine.version}}</span>
              </el-form-item>
            </div>
          </el-collapse-item>
          <el-collapse-item name="3">
            <template #title>
              <span class="package_title">etcd</span>
            </template>
            <div class="package_info">
              <PackageContentField :holder="resourcePackage.etcd" fieldName="etcd_version"></PackageContentField>
              <PackageContentField :holder="resourcePackage.etcd" fieldName="etcd_deployment_type"></PackageContentField>
            </div>
          </el-collapse-item>
          <el-collapse-item name="4">
            <template #title>
              <span class="package_title">{{$t('cni')}}</span>
            </template>
            <div class="package_info">
              <template v-for="(item, index) in resourcePackage.cni" :key="index + 'cni'">
                <el-form-item :label="item.target" label-width="160px">
                  <div class="app_text_mono">
                    {{item.version}}
                  </div>
                </el-form-item>
              </template>
            </div>
          </el-collapse-item>
          <el-collapse-item name="5">
            <template #title>
              <span class="package_title">{{$t('dependency')}}</span>
            </template>
            <div class="package_info">
              <template v-for="(item, index) in resourcePackage.dependency" :key="index + 'dependency'">
                <el-form-item :label="item.target" label-width="160px">
                  <div class="app_text_mono">
                    {{item.version}}
                  </div>
                </el-form-item>
              </template>
            </div>
          </el-collapse-item>
          <el-collapse-item name="6">
            <template #title>
              <span class="package_title">{{$t('addons')}}</span>
            </template>
            <div class="package_info">
              <template v-for="(item, index) in resourcePackage.addons" :key="index + 'addons'">
                <el-form-item label-width="160px">
                  <template #label>
                    <div style="font-weight: bolder;">{{item.name}}</div>
                  </template>
                  <div class="app_text_mono">
                    <el-tag type="success" v-if="item.included">{{$t('included')}}</el-tag>
                    <el-tag type="info" v-else>{{$t('not_included')}}</el-tag>
                  </div>
                </el-form-item>
                <div class="package_info">
                  <template v-for="(value, key) in item.params" :key="key + 'addons' + index">
                    <el-form-item :label="key" label-width="220px">
                      <template #label>
                        <div style="font-size: 12px">{{key}}</div>
                      </template>
                      <div class="app_text_mono" style="font-size: 12px">
                        {{value}}
                      </div>
                    </el-form-item>
                  </template>
                </div>
              </template>
            </div>
          </el-collapse-item>
        </el-collapse>
      </div>
    </ConfigSection>
    <ConfigSection v-model:enabled="proxyEnabled" :label="$t('proxy')">
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
import ConfigSection from './ConfigSection.vue'
import PackageContentField from './common/PackageContentField.vue'

export default {
  props: {
    inventory: { type: Object, required: true }
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
      resourcePackage: undefined,
      activeNames: ['2'],
    }
  },
  computed: {
    inventoryRef: {
      get () {
        return this.inventory
      },
      set () {}
    },
    proxyEnabled: {
      get () {
        if (this.inventoryRef.all.hosts.localhost.vars) {
          return this.inventoryRef.all.hosts.localhost.vars.http_proxy !== undefined
        }
        return false
      },
      set (v) {
        if (v) {
          this.inventoryRef.all.hosts.localhost.vars = this.inventoryRef.all.hosts.localhost.vars || {}
          this.inventoryRef.all.hosts.localhost.vars.http_proxy = ''
        } else {
          delete this.inventoryRef.all.hosts.localhost.vars.http_proxy
        }
      }
    },
    bastionEnabled: {
      get () {
        return this.inventoryRef.all.children.bastion !== undefined
      },
      set (v) {
        if (v) {
          this.inventoryRef.all.hosts.bastion = this.inventoryRef.all.hosts.bastion || {ansible_host: '', ansible_user: ''}
          this.inventoryRef.all.children.bastion = {hosts: {bastion: {}}}
        } else {
          delete this.inventoryRef.all.children.bastion
        }
      }
    }
  },
  components: { ConfigSection, PackageContentField },
  mounted () {
  },
  watch: {
    'inventory.all.hosts.localhost.kuboardspray_resource_package': function(newValue) {
      this.resourcePackage = undefined
      if (newValue) {
        this.kuboardSprayApi.get(`/resources/${newValue}`).then(resp => {
          this.resourcePackage = resp.data.data.package
          this.$emit('update:resourcePackage', this.resourcePackage)
        }).catch(e => {
          console.log(e)
        })
      }
    }
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
.package_title {
  font-weight: bolder;
}
.package_info {
  margin-left: 20px;
}
</style>
