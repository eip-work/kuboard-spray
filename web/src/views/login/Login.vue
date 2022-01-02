<i18n>
en:
  hi: kuboard
zh:
  hi: kuboard
</i18n>

<template>
  <div style="height: 100%;">
    <div class="left-pannel">
      <div class="title">
        Kuboard Spray
      </div>
      <div class="app_block_title" style="margin-left: 5vw; font-size: 20px;">
        离线安装 / 维护 Kubernetes 集群
      </div>
      <div style="display: flex;">
        <img src="./login.png"/>
        <div class="logintext">登 录</div>
      </div>
    </div>

    <div class="right-panel">
      <el-card style="width: 500px; margin: auto; padding: 30px 20px;">
        <h2 class="form-header">
          <span>Kuboard Spray</span>
        </h2>
        <el-form ref="form" :model="form" label-position="left" label-width="80px" style="width: 360px; margin: auto;" size="large">
          <el-form-item prop="username" :rules="[{ required: true, message: '请输入用户名', trigger: 'blur' }]" ref="usernameFormItem">
            <template #label>
              <span style="font-weight: bolder">用户名</span>
            </template>
            <el-input v-model.trim="form.username" placeholder="请输入用户名" @keyup.enter="$refs.passwordInput.focus()" autofocus disabled></el-input>
          </el-form-item>
          <el-form-item prop="password" :rules="[{ required: true, message: '请输入密码', trigger: 'blur' }]" ref="passwordFormItem">
            <template #label>
              <span style="font-weight: bolder">密 码</span>
            </template>
            <el-input ref="passwordInput" v-model="form.password" placeholder="请输入密码" autofocus show-password @keyup.enter="login"></el-input>
          </el-form-item>
          <div style="height: 40px;">
            <!-- <el-checkbox v-if="!mfaPolicy.forceSession" style="float: left;" v-model="rememberLogin" @change="changeRememberLogin">7天内保持登录</el-checkbox> -->
            <!-- <span style="float: right;" v-if="brandConfig && !brandConfig.disableHelpLink">
              <a style="text-decoration: none; font-size: 14px" href="https://kuboard.cn/install/v3/reset-password.html" target="_blank">
                <i class="el-icon-link" style="margin-right: 8px;"></i>重置密码</a>
            </span> -->
          </div>
          <div>
            <el-button style="width: 100%;" size="large" type="primary" icon="el-icon-s-promotion" @click="login">登 录</el-button>
          </div>
        </el-form>
      </el-card>

    </div>
  </div>
</template>

<script>
import {setupCookie} from '../../utils/axios.js'

export default {
  props: {
  },
  data() {
    return {
      form: {
        username: 'admin',
        password: '',
      }
    }
  },
  computed: {
  },
  components: { },
  mounted () {
  },
  methods: {
    login() {
      console.log(this.form)
      this.kuboardSprayApi.post('/login', this.form).then(resp => {
        console.log(resp)
        setupCookie(resp.data.data.token, 7)
        this.$router.replace('/')
      }).catch(e => {
        console.log(e.response)
      })
    }
  }
}
</script>

<style scoped lang="scss">

.form-header {
  font-size: 24px;
  font-weight: bolder;
  padding-bottom: 10px;
  margin-top: 0;
  text-align: center;
}

.left-pannel {
  display: inline-block;
  width: 50vw;
  height: 100%;
  background: #F3F7FD;
  text-align: left;
  overflow: hidden;
}

.left-pannel .title {
  font-size: 42px;
  font-weight: bold;
  color: #234883;
  line-height: 50px;
  text-align: left;
  margin-top: 10vh;
  margin-left: 5vw;
}

.left-pannel img {
  width: 38vw;
  margin-left: 3vw;
  margin-top: 25vh;
}

.left-pannel .logintext {
  height: 8vh;
  line-height: 8vh;
  display: inline-block;
  vertical-align: top;
  margin-top: 32vh;
  margin-left: 2vw;
  width: 7vw;
  color: #234883;
  font-size: 1.2vw;
  font-weight: bolder;
  text-align: center;
  padding-left: 1vw;
  background-color: white;
  border-top-left-radius: 8vh;
  border-bottom-left-radius: 8vh;
}

.right-panel {
  display: inline-block;
  width: calc(50vw - 10px);
  background-color: #fff;
  color: #303133;
  vertical-align: top;
  margin-top: calc(50vh - 200px);
  /* box-shadow: 0 2px 12px 0 rgba(0,0,0,.1); */
  /* border: 1px solid #EBEEF5; */
  /* border-radius: 10px; */
}

</style>
