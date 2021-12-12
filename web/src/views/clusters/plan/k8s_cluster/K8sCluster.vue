<i18n>
en:
  label: Common
  description: Most frequently referenced params to Kubernetes.
  container_manager: container_manager
zh:
  label: 常用参数
  description: Kubernetes 集群的常用参数
  container_manager: 容器引擎
</i18n>

<template>
  <ConfigSection v-model:enabled="enabled" :label="$t('label')" :description="$t('description')" disabled>
    <template v-if="resourcePackage !== undefined">
      <el-form-item label="kube_version">
        <span class="app_text_mono">{{resourcePackage.kubernetes.kube_version}}</span>
      </el-form-item>
      <el-form-item :label="$t('container_manager')">
        <span class="app_text_mono">{{resourcePackage.docker_engine.container_manager}}_{{resourcePackage.docker_engine.version}}</span>
      </el-form-item>
    </template>
    <FieldString :holder="inventory.all.children.k8s_cluster.vars" fieldName="cluster_name"></FieldString>
    <FieldString :holder="inventory.all.children.k8s_cluster.vars" fieldName="event_ttl_duration"></FieldString>
    <!-- <FieldBool :holder="inventory.all.children.k8s_cluster.vars" fieldName="auto_renew_certificates"></FieldBool> -->
  </ConfigSection>
</template>

<script>
import ConfigSection from '../ConfigSection.vue'

export default {
  props: {
    inventory: { type: Object, required: true },
    resourcePackage: { type: Object, required: false, default: undefined },
  },
  data() {
    return {

    }
  },
  computed: {
    enabled: {
      get () {return true},
      set () {},
    },
  },
  components: { ConfigSection },
  mounted () {
  },
  methods: {

  }
}
</script>

<style scoped lang="scss">

</style>
