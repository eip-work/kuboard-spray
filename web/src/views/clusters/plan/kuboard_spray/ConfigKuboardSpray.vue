<i18n>
en:
  createResource: "Add Resource Package"
  bastionUsage: KuboardSpray can access Kubernetes Cluster Nodes through bastion. 
  setSshParam: Bastion is not enabled, please set SSH params in {tabName} tab.
zh:
  createResource: '添加资源包'
  bastionUsage: KuboardSpray 可以通过跳板机或堡垒机访问将要安装 K8S 集群的目标节点。
  setSshParam: 未使用跳板机或堡垒机时，请在 {tabName} 标签页设置 SSH 连接参数。
</i18n>

<template>
  <div>
    <ConfigSection v-model:enabled="useResourcePackage" disabled anti-freeze
      :label="$t('obj.resource')"
      :description="$t('obj.resource') + ' ' + (inventory.all.hosts.localhost.kuboardspray_resource_package ? inventory.all.hosts.localhost.kuboardspray_resource_package : '')">
      <FieldSelect :holder="inventory.all.hosts.localhost" fieldName="kuboardspray_resource_package" :loadOptions="loadResourceList" prop="all.hosts.localhost" required :disabled="isInstalled">
        <template #edit>
          <el-button style="margin-left: 10px;" type="primary" icon="el-icon-plus">{{$t('createResource')}}</el-button>
        </template>
      </FieldSelect>
      <div v-if="resourcePackage">
        <ResourceDetails :resourcePackage="resourcePackage"></ResourceDetails>
      </div>
    </ConfigSection>
    <ConfigSection v-model:enabled="bastionEnabled" :label="$t('obj.bastion')" :description="$t('bastionUsage')" anti-freeze>
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
