<i18n>
en:
  title: Resource Packages
  resourceDescription1: Kuboard provides a set of pre-verified resource packages.
  resourceDescription2: You can also refer to https://github.com/eip-work/kuboard-spray-resource to create your own resource package.
  kuboardspray_version_min: KuboardSpray min version
  kuboardspray_version_max: KuboardSpray max version
  version: Version
  container_engine: Supported Container Engine
  supported_os: Supported OS
  kube_version: K8S Version
  kubespray_version: Kubespray Version
  import: Import
  import_status: Import status
  import_status_true: True
  import_status_false: False
  minVersionRequired: Min version required to kuboardspray
  cannot_reach_online_repository: Current browser cannot reach https://kuboard-spray.cn/support
zh:
  title: 资源包列表
  resourceDescription1: Kuboard 提供一组经过预先测试验证的资源包列表，可以帮助您快速完成集群安装
  resourceDescription2: 您也可以参考项目 https://github.com/eip-work/kuboard-spray-resource 自己创建资源包
  kuboardspray_version_min: KuboardSpray最低版本
  kuboardspray_version_max: KuboardSpray最高版本
  version: 资源包版本
  container_engine: 支持的容器引擎
  supported_os: 支持的操作系统
  kube_version: K8S 版本
  kubespray_version: Kubespray 版本
  import: 导 入
  import_status: 导入状态
  import_status_true: 已导入
  import_status_false: 未导入
  minVersionRequired: KuboardSpray最低版本要求
  cannot_reach_online_repository: 当前浏览器不能在线获取可选的资源包列表，请在可以访问外网的浏览器打开地址 https://kuboard-spray.cn/support
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
    <div style="text-align: right;">
      <el-link v-if="cannot_reach_online_repository" href="https://kuboard-spray.cn/support" target="blank" type="danger" icon="el-icon-link"
        style="margin-right: 10px; color: var(--el-color-danger)">{{ $t('cannot_reach_online_repository') }}</el-link>
      <ResourcesCreateOffline class="app_margin_top"></ResourcesCreateOffline>
    </div>
    <div class="contentList">
      <el-table v-if="mergedPackageList" :data="mergedPackageList" style="width: 100%">
        <el-table-column prop="version" :label="$t('version')" min-width="100px">
          <template #default="scope">
            <template v-if="hideLink">
              {{scope.row.version}}
            </template>
            <template v-else-if="importedPackageMap">
              <router-link v-if="importedPackageMap[scope.row.version]" :to="`/settings/resources/${scope.row.version}`">
                <!-- <i class="el-icon-link"></i> -->
                {{scope.row.version}}
              </router-link>
              <router-link v-else :to="`/settings/resources/${scope.row.version}/on_air`">
                <!-- <i class="el-icon-link"></i> -->
                {{scope.row.version}}
              </router-link>
            </template>
          </template>
        </el-table-column>
        <el-table-column :label="$t('kubespray_version')">
          <template #default="scope">
            <span v-if="packageYaml[scope.row.version]">
              {{ scope.row.yaml.data.kubespray_version }}
            </span>
            <i class="el-icon-loading" v-else></i>
          </template>
        </el-table-column>
        <el-table-column :label="$t('kube_version')">
          <template #default="scope">
            <span v-if="scope.row.yaml">
              {{ scope.row.yaml.data.kubernetes.kube_version }}
            </span>
            <i class="el-icon-loading" v-else></i>
          </template>
        </el-table-column>
        <el-table-column :label="$t('container_engine')">
          <template #default="scope">
            <template v-if="scope.row.yaml">
              <div v-for="(engine, key) in scope.row.yaml.data.container_engine" :key="`c${scope.index}_${key}`">
                <el-tag>{{ engine.container_manager }}_{{ engine.params[engine.container_manager + '_version'] }}</el-tag>
              </div>
            </template>
            <i class="el-icon-loading" v-else></i>
          </template>
        </el-table-column>
        <el-table-column :label="$t('supported_os')">
          <template #default="scope">
            <template v-if="scope.row.yaml">
              <div v-for="(os, key) in scope.row.yaml.metadata.supported_os" :key="`os${scope.index}_${key}`">
                <el-tag>
                  {{ os.distribution }}<span 
                  v-for="(v, i) in os.versions" :key="key + 'v' + i">_{{v}}</span>
                </el-tag>
              </div>
            </template>
            <i class="el-icon-loading" v-else></i>
          </template>
        </el-table-column>
        <el-table-column :label="$t('import_status')">
          <template #default="scope">
            <template v-if="importedPackageMap">
              <template v-if="importedPackageMap[scope.row.version]">
                <el-tag type="success" effect="dark">
                  <i class="el-icon-download" v-if="scope.row.isOffline"></i>
                  <i class="el-icon-cloudy" v-else></i>
                  {{ $t('import_status_true') }}
                </el-tag>
              </template>
              <el-tag v-else-if="scope.row.meetKuboardSprayVersion" type="warning" effect="dark">
                <i class="el-icon-circle-close"></i>
                {{ $t('import_status_false') }}
              </el-tag>
              <template v-else-if="scope.row.yaml">
                <el-tag type="error" effect="dark">{{ $t('minVersionRequired') }}</el-tag>
                <el-tag type="error" class="app_text_mono">{{scope.row.yaml.metadata.kuboard_spray_version.min}}</el-tag>
              </template>
            </template>
            <i class="el-icon-loading" v-else></i>
          </template>
        </el-table-column>
        <slot name="columns">
          <el-table-column :label="$t('msg.operations')" min-width="120px">
            <template #default="scope">
              <template v-if="importedPackageMap">
                <template v-if="importedPackageMap[scope.row.version]">
                  <el-button type="primary" plain icon="el-icon-view" @click="$router.push(`/settings/resources/${scope.row.version}`)">{{ $t('msg.view') }}</el-button>
                  <el-button type="danger" icon="el-icon-delete" @click="$router.push(`/settings/resources/${scope.row.version}`)">{{ $t('msg.delete') }}</el-button>
                </template>
                <template v-else>
                  <el-button type="primary" plain icon="el-icon-view" @click="$router.push(`/settings/resources/${scope.row.version}/on_air`)">{{ $t('msg.view') }}</el-button>
                  <el-button type="primary" v-if="scope.row.meetKuboardSprayVersion" icon="el-icon-download" @click="$router.push(`/settings/resources/${scope.row.version}/on_air`)">{{ $t('import') }}</el-button>
                </template>
              </template>
            </template>
          </el-table-column>
        </slot>
      </el-table>
    </div>

  </div>
