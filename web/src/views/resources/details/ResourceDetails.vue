<i18n>
en:
  os: Operation System
  install_by_default_true: install by default
  install_by_default_false: install on selection
  addon: Addons
  network_plugin: Cni plugins
  dependency: Dependencies
  package_content: Package Content
zh:
  os: 操作系统
  install_by_default_true: 默认安装
  install_by_default_false: 可选安装
  addon: 可选组件
  network_plugin: 网络插件
  dependency: 依赖组件
  package_content: 资源包内容
</i18n>

<template>
  <el-form class="app_form_mini" label-position="left" label-width="200px">
    <div style="text-align: center; margin-bottom: 10px; margin-top: -10px; font-weight: bold;">[ {{$t('package_content')}} ]</div>
    <el-collapse v-model="activeNames">
      <el-collapse-item name="1">
        <template #title>
          <span class="package_title">kuboardspray</span>
        </template>
        <div class="package_info">
          <PackageContentField :holder="resourcePackage.metadata" fieldName="version" label="资源包版本"></PackageContentField>
          <PackageContentField :holder="resourcePackage.metadata" fieldName="issue_date" label="发布时间"></PackageContentField>
          <PackageContentField :holder="resourcePackage.data" fieldName="kubespray_version"></PackageContentField>
        </div>
      </el-collapse-item>
      <el-collapse-item name="2">
        <template #title>
          <span class="package_title">kubernetes</span>
        </template>
        <div class="package_info">
          <el-form-item :label="$t('os')">
            <template v-for="(item, key) in resourcePackage.metadata.supported_os" :key="'os' + key">
              <div>
                <el-tag effect="dark">{{item.distribution}}</el-tag>
                <el-tag style="margin-left: 10px;" v-for="(version, i) in item.versions" :key="'v' + key + '' + i">{{version}}</el-tag>
              </div>
            </template>
          </el-form-item>
          <PackageContentField :holder="resourcePackage.data.kubernetes" fieldName="image_arch"></PackageContentField>
          <PackageContentField :holder="resourcePackage.data.kubernetes" fieldName="gcr_image_repo"></PackageContentField>
          <PackageContentField :holder="resourcePackage.data.kubernetes" fieldName="kube_image_repo"></PackageContentField>
          <PackageContentField :holder="resourcePackage.data.kubernetes" fieldName="kube_version"></PackageContentField>
          <el-form-item label="container_manager">
            <div v-for="(engine, index) in resourcePackage.data.container_engine" :key="'ce' + index">
              <el-tag>
                <span class="app_text_mono">{{engine.container_manager}}_{{engine.params.containerd_version || engine.params.docker_version}}</span>
              </el-tag>
            </div>
          </el-form-item>
        </div>
      </el-collapse-item>
      <el-collapse-item name="3">
        <template #title>
          <span class="package_title">etcd</span>
        </template>
        <div class="package_info">
          <PackageContentField :holder="resourcePackage.data.etcd" fieldName="etcd_version"></PackageContentField>
          <el-form-item :label="$t('field.etcd_deployment_type')">
            <div v-for="(item, key) in resourcePackage.data.etcd.etcd_deployment_type" :key="'k' + key">
              <el-tag>{{$t('field.etcd_deployment_type-' +item)}}</el-tag>
            </div>
          </el-form-item>
        </div>
      </el-collapse-item>
      <el-collapse-item name="4">
        <template #title>
          <span class="package_title">{{$t('network_plugin')}}</span>
        </template>
        <div class="package_info">
          <template v-for="(item, index) in resourcePackage.data.network_plugin" :key="index + 'network_plugin'">
              <div style="font-weight: bolder;">{{item.name}}</div>
              <div class="package_info">
                <el-form-item v-for="(value, key) in item.params" :key="index + 'p' + key" :label="key" label-width="180px">
                  <div class="app_text_mono" style="font-size: 13px" >
                    {{value}}
                  </div>
                </el-form-item>
              </div>
          </template>
        </div>
      </el-collapse-item>
      <el-collapse-item name="5">
        <template #title>
          <span class="package_title">{{$t('dependency')}}</span>
        </template>
        <div class="package_info">
          <template v-for="(item, index) in resourcePackage.data.dependency" :key="index + 'dependency'">
            <el-form-item :label="item.target">
              <div class="app_text_mono" style="font-size: 13px">
                {{item.version}}
              </div>
            </el-form-item>
          </template>
        </div>
      </el-collapse-item>
      <el-collapse-item name="6">
        <template #title>
          <span class="package_title">{{$t('addon')}}</span>
        </template>
        <div class="package_info">
          <template v-for="(item, index) in resourcePackage.data.addon" :key="index + 'addon'">
            <el-form-item>
              <template #label>
                <div style="font-weight: bolder;">{{item.name}}</div>
              </template>
              <div class="app_text_mono">
                <el-tag type="success" v-if="item.lifecycle && item.lifecycle.install_by_default">{{$t('install_by_default_true')}}</el-tag>
                <el-tag type="info" v-else>{{$t('install_by_default_false')}}</el-tag>
              </div>
            </el-form-item>
            <div class="package_info">
              <template v-for="(value, key) in item.params" :key="key + 'addon' + index">
                <el-form-item :label="key" label-width="280px" v-if="typeof value == 'string' && value.indexOf('{{') < 0">
                  <template #label>
                    <div style="font-size: 12px">{{key}}</div>
                  </template>
                  <div class="app_text_mono" style="font-size: 12px">
                    {{value}}
                  </div>
                </el-form-item>
              </template>
            </div>
          </template>
        </div>
      </el-collapse-item>
    </el-collapse>
  </el-form>
</template>

<script>
import PackageContentField from './PackageContentField.vue'

export default {
  props: {
    resourcePackage: { type: Object, required: true },
    expandAll: { type: Boolean, required: false, default: false },
  },
  data() {
    return {
      activeNames: this.expandAll ? ['1','2','3','4','5','6'] : ['2'],
    }
  },
  computed: {
  },
  components: { PackageContentField },
  mounted () {
  },
  methods: {

  }
}
</script>

<style scoped lang="scss">
.package_title {
  font-weight: bolder;
  font-family: Consolas,Menlo,Bitstream Vera Sans Mono,Monaco,"微软雅黑",monospace;
}
.package_info {
  margin-left: 20px;
  font-family: Consolas,Menlo,Bitstream Vera Sans Mono,Monaco,"微软雅黑",monospace;
}
</style>
