<i18n>
en:
  included: Included
  not_included: Not included
  addons: Addons
  cni: Cni plugins
  dependency: Dependencies
  package_content: Package Content
zh:
  included: 已包含
  not_included: 未包含
  addons: 可选组件
  cni: 网络插件
  dependency: 依赖组件
  package_content: 资源包内容
</i18n>

<template>
  <el-form class="app_form_mini" label-position="left" label-width="200px">
    <div style="text-align: center; margin-bottom: 10px; margin-top: -10px; font-weight: bold;">[ {{$t('package_content')}} ]</div>
    <el-collapse v-model="activeNames">
      <el-collapse-item name="1">
        <template #title>
          <span class="package_title">kubespray</span>
        </template>
        <div class="package_info">
          <PackageContentField :holder="resourcePackage.kuboardspray" fieldName="kubespray_version"></PackageContentField>
        </div>
      </el-collapse-item>
      <el-collapse-item name="2">
        <template #title>
          <span class="package_title">kubernetes</span>
        </template>
        <div class="package_info">
          <PackageContentField :holder="resourcePackage.kubernetes" fieldName="image_arch"></PackageContentField>
          <PackageContentField :holder="resourcePackage.kubernetes" fieldName="gcr_image_repo"></PackageContentField>
          <PackageContentField :holder="resourcePackage.kubernetes" fieldName="kube_image_repo"></PackageContentField>
          <PackageContentField :holder="resourcePackage.kubernetes" fieldName="kube_version"></PackageContentField>
          <el-form-item label="container_manager">
            <div v-for="(engine, index) in resourcePackage.container_engine" :key="'ce' + index">
              <el-tag>
                <span class="app_text_mono">{{engine.container_manager}}{{engine.version ? '_' : ''}}{{engine.version}}</span>
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
          <PackageContentField :holder="resourcePackage.etcd" fieldName="etcd_version"></PackageContentField>
          <PackageContentField :holder="resourcePackage.etcd" fieldName="etcd_deployment_type"></PackageContentField>
        </div>
      </el-collapse-item>
      <el-collapse-item name="4">
        <template #title>
          <span class="package_title">{{$t('cni')}}</span>
        </template>
        <div class="package_info">
          <template v-for="(item, index) in resourcePackage.cni" :key="index + 'cni'">
            <el-form-item :label="item.target">
              <div class="app_text_mono">
                {{item.version}}
              </div>
            </el-form-item>
          </template>
        </div>
      </el-collapse-item>
      <el-collapse-item name="5">
        <template #title>
          <span class="package_title">{{$t('dependency')}}</span>
        </template>
        <div class="package_info">
          <template v-for="(item, index) in resourcePackage.dependency" :key="index + 'dependency'">
            <el-form-item :label="item.target">
              <div class="app_text_mono">
                {{item.version}}
              </div>
            </el-form-item>
          </template>
        </div>
      </el-collapse-item>
      <el-collapse-item name="6">
        <template #title>
          <span class="package_title">{{$t('addons')}}</span>
        </template>
        <div class="package_info">
          <template v-for="(item, index) in resourcePackage.addons" :key="index + 'addons'">
            <el-form-item>
              <template #label>
                <div style="font-weight: bolder;">{{item.name}}</div>
              </template>
              <div class="app_text_mono">
                <el-tag type="success" v-if="item.included">{{$t('included')}}</el-tag>
                <el-tag type="info" v-else>{{$t('not_included')}}</el-tag>
              </div>
            </el-form-item>
            <div class="package_info">
              <template v-for="(value, key) in item.params" :key="key + 'addons' + index">
                <el-form-item :label="key" label-width="240px" v-if="value.indexOf('{{') < 0">
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
