<i18n>
en:
  label: OS Mirror
  description: Operation System source repo.
  selectOs: Select OS
  source: Source
  addSource: Add Source
  asis: Use pre-configured source in the OS.
  docker_asis: Use docker official repository download.docker.com
  goToMirrorPage: It's about to open OS Software Source management page in new window, do you confirm?
zh:
  label: 软件源
  description: OS 软件源（为操作系统指定软件源，例如 yum 源、apt 源等）
  selectOs: 请选择操作系统
  source: 源
  addSource: 添加软件源
  asis: 使用操作系统预先配置的软件源
  docker_asis: 使用 Docker 官方软件源 download.docker.com
  goToMirrorPage: 此操作将将在新窗口打开软件源管理页面，完成设置后，可以切换窗口回到本页面，是否继续？
</i18n>

<template>
  <ConfigSection v-model:enabled="enabled" :label="$t('label')" :description="$t('description')" disabled anti-freeze>
    <el-alert title="离线安装请注意" type="warning" :closable="true" class="app_margin_bottom">
      <div style="line-height: 20px;">
        <li>安装 Kubernetes 集群时，需要用到一些系统软件，例如： curl, rsync, ipvsadm, ipset, ethtool 等；</li>
        <li>大部分企业都有自己的系统软件源，为了减小尺寸，KuboardSpray 资源包中不包含这些软件；</li>
      </div>
    </el-alert>
    <FieldCommon :holder="temp" fieldName="os" :label="$t('selectOs')" label-width="150px" anti-freeze>
      <template #edit>
        <el-checkbox-group v-model="os">
          <el-checkbox v-for="(item, index) in supportedOs" :key="'os_e' + index" :label="index.toLowerCase().replaceAll(' ', '_')">{{index}}</el-checkbox>
        </el-checkbox-group>
      </template>
      <template #view>
        <el-checkbox-group v-model="os" disabled>
          <el-checkbox v-for="(item, index) in supportedOs" :key="'os_v' + index" :label="index.toLowerCase().replaceAll(' ', '_')">{{index}}</el-checkbox>
        </el-checkbox-group>
      </template>
    </FieldCommon>
    <template v-for="(item, index) in os" :key="'repo' + index">
      <OsMirrorSelectRepo :holder="vars" :type="item" :prop="prop"></OsMirrorSelectRepo>
      <OsMirrorSelectRepo v-if="vars.container_manager === 'docker'" :holder="vars" :type="'docker_' + item" isdocker :prop="prop"></OsMirrorSelectRepo>
    </template>
  </ConfigSection>
</template>

<script>
import OsMirrorSelectRepo from './OsMirrorSelectRepo.vue'

export default {
  props: {
    cluster: { type: Object, required: true },
  },
  data() {
    return {
    }
  },
  computed: {
    enabled: {
      get () {return true},
      set () {},
    },
    supportedOs () {
      if (this.cluster === undefined) {
        return {}
      }
      if (this.cluster.resourcePackage === undefined) {
        return {}
      }
      let result = {}
      for (let item of this.cluster.resourcePackage.metadata.supported_os) {
        result[item.distribution] = true
      }
      return result
    },
    prop () { return 'all.children.target.vars' },
    vars() {
      return this.cluster.inventory.all.children.target.vars
    },
    os: {
      get () {
        let result = []
        for (let fieldName in this.vars) {
          if (fieldName.indexOf('kuboardspray_repo_') === 0 && fieldName.indexOf('kuboardspray_repo_docker_') !== 0) {
            result.push(fieldName.slice(18))
          }
        }
        return result
      },
      set (v) {
        if (v.length === 0) {
          this.$message.warning('cannot be empty')
          return
        }
        let t = {}
        for (let item of v) {
          t[item] = true
        }
        for (let os in this.supportedOs) {
          let key = os.toLowerCase()
          key = key.replaceAll(' ', '_')
          if (t[key]) {
            this.vars['kuboardspray_repo_' + key] = this.vars['kuboardspray_repo_' + key] || 'AS_IS'
          } else {
            delete this.vars['kuboardspray_repo_' + key]
          }
        }
      }
    },
    temp: {
      get () { return { os: this.os } },
      set () {}
    }
  },
  watch: {
    os: {
      deep: true,
      handler (newValue) {
        for (let item of newValue) {
          this.vars['kuboardspray_repo_docker_' + item] = this.vars['kuboardspray_repo_docker_' + item] || 'AS_IS'
        }
      }
    },
    'vars.container_manager': {
      deep: true,
      handler () {
        for (let item of this.os) {
          this.vars['kuboardspray_repo_docker_' + item] = this.vars['kuboardspray_repo_docker_' + item] || 'AS_IS'
        }
      }
    }
  },
  components: { OsMirrorSelectRepo },
  mounted () {
  },
  methods: {
  }
}
</script>

<style scoped lang="css">

</style>
