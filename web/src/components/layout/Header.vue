<i18n>
en:
  version: Current Version
  homePage: Home Page
  currentLogin: Current User
  logout: Logout
  loadingFailed: Failed in loadinng
zh:
  version: "当前版本"
  homePage: " 首页"
  currentLogin: "当前登录用户"
  logout: "退 出"
  loadingFailed: "加载失败"
</i18n>

<template>
  <div class="header-warpper" :style="transparent ? 'background: transparent': ''">
    <div :class="'header kuboard-header'">
      <el-progress class="progress" :percentage="percentage" status="success" :stroke-width="3"></el-progress>
      <div style="display: flex;" :class="transparent ? '' : 'header-background'">
        <div :class="`slot ${pageLevel ? 'slot-cluster' : ''}`" :style="transparent ? 'background: transparent': ''">
          <transition-group name="el-fade-in-linear" mode="out-in">
            <HeaderBreadCrumb :label="$t('homePage')" to="/" :kind="$t('homePage')" key="home" :hasNext="breadcrumb.length > 0"></HeaderBreadCrumb>
            <HeaderBreadCrumb v-for="(item, index) in breadcrumb" :key="index + 'bc' + item.label"
              :label="item.label" :to="item.to" :kind="item.kind" :hasNext="index < breadcrumb.length - 1"></HeaderBreadCrumb>
            <KbButton v-if="refresh" style="margin-left: 20px;" icon="el-icon-refresh" :loading="percentage < 100" key="refresh"
              @click="refresh.refresh.call(refresh.ref)"></KbButton>
          </transition-group>
        </div>
        <div class="header-right">
          <div class="msg">
            <div class="version" :style="showGithubStar ? '' : 'line-height: 70px;'">
              <div class="dot"></div>
              <div style="margin-right: 10px; width: 80px; display: inline-block; vertical-align: top; text-align: left;" class="nowrap">{{$t('version')}}</div>
              <span class="font-weight">{{ version }}</span>
            </div>
            <div class="version" v-if="showGithubStar">
              <div class="dot"></div>
              <iframe id="github-star-iframe" style="display:inline-block;vertical-align:middle;" src="https://addons.kuboard.cn/downloads/github-star-kuboard-spray.html" frameborder="0" scrolling="0" width="120" height="20" title="GitHub"></iframe>
            </div>
          </div>
          <LoginInfo></LoginInfo>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import HeaderBreadCrumb from './HeaderBreadCrumb.vue'
import KbButton from './KbButton.vue'
import LoginInfo from './LoginInfo.vue'
import axios from 'axios'

export default {
  props: {
    cluster: { type: String, required: false, default: undefined },
    namespace: { type: String, required: false, default: undefined },
    pageLevel: { type: Boolean, required: false, default: false },
    switchToNamespaceHome: { type: Boolean, required: false, default: false },
    transparent: { type: Boolean, required: false, default: false },
  },
  data () {
    return {
      dialogVisible: false,
      current: 0,
      showGithubStar: false,
    }
  },
  computed: {
    ...mapGetters({
      breadcrumb: 'header/breadcrumb',
      percentage: 'header/percentage',
      isClusterLevel: 'header/isClusterLevel',
      refresh: 'header/refresh',
    }),
    loading() {
      return this. percentage < 100
    },
    version() {
      return window.KuboardSpray.version.version
    }
  },
  // watch: {
  //   loading (value) {
  //     if (value) {
  //       this.$showErrorMsg()
  //     }
  //   }
  // },
  mounted () {
    this.interval = setInterval(() => {
      this.current = this.current + 1
    }, 5000)
    axios.get('https://addons.kuboard.cn/downloads/github-star-kuboard-spray.html').then(() => {
      this.showGithubStar = true
    }).catch(e => {
      console.log('hide github-star', e)
    })
  },
  beforeUnmount () {
    clearInterval(this.interval)
  },
  components: { HeaderBreadCrumb, KbButton, LoginInfo },
  methods: {
  }
}
</script>

<style>
.kuboard-header .el-progress__text {
  display: none;
}
.kuboard-header .el-progress {
  line-height: 2px;
  height: 2px;
}
.kuboard-header .el-progress .el-progress-bar {
  top: 0;
  position: fixed;
  left: 0;
  right: 0;
  margin-right: 0;
  padding-right: 0;
}
</style>

