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
  install_docker_from_distro: SUSE distributions always install Docker from the distro repos
zh:
  label: 软件源
  description: OS 软件源（为操作系统指定软件源，例如 yum 源、apt 源等）
  selectOs: 请选择操作系统
  source: 源
  addSource: 添加软件源
  asis: 使用操作系统预先配置的软件源
  docker_asis: 使用 Docker 官方软件源 download.docker.com
  goToMirrorPage: 此操作将将在新窗口打开软件源管理页面，完成设置后，可以切换窗口回到本页面，是否继续？
  install_docker_from_distro: SUSE 将始终从操作系统的软件源中安装 Docker
</i18n>

<template>
  <FieldCommon :holder="vars" :fieldName="fieldName" label-width="150px" :placeholder="type.indexOf('opensuse') >= 0 ? $t('install_docker_from_distro') : $t('asis')"
    :label="type + ' ' + $t('source')" :required="type.indexOf('opensuse') < 0" :prop="prop" anti-freeze>
    <template #edit>
      <template v-if="type.indexOf('opensuse') >= 0">
        <span v-if="type.indexOf('docker_') === 0">{{ $t('install_docker_from_distro') }}</span>
        <span v-else>{{ $t('asis') }}</span>
      </template>
      <div style="display: flex;" v-else>
        <el-select v-model="vars[fieldName]" @visible-change="loadKuboardsprayRepoOptions" style="flex-grow: 1;">
          <template v-for="(option, index) in options" :key="fieldName + index">
            <el-option :value="option.value" class="kuboardspray_mirror_select" :label="option.label">
              <div>
                <div>{{option.label}}</div>
                <div v-if="option.data">
                  <div v-for="(value, key) in option.data.status.params" :key="fieldName + index + key"
                    style="margin-left: 20px;">
                    <el-tag effect="dark">{{key}}</el-tag>
                    <el-tag>{{value}}</el-tag>
                  </div>
                </div>
                <div v-else-if="option.value !== 'AS_IS'">
                  <el-icon class="is-loading">
                    <el-icon-loading></el-icon-loading>
                  </el-icon>
                </div>
              </div>
            </el-option>
          </template>
        </el-select>
        <ConfirmButton @confirm="openUrlInBlank('/#/settings/mirrors')" buttonStyle="margin-left: 10px;" placement="top-end"
          icon="el-icon-plus" plain :text="$t('addSource')" :message="$t('goToMirrorPage')"></ConfirmButton>
      </div>
      <template v-for="(option, index) in options" :key="fieldName + 'ep' + index">
        <div v-if="option.value === vars[fieldName]">
          <div v-if="option.data" class="preview">
            <div v-for="(value, key) in option.data.status.params" :key="fieldName + index + key">
              <el-tag effect="dark">{{key}}</el-tag>
              <el-tag>{{value}}</el-tag>
            </div>
          </div>
        </div>
      </template>
    </template>
    <template #view>
      <span v-if="type.indexOf('opensuse') >= 0" style="color: var(--el-text-color-placeholder)">
        {{$t('asis')}}
      </span>
      <template v-for="(option, index) in options" :key="fieldName + 'ep' + index">
        <div v-if="option.value === vars[fieldName]">
          <div>
            <div>{{option.label}}</div>
            <div v-if="option.data" class="preview">
              <div v-for="(value, key) in option.data.status.params" :key="fieldName + index + key">
                <el-tag effect="dark">{{key}}</el-tag>
                <el-tag>{{value}}</el-tag>
              </div>
            </div>
          </div>
        </div>
      </template>
    </template>
  </FieldCommon>
</template>

<script>
export default {
  props: {
    holder: { type: Object, required: true },
    type: { type: String, required: true },
    prop: { type: String, required: false, default: '' },
    isdocker: { type: Boolean, required: false, default: false },
  },
  data() {
    return {
      options: []
    }
  },
  computed: {
    fieldName() {
      return 'kuboardspray_repo_' + (this.isdocker ? 'docker_': '')+ this.type
    },
    vars: {
      get () { return this.holder },
      set () {},
    }
  },
  components: { },
  mounted () {
    this.loadKuboardsprayRepoOptions(true, true)
  },
  methods: {
    loadKuboardsprayRepoOptions(flag, loadAll) {
      if (flag) {
        this.options = []
        let type = (this.isdocker ? 'docker_': '')+ this.type
        this.kuboardSprayApi.get(`/mirrors`, { params: { type: type } }).then(async resp => {
          if (type.indexOf('docker_') !== 0) {
            this.options.push({ label: this.$t('asis'), value: 'AS_IS'})
          } else {
            this.options.push({ label: this.$t('docker_asis'), value: 'AS_IS'})
          }
          for (let item of resp.data.data) {
            let name = item.slice(item.indexOf('-') + 1)
            let option = { label: name, value: item }
            if (!loadAll || item === this.vars[this.fieldName]) {
              await this.kuboardSprayApi.get(`/mirrors/${item}`).then(resp => {
                option.data = resp.data.data
              })
            }
            this.options.push(option)
          }
        }).catch(e => {
          console.log(e)
        })
      }
    },
  }
}
</script>

<style lang="scss">
.kuboardspray_mirror_select.el-select-dropdown__item {
  height: auto;
}
</style>

<style scoped lang="scss">
.preview {
  padding: 10px;
  background-color: var(--el-color-info-lighter);
}
</style>
