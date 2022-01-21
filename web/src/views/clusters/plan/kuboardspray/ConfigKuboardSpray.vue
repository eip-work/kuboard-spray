<i18n>
en:
  createResource: "Add Resource Package"
  goToResourcePage: It's about to open Resource Package Download page in new window, do you confirm?
zh:
  createResource: '添加资源包'
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
    <SshParamsBastion v-if="cluster && cluster.inventory" :cluster="cluster" nodeName="bastion" :holder="inventory.all.hosts.bastion || {}" 
      prop="all.hosts.bastion"></SshParamsBastion>
  </div>
</template>

<script>
import ResourceDetails from '../../../resources/details/ResourceDetails.vue'
import SshParamsBastion from './SshParamsBastion.vue'

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
  },
  components: { ResourceDetails, SshParamsBastion },
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
