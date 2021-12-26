<i18n>
en:
  title: Resource Packages
  resourceDescription1: Kuboard provides a set of pre-verified resource packages.
  resourceDescription2: You can also refer to https://github.com/eip-work/kuboard-spray-resource to create your own resource package.
  kuboard_spray_version_min: KuboardSpray min version
  kuboard_spray_version_max: KuboardSpray max version
  version: Version
  container_engine: Container Engine
  supported_os: Supported OS
  kube_version: K8S Version
  kubespray_version: Kubespray Version
zh:
  title: 资源包列表
  resourceDescription1: Kuboard 提供一组经过预先测试验证的资源包列表，可以帮助您快速完成集群安装
  resourceDescription2: 您也可以参考项目 https://github.com/eip-work/kuboard-spray-resource 自己创建资源包
  kuboard_spray_version_min: KuboardSpray最低版本
  kuboard_spray_version_max: KuboardSpray最高版本
  version: 资源包版本
  container_engine: 容器引擎
  supported_os: 支持的操作系统
  kube_version: K8S 版本
  kubespray_version: Kubespray 版本
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
      <el-table v-if="packageList" :data="packageList.items" height="250" style="width: 100%">
        <el-table-column prop="version" :label="$t('version')" min-width="100px"/>
        <!-- <el-table-column prop="kuboard_spray_version.min" :label="$t('kuboard_spray_version_min')" align="center"/>
        <el-table-column prop="kuboard_spray_version.max" :label="$t('kuboard_spray_version_max')" align="center">
          <template #default="scope">
            {{ scope.row.kuboard_spray_version.max || '-' }}
          </template>
        </el-table-column> -->
        <el-table-column :label="$t('kubespray_version')">
          <template #default="scope">
            <span v-if="packageMap[scope.row.version]">
              {{ packageMap[scope.row.version].metadata.kubespray_version }}
            </span>
            <i class="el-icon-loading" v-else></i>
          </template>
        </el-table-column>
        <el-table-column :label="$t('kube_version')">
          <template #default="scope">
            <span v-if="packageMap[scope.row.version]">
              {{ packageMap[scope.row.version].kubernetes.kube_version }}
            </span>
            <i class="el-icon-loading" v-else></i>
          </template>
        </el-table-column>
        <el-table-column :label="$t('container_engine')">
          <template #default="scope">
            <template v-if="packageMap[scope.row.version]">
              <div v-for="(engine, key) in packageMap[scope.row.version].container_engine" :key="`c${scope.index}_${key}`">
                <el-tag>{{ engine.container_manager }} {{ engine.version }}</el-tag>
              </div>
            </template>
            <i class="el-icon-loading" v-else></i>
          </template>
        </el-table-column>
        <el-table-column :label="$t('supported_os')">
          <template #default="scope">
            <template v-if="packageMap[scope.row.version]">
              <div v-for="(os, key) in packageMap[scope.row.version].supported_os" :key="`os${scope.index}_${key}`">
                <el-tag>{{ os.distribution }} {{ os.version }}</el-tag>
              </div>
            </template>
            <i class="el-icon-loading" v-else></i>
          </template>
        </el-table-column>
      </el-table>
    </div>

  </div>
</template>

<script>
import mixin from '../../mixins/mixin.js'
import axios from 'axios'
import yaml from 'js-yaml'

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
  refresh () {
    this.refresh()
  },
  data () {
    return {
      packageList: undefined,
      packageMap: {},
    }
  },
  computed: {
  },
  components: { },
  mounted () {
    this.refresh()
  },
  methods: {

    refresh () {
      axios.get('https://addons.kuboard.cn/v-kuboard-spray-resources/package-list.yaml').then(resp => {
        this.packageList = yaml.load(resp.data)
        for (let i in this.packageList.items) {
          let packageVersion = this.packageList.items[i]
          this.loadPackage(packageVersion)
        }
      }).catch(e => {
        console.log(e)
        this.$message.error('离线环境')
      })
    },
    loadPackage (packageVersion) {
      axios.get(`https://addons.kuboard.cn/v-kuboard-spray-resources/${packageVersion.version}/package.yaml`).then(resp => {
        setTimeout(() => {
          this.packageMap[packageVersion.version] = yaml.load(resp.data)
        }, 500)
      }).catch(e => {
        this.$message.error(e + '')
      })
    }
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
