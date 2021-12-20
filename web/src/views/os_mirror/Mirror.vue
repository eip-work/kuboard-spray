<i18n>
en:
  titleRepo: OS mirrors
  repoDescription1: You muse provide a mirror address (e.g. CentOS yum mirror, Ubuntu apt mirror) that all the nodes in your mirror could access.
  repoDescription2: Usually, a enterprise has its own os mirror, if not, Kuboard alos provide a wizard, so that you can create an os mirror quickly.
zh:
  titleRepo: OS 镜像源
  repoDescription1: 您必须定义您集群机器可以访问的操作系统镜像地址（例如 CentOS 的 yum 源、Ubuntu 的 apt 源等）以供使用
  repoDescription2: 通常企业都有自己的常用操作系统的本地镜像仓库，Kuboard 提供向导，帮助您快速设置一个操作系统镜像源
</i18n>

<template>
  <div>
    hello
  </div>
</template>

<script>
import mixin from '../../mixins/mixin.js'
export default {
  mixins: [mixin],
  props: {
    name: { type: String, required: true },
  },
  data() {
    return {
      percentage: 0,
    }
  },
  refresh () {
    this.refresh()
  },
  percentage () {
    return 100
  },
  breadcrumb () {
    return [
      { label: this.$t('titleRepo'), to: '/settings/mirrors' },
      { label: this.name },
    ]
  },
  computed: {
  },
  components: { },
  mounted () {
    this.refresh()
  },
  methods: {
    refresh () {

    }
  }
}
</script>

<style scoped lang="scss">

</style>