<style scoped lang="scss">
.version {
  // display: inline-block;
  width: 400px;
  text-overflow: ellipsis;
  overflow: hidden;
  white-space: nowrap;
  margin-right: 30px;
  user-select: none;
  vertical-align: top;
  font-size: 14px;
  color: $--color-text-regular;
  transition: 0.2s;
  line-height: 35px;
}
.version:hover {
  color: $--color-primary-light-5;
}
.version .dot {
  width: 6px;
  height: 6px;
  background-color: $--color-text-regular;
  display: inline-block;
  border-radius: 3px;
  margin-right: 8px;
  margin-bottom: 2px;
}
.msg {
  display: inline-block;
  width: 300px;
  text-align: left;
}
.setting {
  height: 70px;
  vertical-align: top;
  display: inline-block;
}
.font-weight {
  font-weight: 600;
}

.dropdown-menu {
  cursor: pointer;
}
.dropdown-menu :hover {
  color: $--color-primary;
}
.header {
  height: 72px;
  position: fixed;
  top: 0;
  z-index: 1000;
}
.header-background {
  // background:linear-gradient(0deg,#16386E,#1E4A90);
  background: transparent;
}
.header-warpper {
  height: 72px;
  background: transparent;
}

.bg {
  height: 70px;
  overflow: hidden;
}
.progress {
  width: 100%;
  overflow: hidden;
  // background-color: white;
}
.slot {
  height: 70px;
  overflow: hidden;
  padding: 0px 0 0px 20px;
  text-align: left;
  line-height: 70px;
  background: transparent;
  z-index: 1;
  color: $--color-white;
}
.slot-cluster {
  color: $--color-white;
}
.header-right {
  padding-right: 10px;
  position: fixed;
  top: 2px;
  right: 0;
  overflow: hidden;
  color: $--color-white;
  line-height: 70px;
  height: 70px;
  text-align: right;
}
.logo_label {
  font-size: 12px;
  color: mix($--color-white, $--color-kuboard-header, 70%);
}
.logo {
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  display: inline-block;
  color: white;
  padding: 0px 10px;
  width: 65px;
  border-radius: 4px;
  vertical-align: top;
  line-height: 50px;
  height: 50px;
  margin-top: 10px;
}
.logo:hover {
  color: $--color-kuboard-header;
  background-color: white;
}
.logo:active {
  color: $--color-primary;
  outline: none;
}
.yaml {
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  display: inline-block;
  padding: 0px 10px;
  width: 30px;
  margin-right: 8px;
  border-radius: 4px;
  vertical-align: top;
  line-height: 50px;
  margin-top: 10px;
  border: 1px solid rgb(217, 236, 255);
}

.settings {
  display: inline-block;
  height: 60px;
  margin-top: 5px;
}
.settings .button {
  height: 30px;
  line-height: 30px;
  color: white;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  padding: 0 15px;
  text-align: left;
  margin-right: 10px;
  width: 75px;
}
.settings .button:hover {
  color: $--color-kuboard-header;
  background-color: white;
}
.userInfo {
  font-size: 13px;
  width: 110px;
  display: inline-block;
  vertical-align: middle;
  cursor: pointer;
  font-weight: 800;
  color: white;
}
.userInfo:hover {
  color: $--color-primary;
}
.item {
  line-height: 1;
  height: 32px;
  vertical-align: top;
}
.label {
  color: $--color-text-secondary;
  width: 120px;
  display: inline-block;
  vertical-align: top;
  font-size: 13px;
  padding: 7px 0;
}
.value {
  color: $--color-text-primary;
  width: calc(100% - 150px);
  margin-left: 10px;
  margin-bottom: 10px;
  display: inline-block;
  vertical-align: top;
  font-size: 13px;
  height: 28px;
  padding: 7px 0;
  font-family: Monaco,Menlo,Consolas,Bitstream Vera Sans Mono,monospace;
}
.value button{
  margin-top: -7px;
}
.alertWrapper {
  text-align: left; padding: 0;
  position: absolute;
  top: 2px;
  z-index: 5000;
  width: calc(100vw - 0px);
}
</style>
