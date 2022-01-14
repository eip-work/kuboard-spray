<i18n>
en:
  createResource: "Add Resource Package"
  bastionUsage: KuboardSpray can access Kubernetes Cluster Nodes through bastion. 
  setSshParam: Bastion is not enabled, please set SSH params in {tabName} tab.
  goToResourcePage: It's about to open Resource Package Download page in new window, do you confirm?
zh:
  createResource: '添加资源包'
  bastionUsage: KuboardSpray 可以通过跳板机或堡垒机访问将要安装 K8S 集群的目标节点。
  setSshParam: 未使用跳板机或堡垒机时，请在 {tabName} 标签页设置 SSH 连接参数。
  goToResourcePage: '此操作将在新窗口打开资源包导入页面，完成设置后，您可以切换窗口会到当前页面，是否继续？'
</i18n>

<template>
  <div>
    <ConfigSection v-model:enabled="useResourcePackage" disabled anti-freeze
      :label="$t('obj.resource')"
      :description="$t('obj.resource') + ' ' + (inventory.all.hosts.localhost.kuboardspray_resource_package ? inventory.all.hosts.localhost.kuboardspray_resource_package : '')">
      <FieldSelect :holder="inventory.all.hosts.localhost" fieldName="kuboardspray_resource_package" :loadOptions="loadResourceList" prop="all.hosts.localhost" required :disabled="isClusterInstalled">
        <template #edit>
          <ConfirmButton buttonStyle="margin-left: 10px;" icon="el-icon-plus" 
            @confirm="openUrlInBlank('/#/settings/resources')"
            :text="$t('createResource')" :message="$t('goToResourcePage')"></ConfirmButton>
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
  inject: ['isClusterInstalled', 'isClusterOnline'],
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
          delete this.inventory.all.hosts.bastion
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
