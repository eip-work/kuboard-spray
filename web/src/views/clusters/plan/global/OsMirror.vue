<i18n>
en:
  label: OS Mirror
  description: Operation System source repo.
  selectOs: Select OS
zh:
  label: 软件源
  description: OS 软件源（为操作系统指定软件源，例如 yum 源、apt 源等）
  selectOs: 请选择操作系统
</i18n>

<template>
  <ConfigSection v-model:enabled="enabled" :label="$t('label')" :description="$t('description')" disabled anti-freeze>
    <el-alert title="离线安装请注意" type="warning" :closable="true" class="app_margin_bottom">
      <div style="line-height: 20px;">
        <li>安装 Kubernetes 集群时，需要用到一些系统软件，例如： curl, rsync, ipvsadm, ipset, ethtool 等；</li>
        <li>大部分企业都有自己的系统软件源，为了减小尺寸，KuboardSpray 资源包中不包含这些软件；</li>
      </div>
    </el-alert>
    <FieldCommon :holder="temp" fieldName="os" :label="$t('selectOs')" required>
      <template #edit>
        <el-checkbox-group v-model="os">
          <el-checkbox v-for="(item, index) in supportedOs" :key="'os_e' + index" :label="index"></el-checkbox>
        </el-checkbox-group>
      </template>
      <template #view>
        <el-checkbox-group v-model="temp.os" disabled>
          <el-checkbox v-for="(item, index) in supportedOs" :key="'os_v' + index" :label="index"></el-checkbox>
        </el-checkbox-group>
      </template>
    </FieldCommon>
    <template v-for="(item, index) in os" :key="'repo' + index">
      <FieldString :holder="vars" :fieldName="repoMap[item].repo"></FieldString>
      <template v-if="vars.container_manager === 'docker'">
        <FieldString :holder="vars" :fieldName="repoMap[item].dockerRepo"></FieldString>
        <FieldString :holder="vars" :fieldName="repoMap[item].dockerRepoGpgkey"></FieldString>
        <FieldString v-if="repoMap[item].dockerRepoRepokey" :holder="vars" :fieldName="repoMap[item].dockerRepoRepokey"></FieldString>
      </template>
    </template>
    <!-- <FieldString :holder="cluster.inventory.all.children.target.vars" fieldName="yum_repo"></FieldString>
    <FieldString :holder="cluster.inventory.all.children.target.vars" fieldName="docker_ubuntu_repo_base_url"></FieldString>
    <FieldString :holder="cluster.inventory.all.children.target.vars" fieldName="docker_ubuntu_repo_gpgkey"></FieldString>
    <FieldString :holder="cluster.inventory.all.children.target.vars" fieldName="docker_ubuntu_repo_repokey"></FieldString> -->
  </ConfigSection>
</template>

<script>
export default {
  props: {
    cluster: { type: Object, required: true },
  },
  data() {
    return {
      repoMap: {
        Ubuntu: {
          repo: 'ubuntu_repo',
          dockerRepo: 'docker_ubuntu_repo_base_url',
          dockerRepoGpgkey: 'docker_ubuntu_repo_gpgkey',
          dockerRepoRepokey: 'docker_ubuntu_repo_repokey',
        },
        CentOS: {
          repo: 'centos_repo',
          dockerRepo: 'docker_rh_repo_base_url',
          dockerRepoGpgkey: 'docker_rh_repo_gpgkey',
        }
      }
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
      for (let item of this.cluster.resourcePackage.supported_os) {
        result[item.distribution] = true
      }
      return result
    },
    vars() {
      return this.cluster.inventory.all.children.target.vars
    },
    os: {
      get () {
        let result = []
        for (let fieldName in this.vars) {
          for (let os in this.repoMap) {
            if (fieldName === this.repoMap[os].repo) {
              result.push(os)
            }
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
        for (let key in this.repoMap) {
          if (t[key]) {
            this.vars[this.repoMap[key].repo] = this.vars[this.repoMap[key].repo] || undefined
          } else {
            delete this.vars[this.repoMap[key].repo]
          }
        }
      }
    },
    temp: {
      get () { return { os: this.os } },
      set () {}
    }
  },
  components: { },
  mounted () {
  },
  methods: {

  }
}
</script>

<style scoped lang="scss">

</style>