</template>

<script>
import axios from 'axios'
import yaml from 'js-yaml'
import ResourcesCreateOffline from './ResourcesCreateOffline.vue'
import compareVersions from 'compare-versions'

export default {
  props: {
    hideLink: { type: Boolean, required: false, default: false }
  },
  data () {
    return {
      availablePackageList: undefined,
      importedPackageMap: {},
      packageYaml: {},
      cannot_reach_online_repository: false,
    }
  },
  computed: {
    mergedPackageList () {
      let result = []
      for (let i in this.importedPackageMap) {
        let flag = false
        for (let j in this.availablePackageList) {
          if (this.availablePackageList[j].version === i) {
            flag = true
          }
        }
        if (flag === false) {
          result.push({
            version: i,
            isOffline: true,
          })
        }
      }
      for (let i in this.availablePackageList) {
        result.push(this.availablePackageList[i])
      }
      return result
    }
  },
  components: { ResourcesCreateOffline },
  mounted () {
    this.refresh()
  },
  methods: {
    async refresh () {
      this.importedPackageMap = {}
      this.availablePackageList = undefined
      this.cannot_reach_online_repository = false
      await axios.get('https://addons.kuboard.cn/v-kuboard-spray/package-list.yaml?nocache=' + new Date().getTime()).then(resp => {
        this.availablePackageList = yaml.load(resp.data).items
      }).catch(e => {
        console.log(e)
        // this.$message.error('离线环境')
        this.cannot_reach_online_repository = true
      })
      await this.kuboardSprayApi.get(`/resources`).then(resp => {
        for (let i in resp.data.data) {
          this.importedPackageMap[resp.data.data[i]] = true
          for (let j in this.availablePackageList) {
            let pkg = this.availablePackageList[j]
            if (pkg.version === resp.data.data[i]) {
              pkg.imported = true
            }
          }
        }
      }).catch(e => {
        console.log(e)
      })
      for (let i in this.mergedPackageList) {
        let packageVersion = this.mergedPackageList[i]
        if (packageVersion.isOffline || packageVersion.imported) {
          this.loadPackageLocal(packageVersion)
        } else {
          this.loadPackageFromRepository(packageVersion)
        }
      }
    },
    loadPackageLocal(packageVersion) {
      this.kuboardSprayApi.get(`/resources/${packageVersion.version}`).then(resp => {
        this.packageYaml[packageVersion.version] = resp.data.data.package
        packageVersion.yaml = resp.data.data.package
        packageVersion.meetKuboardSprayVersion = compareVersions(window.KuboardSpray.version.trimed, packageVersion.yaml.metadata.kuboard_spray_version.min) >= 0
      })
    },
    loadPackageFromRepository (packageVersion) {
      axios.get(`https://addons.kuboard.cn/v-kuboard-spray/${packageVersion.version}/package.yaml?nocache=${new Date().getTime()}`).then(resp => {
        setTimeout(() => {
          let temp = yaml.load(resp.data)
          this.packageYaml[packageVersion.version] = temp
          packageVersion.yaml = temp
          packageVersion.meetKuboardSprayVersion = compareVersions(window.KuboardSpray.version.trimed, packageVersion.yaml.metadata.kuboard_spray_version.min) >= 0
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

