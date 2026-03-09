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
  <el-form class="app_form_mini" label-position="left" label-width="300px">
    <div style="text-align: center; margin-bottom: 10px; margin-top: -10px; font-weight: bold;">[
      {{ t('package_content') }} ]</div>
    <el-collapse v-model="activeNames">
      <el-collapse-item name="1">
        <template #title>
          <span class="package_title">Pangee Cluster</span>
        </template>
        <div class="package_info">
          <PackageContentField :holder="resourcePackage.metadata" fieldName="version" label="资源包版本">
          </PackageContentField>
          <PackageContentField :holder="resourcePackage.metadata" fieldName="issue_date" label="发布时间">
          </PackageContentField>
          <el-form-item :label="t('os')">
            <div>
              <template v-for="(item, key) in resourcePackage.metadata.supported_os" :key="'os' + key">
                <div>
                  <el-tag effect="dark">{{ item.distribution }}</el-tag>
                  <el-tag style="margin-left: 10px;" v-for="(version, i) in item.versions" :key="'v' + key + '' + i">{{
                    version }}</el-tag>
                </div>
              </template>
            </div>
          </el-form-item>
        </div>
      </el-collapse-item>
      <el-collapse-item name="2">
        <template #title>
          <span class="package_title">{{ t('dependency') }}</span>
        </template>
        <div class="package_info">
          <template v-for="(item, index) in resourcePackage.data.dependency" :key="index + 'dependency'">
            <el-form-item :label="item.name">
              <div class="app_text_mono" style="font-size: 13px">
                {{ item.version }}
              </div>
            </el-form-item>
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
      activeNames: this.expandAll ? ['1', '2', '3', '4', '5', '6'] : ['1', '2'],
    }
  },
  computed: {
  },
  components: { PackageContentField },
  mounted() {
  },
  methods: {

  }
}
</script>

<style scoped lang="css">
.package_title {
  font-weight: bolder;
  font-family: Consolas, Menlo, Bitstream Vera Sans Mono, Monaco, "微软雅黑", monospace;
}

.package_info {
  margin-left: 20px;
  font-family: Consolas, Menlo, Bitstream Vera Sans Mono, Monaco, "微软雅黑", monospace;
}
</style>
