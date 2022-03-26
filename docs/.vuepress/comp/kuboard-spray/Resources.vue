<template>
  <div>
    <el-alert v-if="arch === ''" title="资源包列表" type="warning" :closable="false" show style="margin-bottom: 10px;">
      <div class="description">
        <li>Kuboard 提供一组经过预先测试验证的资源包列表，可以帮助您快速完成集群安装</li>
        <li>您也可以参考项目 https://github.com/eip-work/kuboard-spray-resource 自己创建资源包</li>
      </div>
    </el-alert>
    <div v-if="arch === '-arm64'" class="type">ARM 架构</div>
    <div v-else class="type">AMD 架构</div>
    <div class="contentList" v-if="mergedPackageList">
      <el-table :data="mergedPackageList">
        <el-table-column prop="version" label="版 本"></el-table-column>
        <el-table-column prop="package.data.kubespray_version" label="kubespray"></el-table-column>
        <el-table-column prop="package.data.kubernetes.kube_version" label="kubernetes"></el-table-column>
        <el-table-column prop="package.data.container_engine" label="容器引擎">
          <template slot-scope="scope">
            <div v-if="scope.row.package" style="margin-bottom: -10px;">
              <div v-for="(engine, key) in scope.row.package.data.container_engine" :key="`c${scope.index}_${key}`">
                <el-tag style="margin-bottom: 10px;">{{ engine.container_manager }}_{{ engine.params[engine.container_manager + '_version'] }}</el-tag>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="package.metadata.supported_os" label="操作系统">
          <template #default="data">
            <div v-if="data.row.package" style="margin-bottom: -10px;">
              <div v-for="(os, key) in data.row.package.metadata.supported_os" :key="`os${data.index}_${key}`">
                <el-tag style="margin-bottom: 10px;">
                  {{ os.distribution }}<span 
                  v-for="(v, i) in os.versions" :key="key + 'v' + i">_{{v}}</span>
                </el-tag>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="package.metadata.kuboard_spray_version.min" label="KuboardSpray最低版本"></el-table-column>
        <el-table-column label="操 作">
          <template #default="data">
            <el-button type="primary" @click="showVersion(data.row)">离线导入</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-dialog :visible.sync="dialogVisible" top="3vh" width="1200px" title="离线导入资源包">
      <el-tabs v-model="activeName" type="card">
        <el-tab-pane name="details" label="资源包内容">
          <div class="app_form_mini resource_content">
            <ResourceDetails v-if="currentVersion" ref="details" expandAll :resourcePackage="currentVersion"></ResourceDetails>
          </div>
        </el-tab-pane>
        <el-tab-pane name="offline" label="离线导入">
          <div class="app_form_mini resource_content" v-if="activeName == 'offline'">
            <ClientOnly>
            <ResourceOffline v-if="currentVersion" :resourcePackage="currentVersion"></ResourceOffline>
            </ClientOnly>
          </div>
        </el-tab-pane>
      </el-tabs>
      <template #modal-title>
        离线导入资源包
      </template>
      <template #footer>
        <el-button @click="dialogVisible = false">关闭对话框</el-button>
        <el-button type="primary" @click="activeName = 'offline'; $message.success('已切换到离线导入页，请按提示操作。')">离线导入</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import axios from 'axios'
import yaml from 'js-yaml'
import ResourceDetails from './ResourceDetails.vue'
import ResourceOffline from './ResourceOffline.vue'


export default {
  props: {
    arch: { type: String, required: false, default: '' }
  },
  data () {
    return {
      dialogVisible: false,
      availablePackageList: undefined,
      packageMap: {},
      importedPackageMap: {},
      currentVersion: undefined,
      activeName: 'details',
    }
  },
  computed: {
    mergedPackageList () {
      let result = []
      for (let i in this.availablePackageList) {
        result.push(this.availablePackageList[i])
      }
      return result
    }
  },
  components: { 
    ResourceDetails,
    ResourceOffline,
  },
  mounted () {
    this.refresh()
  },
  methods: {

    async refresh () {
      this.importedPackageMap = {}
      this.packageMap = {}
      this.availablePackageList = undefined
      await axios.get(`https://addons.kuboard.cn/v-kuboard-spray${this.arch}/package-list.yaml?nocache=` + new Date().getTime()).then(resp => {
        this.availablePackageList = yaml.load(resp.data).items
      }).catch(e => {
        console.log(e)
        this.$message.error('离线环境')
      })

      for (let i in this.mergedPackageList) {
        let packageVersion = this.mergedPackageList[i]
        this.loadPackageFromRepository(packageVersion)
      }
    },
    loadPackageFromRepository (packageVersion) {
      axios.get(`https://addons.kuboard.cn/v-kuboard-spray${this.arch}/${packageVersion.version}/package.yaml?nocache=${new Date().getTime()}`).then(resp => {
        setTimeout(() => {
          // packageVersion.package = yaml.load(resp.data)
          // packageVersion.loaded = true
          this.$set(packageVersion, 'package', yaml.load(resp.data))
          this.$set(packageVersion, 'loaded', true)
          // this.packageMap[packageVersion.version].loaded = true
        }, 500)
      }).catch(e => {
        this.$message.error(e + '')
        this.packageMap[packageVersion.version].loaded = false
      })
    },
    showVersion (version) {
      this.currentVersion = version.package
      console.log(version)
      this.dialogVisible = true
    }
  }
}
</script>

<style>
.app_form_mini .el-form-item {
  margin-bottom: 0px !important;
}
.app_form_mini .el-form-item__label {
  margin-bottom: 0px !important;
}
.el-table table {
  margin: 0px;
}
.el-table th {
  padding: 0px !important;
}

</style>

<style scoped lang="css">
.description {
  line-height: 28px;
}
.contentList {
  width: 100%;
  margin-bottom: 40px;
}
.app_text_mono {
  font-family: Consolas,Menlo,Bitstream Vera Sans Mono,Monaco,"微软雅黑",monospace !important;
}
.resource_tab_pane {
  height: calc(100vh - 300px);
  overflow: hidden;
  overflow-y: auto;
}
.resource_content {
  height: calc(94vh - 270px);
  overflow: hidden;
  overflow-y: auto;
}
.type {
  background-color: #007af5;
  color: white;
  padding: 10px 20px;
  font-family: Consolas,Menlo,Bitstream Vera Sans Mono,Monaco,"微软雅黑",monospace !important;
  font-weight: bolder;
}
</style>

