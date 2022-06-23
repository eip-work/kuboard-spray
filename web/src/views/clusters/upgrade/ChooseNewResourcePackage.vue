<i18n>
en:
  title: Choose resource package to upgrade to.
  currentVersion: Currently selected
  cannot_down_grade: K8S version in the package is lower
  switchToTargetVersion: Replace with this version
  switchToTargetVersionDesc: Component version in this resource package are the same with the currently installed one, but ansible script are different. You don't need to perform a cluster upgrade by selecting this resource package.
  upgradeToTargetVersion: Upgrade to this version
  upgradeToTargetVersionDesc: After select this resource package you need to perform a cluster upgrade to keep with the latest version.
  cannotUpgradeToTargetVersion: Cannot upgrade to this version
  downloadFirstCanUpgrade: Can upgrade, go download
  downloadFirstCanReplace: Can replace, go download
  change_resource_package_version_success: Succeed in changing resource package version
  change_resource_package_version_failed: Failed to change resource package version {msg}

  not_supported_cm: Container manager not match
zh:
  title: 选择升级时使用的资源包
  currentVersion: 当前版本
  cannot_down_grade: 资源包中K8S版本过低
  switchToTargetVersion: 用此版本替代
  switchToTargetVersionDesc: 此资源包中的组件版本与当前已安装版本相同，但是 ansible 脚本有所不同，替换后，不需要对集群执行升级操作。
  upgradeToTargetVersion: 升级到此版本
  upgradeToTargetVersionDesc: 选择此版本后，您需要执行集群升级动作，以更新集群的版本。
  cannotUpgradeToTargetVersion: 不能升级到此版本
  downloadFirstCanUpgrade: 可升级，去下载
  downloadFirstCanReplace: 可替代，去下载
  change_resource_package_version_success: 已成功修改资源包版本
  change_resource_package_version_failed: 修改资源包版本失败 {msg}

  not_supported_cm: 不支持当前容器引擎
</i18n>

<template>
  <div>
    <el-dialog v-model="dialogVisible" :title="$t('title')" width="80%" top="5vh" :close-on-click-modal="false">
      <Resources v-if="dialogVisible" hide-link>
        <template #columns>
          <el-table-column :label="$t('msg.operations')" min-width="130px">
            <template #default="scope">
              <el-tag v-if="currentVersion === scope.row.version" type="success" effect="dark">
                <el-icon :size="14" style="width: 14px; height: 14px; vertical-align: bottom;">
                  <el-icon-circle-check></el-icon-circle-check>
                </el-icon>
                {{ $t('currentVersion') }}
              </el-tag>
              <template v-else-if="scope.row.yaml">
                <el-tag v-if="compareVersions(cluster.resourcePackage.data.kubernetes.kube_version, scope.row.yaml.data.kubernetes.kube_version) > 0"
                  type="info" effect="dark">
                  {{$t('cannot_down_grade')}}
                </el-tag>
                <template v-else>
                  <el-button icon="el-icon-view" type="primary" plain @click="$refs.comparePackage.show(scope.row)">{{ $t('msg.view') }}</el-button>
                  <template v-if="!matchComponents(scope.row)">
                    <el-button disabled icon="el-icon-circle-close" type="info">{{ $t('not_supported_cm') }}</el-button>
                  </template>
                  <template v-else-if="canReplaceBy(scope.row)">
                    <el-button type="primary" icon="el-icon-download" v-if="!scope.row.isOffline && !scope.row.imported"
                      @click="$router.push(`/settings/resources/${scope.row.version}/on_air`)">
                      {{ $t('downloadFirstCanReplace') }}
                    </el-button>
                    <confirm-button v-else type="warning" icon="el-icon-right" :text="$t('switchToTargetVersion')"
                      :message="$t('switchToTargetVersionDesc')"
                      @confirm="changeToVersion('replace', scope.row.version)">
                    </confirm-button>
                  </template>
                  <template v-else-if="canUpgradeTo(scope.row)">
                    <el-button type="primary" icon="el-icon-download" v-if="!scope.row.isOffline && !scope.row.imported"
                      @click="$router.push(`/settings/resources/${scope.row.version}/on_air`)">
                      {{ $t('downloadFirstCanUpgrade') }}
                    </el-button>
                    <confirm-button v-else type="warning" icon="el-icon-right" :text="$t('upgradeToTargetVersion')"
                      :message="$t('upgradeToTargetVersionDesc')"
                      @confirm="changeToVersion('upgrade', scope.row.version)">
                    </confirm-button>
                  </template>
                  <el-button v-else type="info" disabled icon="el-icon-circle-close">{{ $t('cannotUpgradeToTargetVersion') }}</el-button>
                </template>
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
    },
  },
  components: { Resources, ComparePackage },
  mounted () {
  },
  methods: {
    compareVersions,
    show() {
      this.dialogVisible = true
    },
    canUpgradeTo(target) {
      return this.canTarget(target.yaml.metadata.can_upgrade_from, this.cluster.resourcePackage.metadata.version)
    },
    canReplaceBy(target) {
      return this.canTarget(target.yaml.metadata.can_replace_to, this.cluster.resourcePackage.metadata.version)
    },
    matchComponents(target) {
      let container_manager = this.cluster.inventory.all.children.target.vars.container_manager
      for (let ce_index in target.yaml.data.container_engine) {
        if (target.yaml.data.container_engine[ce_index].container_manager == container_manager) {
          return true
        }
      }
      return false
    },
    canTarget(rule, originalVersion) {
      if (!rule) {
        return false
      }
      for (let k in rule.exclude) {
        let p = new RegExp('^' + rule.exclude[k] + '$')
        // console.log('exclude: ', p, ' ', originalVersion, '', p.test(originalVersion))
        if (p.test(originalVersion)) {
          return false
        }
      }
      for (let k in rule.include) {
        let p = new RegExp('^' + rule.include[k] + '$')
        // console.log('include: ', p, ' ', originalVersion, '', p.test(originalVersion))
        if (p.test(originalVersion)) {
          return true
        }
      }
      return false
    },
    changeToVersion(type, targetVersion) {
      let req = {
        type: type,
        target_version: targetVersion,
      }
      this.kuboardSprayApi.post(`/clusters/${this.cluster.name}/change_resource_package_version`, req).then(resp => {
        console.log(resp.data)
        this.dialogVisible = false
        this.$message.success(this.$t('change_resource_package_version_success'))
        this.$emit('refresh')
      }).catch(e => {
        console.log(e)
        let msg = e + ''
        if (e.response && e.response.data && e.response.data.message) {
          msg = e.response.data.message
        }
        this.$message.error(this.$t('change_resource_package_version_failed', {msg}))
      })
    }
  }
}
</script>

<style scoped lang="css">

</style>
