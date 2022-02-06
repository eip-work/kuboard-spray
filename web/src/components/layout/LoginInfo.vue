<i18n>
en:
  currentLogin: "Current Login User"
  logout: "Logout"
  loadingFailed: "Failed in loadinng"
  username: User Name
  currentCluster: "Current Cluster"
  logoutAction: "Please go to {type} to logout"
  stayLogin: "Cancel logout, stay login"
  changePassword: Change Password
zh:
  currentLogin: "当前登录用户"
  logout: "退 出"
  loadingFailed: "加载失败"
  username: 用户名
  currentCluster: "当前集群"
  logoutAction: "请到 {type} 退出登录"
  stayLogin: "取消退出，保持登录状态"
  changePassword: 修改密码
</i18n>

<template>
  <div class="user">
    <el-popover v-model:visible="visible" placement="bottom-end" :width="540" trigger="mannual">
      <template #reference>
        <div @click="visible = !visible">
          <span class="font-weight" style="margin-right: 5px;">{{displayName}}</span>
          <el-icon>
            <el-icon-caret-bottom></el-icon-caret-bottom>
          </el-icon>
        </div>
      </template>
      <div style="padding: 20px; max-height: calc(100vh - 160px); overflow: hidden; overflow-y: auto;">
        <div class="item">
          <div class="label">{{$t('username')}}</div>
          <div class="value">
            admin
          </div>
        </div>
        <div class="item">
          <div class="label">{{$t('changePassword')}}</div>
          <div class="value">
            <el-button type="text" @click="visible = false; $router.push('/settings/profile')" icon="el-icon-link">{{$t('changePassword')}}</el-button>
          </div>
        </div>
        <div class="app_text_right" style="float: right; margin-top: 20px;">
          <el-button type="danger" @click="logout" icon="el-icon-lock">{{$t('logout')}}</el-button>
        </div>
      </div>
    </el-popover>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import {clearAllCookie} from '../../utils/axios.js'

export default {
  props: {
  },
  data() {
    return {
      visible: false,
    }
  },
  computed: {
    ...mapGetters({
      userInfo: 'user/userInfo',
    }),
    displayName () {
      return 'admin'
    }
  },
  components: {},
  mounted () {
  },
  methods: {
    logout () {
      clearAllCookie()
      this.$router.push('/login')
    },
    changeActAs() {
      location.reload()
    }
  }
}
</script>

<style scoped lang="css">
.user {
  display: inline-block;
  user-select: none;
  cursor: pointer;
  transition: 0.2s;
  margin-right: 15px;
  vertical-align: top;
  color: var(--el-color-primary);
}
.user:hover {
  color: var(--el-color-primary-light-3);
}
.user .img {
  font-size: 22px;
  font-weight: 500;
  width: 30px;
  height: 30px;
  display: inline-block;
  background-color: rgba(var(--el-color-primary), 0.8);
  color: var(--el-color-white) !important;
  vertical-align: middle;
  text-align: center;
  line-height: 30px;
  border-radius: 15px;
  margin-right: 15px;
  margin-left: 10px;
  font-weight: 600;
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
  color: var(--el-color-primary);
}
.item {
  line-height: 1;
  height: 32px;
  vertical-align: top;
}
.label {
  color: var(--el-text-color-secondary);
  width: 120px;
  display: inline-block;
  vertical-align: top;
  font-size: 13px;
  padding: 7px 0;
}
.value {
  color: var(--el-text-color-primary);
  width: 200px;
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
.font-weight {
  font-weight: 600;
}
.clusterName {
  color: var(--el-color-primary) !important;
}
</style>
