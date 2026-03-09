<i18n>
en:
  createResource: "Add Resource Package"
  goToResourcePage: It's about to open Resource Package Download page in new window, do you confirm?
  pangeecluster_resource_package: 'Resource Package'
zh:
  createResource: '添加资源包'
  goToResourcePage: '此操作将在新窗口打开资源包导入页面，完成设置后，您可以切换窗口会到当前页面，是否继续？'
  pangeecluster_resource_package: '资源包'
  pangeecluster_resource_package_placeholder: '请选择资源包'
</i18n>

<template>
  <div>
    <ConfigSection v-model:enabled="useResourcePackage" disabled anti-freeze :label="$t('obj.resource')" :description="$t('obj.resource') +
      ' ' +
      (inventory.all.hosts.localhost.pangeecluster_resource_package
        ? inventory.all.hosts.localhost.pangeecluster_resource_package
        : '')
      ">
      <EditSelect v-model="inventory.all.hosts.localhost.pangeecluster_resource_package"
        :label="t('pangeecluster_resource_package')" :placeholder="t('pangeecluster_resource_package_placeholder')"
        :loadOptions="loadResourceList" prop="all.hosts.localhost" required :disabled="isClusterInstalled">
        <template #edit>
          <ConfirmButton buttonStyle="margin-left: 10px;" icon="el-icon-plus"
            @confirm="openUrlInBlank('/#/settings/resources')" :text="t('createResource')"
            :message="t('goToResourcePage')"></ConfirmButton>
        </template>
      </EditSelect>
      <div v-if="resourcePackage">
        <ResourceDetails :resourcePackage="resourcePackage"></ResourceDetails>
      </div>
    </ConfigSection>
  </div>
</template>

<script>
import ResourceDetails from "../../../resources/details/ResourceDetails.vue";
import { computed } from 'vue'

export default {
  props: {
    cluster: { type: Object, required: true }
  },
  data() {
    return {
      localhostRules: [
        {
          validator: (rule, value, callback) => {
            callback("ErrorMessage");
          },
          trigger: "blur"
        }
      ],
      useResourcePackage: true
    };
  },
  provide() {
    return {
      editMode: computed(() => {
        return "view";
      })
    };
  },
  inject: ["isClusterInstalled", "isClusterOnline"],
  computed: {
    inventory: {
      get() {
        return this.cluster.inventory;
      },
      set() { }
    },
    resourcePackage: {
      get() {
        return this.cluster.resourcePackage;
      },
      set() { }
    }
  },
  components: { ResourceDetails },
  mounted() { },
  methods: {
    async loadResourceList() {
      let result = [];
      await this.pangeeClusterApi
        .get("/resources")
        .then(resp => {
          for (let res of resp.data.data) {
            result.push({ label: res, value: res });
          }
        })
        .catch(e => {
          console.log(e);
        });
      return result;
    }
  }
};
</script>

<style scoped lang="css"></style>
