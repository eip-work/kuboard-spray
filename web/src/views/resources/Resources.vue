<i18n>
en:
  title: Resource Packages
  resourceDescription1: Kuboard provides a set of pre-verified resource packages.
  resourceDescription2: You can also refer to https://github.com/eip-work/kuboard-spray-resource to create your own resource package.
zh:
  title: 资源包列表
  resourceDescription1: Kuboard 提供一组经过预先测试验证的资源包列表，可以帮助您快速完成集群安装
  resourceDescription2: 您也可以参考项目 https://github.com/eip-work/kuboard-spray-resource 自己创建资源包
  titleRepo: 操作系统镜像
  repoDescription1: 您必须定义您集群机器可以访问的操作系统镜像地址（例如 CentOS 的 yum 源、Ubuntu 的 apt 源等）以供使用
  repoDescription2: 通常企业都有自己的常用操作系统的本地镜像仓库，Kuboard 提供向导，帮助您快速设置一个操作系统镜像源
</i18n>

<template>
  <div>
    <div class="app_block_title">{{$t('title')}}</div>
    <el-alert :title="$t('title')" type="default" :closable="false">
      <div class="description">
        <li>{{$t('resourceDescription1')}}</li>
        <li>{{$t('resourceDescription2')}}</li>
      </div>
    </el-alert>
    <div class="contentList">
      
    </div>

    <div class="app_block_title">{{$t('titleRepo')}}</div>
    <el-alert :title="$t('titleRepo')" type="default" :closable="false">
      <div class="description">
        <li>{{$t('repoDescription1')}}</li>
        <li>{{$t('repoDescription2')}}</li>
      </div>
    </el-alert>

    <div class="contentList">

    </div>

  </div>
</template>

<script>
import mixin from '../../mixins/mixin.js'

export default {
  mixins: [mixin],
  props: {
  },
  percentage () {
    return 100
  },
  breadcrumb () {
    return [
      { label: '资源包管理' }
    ]
  },
  computed: {
  },
  components: { },
  mounted () {
  },
  methods: {

  }
}
</script>

<style scoped lang="scss">
.description {
  line-height: 28px;
}
.contentList {
  margin: 10px 0;
}
</style>
