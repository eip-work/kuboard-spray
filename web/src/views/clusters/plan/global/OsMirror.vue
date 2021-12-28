<i18n>
en:
  label: OS Mirror
  description: Operation System source repo.
zh:
  label: OS软件源
  description: OS 软件（指定操作系统软件源）
</i18n>

<template>
  <ConfigSection v-model:enabled="enabled" :label="$t('label')" :description="$t('description')" disabled anti-freeze>
    <el-alert title="离线安装请注意" type="warning" :closable="false" class="app_margin_bottom">
      <li>安装 Kubernetes 集群时，需要在 K8S 节点上安装一些系统软件，例如：ipset / ethtool / docker-ce 等；</li>
      <li>大部分企业对于自己常用的操作系统都有离线环境下的软件源，在此情况下，为了减小资源包的尺寸，KuboardSpray 并不将这些依赖的系统软件打包到资源包里；</li>
      <li>如果您的所有 K8S 节点都可以访问到下面的默认软件源，您可以忽略此处的参数。</li>
    </el-alert>
    <FieldString :holder="cluster.inventory.all.children.target.vars" fieldName="docker_ubuntu_repo_base_url"></FieldString>
    <FieldString :holder="cluster.inventory.all.children.target.vars" fieldName="docker_ubuntu_repo_gpgkey"></FieldString>
    <FieldString :holder="cluster.inventory.all.children.target.vars" fieldName="docker_ubuntu_repo_repokey"></FieldString>
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
    enabled: {
      get () {return true},
      set () {},
    },
  },
  components: { },
  mounted () {
  },
  methods: {

  }
}
</script>

<style scoped lang="scss">

</style>
