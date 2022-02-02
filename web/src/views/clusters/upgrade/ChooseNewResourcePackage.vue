<i18n>
en:
  title: Choose resource package to upgrade to.
  currentVersion: Currently installed
  cannot_down_grade: K8S version in the package is lower
zh:
  title: 选择升级时使用的资源包
  currentVersion: 当前已安装的版本
  cannot_down_grade: 资源包中K8S版本过低
</i18n>

<template>
  <div>
    <el-dialog v-model="dialogVisible" :title="$t('title')" width="80%" top="5vh" :close-on-click-modal="false">
      <Resources v-if="dialogVisible" hide-link>
        <template #columns>
          <el-table-column :label="$t('msg.operations')" min-width="120px">
            <template #default="scope">
              <el-tag v-if="currentVersion === scope.row.version" type="success" effect="dark">
                <i class="el-icon-circle-check"></i>
                {{ $t('currentVersion') }}
              </el-tag>
              <template v-else-if="scope.row.yaml">
                <el-tag v-if="compareVersions(cluster.resourcePackage.data.kubernetes.kube_version, scope.row.yaml.data.kubernetes.kube_version) > 0"
                  type="info" effect="dark">
                  {{$t('cannot_down_grade')}}
                </el-tag>
                <template v-else>
                  <el-button icon="el-icon-view" type="primary" plain @click="$refs.comparePackage.show(scope.row)">{{ $t('msg.view') }}</el-button>
                </template>
                <el-button icon="el-icon-view" type="primary" plain @click="$refs.comparePackage.show(scope.row)">{{ $t('msg.view') }}</el-button>
              </template>
            </template>
          </el-table-column>
        </template>
      </Resources>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false" icon="el-icon-close" type="primary">{{$t('msg.cancel')}}</el-button>
        </span>
      </template>
    </el-dialog>
    <ComparePackage ref="comparePackage" :cluster="cluster"></ComparePackage>
  </div>
</template>

<script>
import Resources from '../../resources/Resources.vue'
import ComparePackage from './ComparePackage.vue'
import compareVersions from 'compare-versions'

export default {
  props: {
    cluster: { type: Object, required: true },
  },
  data() {
    return {
      dialogVisible: false,
    }
  },
  computed: {
    currentVersion () {
      return this.cluster.inventory.all.hosts.localhost.kuboardspray_resource_package
    }
  },
  components: { Resources, ComparePackage },
  mounted () {
  },
  methods: {
    compareVersions,
    show() {
      this.dialogVisible = true
    },
    upgradeToType(target) {
      console.log(target)
    }
  }
}
</script>

<style scoped lang="scss">

</style>
