<i18n>
en:
  clusters: Kubernetes Cluster
  resources: Resources Package
  kuboard: Kuboard Integration
zh:
  clusters: 集群管理
  resources: 资源包管理
  kuboard: Kuboard 集成
</i18n>

<template>
  <div>
    <div class="app_block_title">Kuboard Spray</div>
    <el-alert type="info" title="README" :closable="false">
      <div style="line-height: 28px;">
        <li>Kuboard Spray 基于 <a href="https://github.com/kubernetes-sigs/kubespray" target="_blank">https://github.com/kubernetes-sigs/kubespray</a> 提供 Kubernetes 集群的离线安装、升级、维护等功能。</li>
        <li>您可以在此找到 Kuboard Spray 的源代码 <a href="https://github.com/eip-work/kuboard-spray" target="_blank">https://github.com/eip-work/kuboard-spray</a></li>
        <li>如果您发现 Kuboard Spray 的 BUG，请及时提出 <a href="https://github.com/eip-work/kuboard-spray/issues" target="_blank">https://github.com/eip-work/kuboard-spray/issues</a></li>
      </div>
    </el-alert>
    <div style="display: flex; margin-top: 20px;">
      <el-card class="card">
        <template #header>
          <el-icon :size="20" style="margin-right: 10px;"><el-icon-box /></el-icon>
          {{$t('resources')}}
        </template>
        <div class="description">
          下载不同版本 Kubernetes 集群所对应的资源包，以实现离线安装 Kubernetes 集群。
        </div>
        <div class="buttons">
          <el-button type="primary" icon="el-icon-position" @click="$router.push(`/settings/resources`)">{{$t('resources')}}</el-button>
        </div>
      </el-card>
      <el-card style="margin: 0 20px;" class="card">
        <template #header>
          <el-icon :size="20" style="margin-right: 10px;"><img src="./images/kubernetes.png" style="width: 20px; height: 20px;"></el-icon>
          {{$t('clusters')}}
        </template>
        <div class="description">
          通过向导的方式规划和安装您的 Kubernetes 集群，并对集群进行后续的维护管理。
        </div>
        <div class="buttons">
          <el-button type="primary" icon="el-icon-position" @click="$router.push(`/clusters`)">
            {{$t('clusters')}}
          </el-button>
        </div>
      </el-card>
      <el-card class="card">
        <template #header>
          <el-icon :size="20" style="margin-right: 10px;">
            <img style="width: 20px; height: 20px;" :src="kuboardIcon">
          </el-icon>
          {{$t('kuboard')}}
        </template>
        <div class="description">
          将 Kuboard Spray 集成到 Kuboard 界面中，更方便地管理您的 Kubernetes 集群。
        </div>
        <div class="buttons">
          <el-button type="primary" icon="el-icon-position" @click="$router.push(`/settings/kuboard`)">{{$t('kuboard')}}</el-button>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script>
import mixin from '../../mixins/mixin.js'
import { kuboardIcon } from './images/icon.js'

export default {
  mixins: [mixin],
  props: {
  },
  percentage () {
    return 100
  },
  breadcrumb () {
    return []
  },
  data() {
    return {
      kuboardIcon,
    }
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

<style scoped lang="css">
.description {
  font-size: 13px;
  line-height: 20px;
}
.buttons {
  margin-top: 20px;
  text-align: right;
}
.buttons button {
  width: 150px;
}
.card {
  width: 33%;
  flex-grow: 1;
}
</style>
