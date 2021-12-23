<i18n>
en:
  manageSshKey: Manage Private Keys
  addSshKey: Add Private Key
  existingSshKeyFiles: Existing SSH private key files to {ownerType}/{ownerName}
  noKeyFiles: No existing SSH private keyfiles
  copyFile: If private key file is on another machine, suggest you use scp or other tool to transfer it to the current client machine, then upload it here. If you copy content into a editor and save it as a file, it may causes format error.
  nameRequired: Please input private key name
  duplicatedPrivateKeyFile: Conflict with a existing one
  selectFile: Select
  deleteSshKey: You are about to delete SSH private key {name}, do you confirm ?
zh:
  manageSshKey: 管理私钥
  addSshKey: 添加私钥
  existingSshKeyFiles: 已有的私钥文件（{ownerType}/{ownerName}）
  noKeyFiles: 暂时没有可用的私钥文件
  copyFile: 如果私钥文件在另外一台机器，建议用 scp 或其他工具将文件传送到当前浏览器所在机器再上传，直接拷贝内容再通过编辑器另存为一个文件，可能导致格式错误。
  nameRequired: 请输入 private key 的名称
  duplicatedPrivateKeyFile: 已经存在同名的 private key 文件
  selectFile: 选择文件
  deleteSshKey: 此操作将要删除私钥 {name}，是否继续？
</i18n>

<template>
  <div>
    <el-dialog v-model="dialogVisible" :close-on-click-modal="false" :modal="true"
      :title="$t('manageSshKey')" width="60%">
      <div class="app_block_title_small">
        {{$t('existingSshKeyFiles', {ownerType, ownerName})}}
        <el-button circle type="primary" icon="el-icon-refresh" @click="load"></el-button>
      </div>
      <el-skeleton :rows="1" animated :loading="loading">
        <el-alert v-if="keyFiles.length === 0" :title="$t('noKeyFiles')" :closable="false" type="warning"></el-alert>
        <template v-else>
          <el-tag v-for="(item, index) in keyFiles" :key="'file' + index"
            @close="deleteSshKey(item)"
            size="medium" style="margin-top: 10px; margin-right: 10px;" closable>
            {{item}}
          </el-tag>
        </template>
      </el-skeleton>

      <div class="app_margin_top uploadSection">
        <div class="app_block_title_small">
          {{$t('addSshKey')}}
        </div>
        <el-upload ref="upload" :action="`${baseURL}/private-keys/${ownerType}/${ownerName}/${form.name}`"
          :headers="headers"
          :auto-upload="false">
          <template #trigger>
            <div class="step">1</div>
            <el-button size="small" type="primary" icon="el-icon-files">{{$t('selectFile')}}</el-button>
          </template>
          <div style="display: inline-block;">
            <el-form :model="form" class="app_form_mini" ref="form">
              <div class="step" style="margin-left: 20px">2</div>
              <el-form-item prop="name" :rules="rules" style="display: inline-block;">
                <el-input style="width: 300px;" v-model.trim="form.name" :placeholder="$t('nameRequired')"></el-input>
              </el-form-item>
            </el-form>
          </div>
          <div class="step" style="margin-left: 20px">3</div>
          <el-button size="small" type="success" :loading="uploading" @click="submitUpload" icon="el-icon-upload2">{{$t('msg.upload')}}</el-button>
          <template #tip>
            <div class="el-upload__tip">
              {{$t('copyFile')}}
            </div>
          </template>
        </el-upload>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false" icon="el-icon-close" type="warning" plain>{{$t('msg.close')}}</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import {baseURL} from '../../utils/axios.js'
import Cookie from 'js-cookie'

export default {
  props: {
    ownerType: { type: String, required: true },
    ownerName: { type: String, required: true },
  },
  data() {
    return {
      dialogVisible: false,
      keyFiles: [],
      loading: false,
      uploading: false,
      form: {
        name: undefined
      },
      rules: [
        { required: true, message: this.$t('nameRequired'), trigger: 'blur' },
        {
          validator: async (rule, value, callback) => {
            await this.kuboardSprayApi.get(`/private-keys/${this.ownerType}/${this.ownerName}/${value}`).then(() => {
              callback(this.$t('duplicatedPrivateKeyFile'))
            }).catch(e => {
              if (e.response.status === 500) {
                callback()
              }
            })
          },
          trigger: 'blur',
        }
      ],
      baseURL,
      headers: {
        Authorization: `Bearer ${Cookie.get('KuboardToken')}`
      }
    }
  },
  computed: {
  },
  components: { },
  mounted () {
  },
  methods: {
    async show () {
      this.dialogVisible = true
      this.load ()
    },
    async load () {
      this.loading = true
      await this.kuboardSprayApi.get(`/private-keys/${this.ownerType}/${this.ownerName}`).then(resp => {
        this.keyFiles = resp.data.data
      }).catch(e => {
        console.log(e)
      })
      setTimeout(() => {
        this.loading = false
      }, 400)
    },
    submitUpload () {
      this.$refs.form.validate(flag => {
        if (flag) {
          this.uploading = true
          this.$refs.upload.submit()
          setTimeout(() => {
            this.$message.success(this.$t('msg.save_succeeded'))
            this.load()
            this.uploading = false
            this.$refs.upload.clearFiles()
            this.form.name = undefined
          }, 500)
        }
      })
    },
    deleteSshKey (name) {
      this.$confirm(
        this.$t('deleteSshKey', {name}),
        'Warning',
        {
          confirmButtonText: this.$t('msg.ok'),
          cancelButtonText: this.$t('msg.cancel'),
          type: 'warning',
        }
      ).then(() => {
        this.kuboardSprayApi.delete(`/private-keys/${this.ownerType}/${this.ownerName}/${name}`).then(resp => {
          this.$message.success("Delete " + resp.data.message)
          this.load()
        }).catch(e => {
          this.$message.error("Failed " + e.response.data.message)
        })
      }).catch(() => {
        this.$message.info("canceled")
      })
    }
  }
}
</script>

<style scoped lang="scss">
.uploadSection {
  margin-top: 20px;
  border: 1px solid $--border-color-light;
  border-radius: 10px;
  padding: 20px;
}
.step {
  border: solid 1px $--color-warning;
  border-radius: 30px;
  width: 28px;
  height: 28px;
  background-color: $--color-warning;
  color: white;
  display: inline-block;
  text-align: center;
  margin: 0 10px 0px 0px;
}
</style>
