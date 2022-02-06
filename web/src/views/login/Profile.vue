<i18n>
en:
  changePassword: Change Password
  username: Username
  password: Password
  password2: Repeat Password
  password_placeholer: Please input new password
  password2_placeholder: Please repeat the new password
  passwordNotMatch: Password doesnot match.
zh:
  changePassword: 修改密码
  username: 用户名
  password: 密 码
  password2: 重复密码
  password_placeholer: 请输入新密码
  password2_placeholder: 请再次输入新密码
  passwordNotMatch: 两次输入的密码不匹配
</i18n>

<template>
  <div>
    <div class="app_block_title">{{ $t('changePassword') }}</div>
    <div class="block">
      <el-form ref="form" :model="form" label-position="left" label-width="120px" v-if="userInfo" style="width: 420px;">
        <el-form-item :label="$t('username')">
          <span>{{userInfo.username}}</span>
        </el-form-item>
        <el-form-item :label="$t('password')" prop="password" :rules="passwordRules">
          <el-input v-model.trim="form.password" show-password :placeholder="$t('password_placeholer')"></el-input>
        </el-form-item>
        <el-form-item :label="$t('password2')" prop="password2" :rules="passwordRules">
          <el-input v-model.trim="form.password2" show-password :placeholder="$t('password2_placeholder')"></el-input>
        </el-form-item>
        <div style="text-align: right;">
          <el-button type="primary" @click="save">{{$t('changePassword')}}</el-button>
        </div>
      </el-form>
    </div>
  </div>
</template>

<script>
import mixin from '../../mixins/mixin.js'

export default {
  mixins: [mixin],
  percentage () {
    return this.userInfo ? 100 : 0
  },
  breadcrumb () {
    return [ { label: this.$t('changePassword') } ]
  },
  refresh () {
    this.refresh()
  },
  props: {
  },
  data() {
    return {
      userInfo: undefined,
      form: {
        password: undefined,
        password2: undefined,
      },
      passwordRules: [
        { required: true, message: 'cannot be empty', trigger: 'blur' },
        {
          validator: (rule, value, callback) => {
            console.log(rule)
            if (this.form.password === this.form.password2) {
              return callback()
            } else if (rule.field === 'password2') {
              return callback(new Error(this.$t('passwordNotMatch')))
            } else {
              return callback()
            }
          },
          trigger: 'blur',
        }
      ]
    }
  },
  computed: {
  },
  components: { },
  mounted () {
    this.refresh()
  },
  methods: {
    refresh () {
      this.kuboardSprayApi.get(`/profile`).then(resp => {
        this.userInfo = resp.data.data
      }).catch(e => {
        console.log(e)
      })
    },
    save () {
      this.$refs.form.validate(flag => {
        if (flag) {
          this.kuboardSprayApi.post(`/profile/change_password`, this.form).then(resp => {
            console.log(resp.data)
            this.$message.success(this.$t('msg.save_succeeded'))
            this.form.password = undefined
            this.form.password2 = undefined
          }).catch(e => {
            console.log(e)
            this.$message.error(this.$t('msg.save_failed'))
          })
        }
      })
    }
  }
}
</script>

<style scoped lang="scss">
.block {
  background-color: var(--el-color-white);
  padding: 20px;
  border-radius: 10px;
}
</style>
