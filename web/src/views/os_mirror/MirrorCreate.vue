<i18n>
en:
  addMirror: Add OS Mirror
  name: Name
  conflict: conflict with existing one.
  provision_msg: Requires one server, and at least 200G disk space.
  existing_msg: If you already have a existing OS mirror source that is available to all the machines used to install K8S.
  canCreateItem: You  can input a mirror url that does not exist in the list.
  ubuntu: Ubuntu apt mirror
  centos: Centos yum repo
  docker_ubuntu: Ubuntu apt mirror for docker
  docker_centos: Centos yum repo for docker
zh:
  addMirror: 添加 OS 软件源
  name: 名称
  conflict: 与已有的重复 {name}
  provision_msg: 需要您提供一台机器，包含至少 200G 磁盘空间
  existing_msg: 如果您已经有一个安装集群所用机器都可以访问到的 OS 软件源，可以直接在此处输入该软件源的访问地址。
  canCreateItem: 您可以在此直接输入下拉列表中没有的条目
  ubuntu: Ubuntu apt mirror
  centos: Centos yum repo
  docker_ubuntu: Docker 的 Ubuntu apt mirror
  docker_centos: Docker 的 Centos yum repo
</i18n>


<template>
  <div>
    <el-dialog v-model="dialogVisible" :close-on-click-modal="false" :modal="true" top="20vh"
      :title="$t('addMirror')" width="45%">
      <el-form :model="form" label-position="left" label-width="120px" v-if="dialogVisible" ref="form">
        <FieldSelect :holder="form" fieldName="kuboardspray_os_mirror_type" :loadOptions="loadMirrorTypeList" required>
        </FieldSelect>
        <FieldString :holder="form" fieldName="kuboardspray_os_mirror_name" required :rules="nameRules" :disabled="!form.kuboardspray_os_mirror_type"></FieldString>
        <FieldRadio :holder="form" fieldName="kuboardspray_os_mirror_kind" required :options="['existing', 'provision']">
          <template #edit>
            <div class="create_kind">
              <div v-if="form.kuboardspray_os_mirror_kind === 'provision'">
                <el-alert :closable="false" type="warning" effect="dark" :title="$t('msg.prompt')">{{$t('provision_msg')}}</el-alert>
              </div>
              <div v-else>
                <FieldSelect :holder="form" fieldName="kuboardspray_os_mirror_url" required :disabled="!form.kuboardspray_os_mirror_type"
                  allow-create filterable :loadOptions="loadMirrorList" :placeholder="$t('canCreateItem')"></FieldSelect>
                <el-alert :closable="false" type="warning" :title="$t('msg.prompt')">{{$t('existing_msg')}}</el-alert>
              </div>
            </div>
          </template>
        </FieldRadio>
      </el-form>
      <div class="dialog-footer">
        <el-button @click="dialogVisible = false" icon="el-icon-close" type="info" plain>{{$t('msg.cancel')}}</el-button>
        <el-button @click="save" icon="el-icon-check" type="primary" :loading="saving">{{$t('msg.ok')}}</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import {ref} from 'vue'
import clone from 'clone'
import {syncParams} from './params/sync_params.js'

export default {
  setup () {
    return {
      dialogVisible: ref(false),
    }
  },
  props: {
  },
  provide () {
    return {
      editMode: 'edit',
    }
  },
  data() {
    return {
      saving: false,
      form: {
        kuboardspray_os_mirror_name: '',
        kuboardspray_os_mirror_type: '',
        kuboardspray_os_mirror_kind: 'existing',
        kuboardspray_os_mirror_url: '',
      },
      nameRules: [
        {
          validator: (rule, value, callback) => {
            if (value.length < 4) {
              return callback('min length ' + 4)
            }
            if (value.length > 20) {
              return callback('max length ' + 20)
            }
            if (!/^[a-zA-Z][a-zA-Z0-9_]{3,21}$/.test(value)) {
              return callback('必须以字母开头，可以包含数字，字母，下划线')
            }
            let v = this.form.kuboardspray_os_mirror_type + '-' + this.form.kuboardspray_os_mirror_name
            this.kuboardSprayApi.get(`/mirrors/${v}`).then(() => {
              callback(this.$t('conflict', {name: v}))
            }).catch(e => {
              if (e.response && e.response.data.code === 500) {
                callback()
                return
              }
              callback()
            })
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
  },
  methods: {
    show () {
      this.dialogVisible = true
    },
    async loadMirrorTypeList () {
      let result = [
        { label: this.$t('ubuntu'), value: 'ubuntu' },
        { label: this.$t('centos'), value: 'centos' },
        { label: this.$t('docker_ubuntu'), value: 'docker_ubuntu' },
        { label: this.$t('docker_centos'), value: 'docker_centos' },
      ]
      return result
    },
    async loadMirrorList () {
      let temp = {
        ubuntu: [
            'https://repo.huaweicloud.com/ubuntu/',
            'https://mirrors.aliyun.com/ubuntu/',
            'https://mirrors.tuna.tsinghua.edu.cn/ubuntu/',
            'http://cn.archive.ubuntu.com/ubuntu/',
            'https://mirrors.cloud.tencent.com/ubuntu/',
            'http://ftp.sjtu.edu.cn/ubuntu',
            'http://mirrors.163.com/ubuntu/',
            'http://mirrors.nju.edu.cn/ubuntu/',
          ],
        docker_ubuntu: [
            'https://repo.huaweicloud.com/docker-ce/linux/ubuntu/',
            'https://mirrors.tuna.tsinghua.edu.cn/docker-ce/linux/ubuntu/',
            'https://mirrors.aliyun.com/docker-ce/linux/ubuntu/',
            'https://mirrors.cloud.tencent.com/docker-ce/linux/ubuntu/',
          ],
        centos: [
          'http://mirrors.aliyun.com/repo/',
          'https://repo.huaweicloud.com/centos/',
          'https://mirrors.tuna.tsinghua.edu.cn/centos/',
          'https://mirrors.cloud.tencent.com/centos/',
        ],
        docker_centos: [
          'https://download.docker.com/linux/centos/',
        ],
      }
      let result = []
      for (let item of temp[this.form.kuboardspray_os_mirror_type]) {
        result.push({ label: item, value: item })
      }
      return result
    },
    save () {
      this.$refs.form.validate(async flag => {
        if (flag) {
          this.saving = true
          
          let temp = clone(this.form)
          temp.params = {}
          
          syncParams(temp.params, this.form.kuboardspray_os_mirror_type, this.form.kuboardspray_os_mirror_url)

          await this.kuboardSprayApi.post('/mirrors', temp).then(() => {
            this.$message.success(this.$t('msg.save_succeeded'))
            if (this.form.kuboardspray_os_mirror_kind === 'existing') {
              this.$router.push(`/settings/mirrors/${this.form.kuboardspray_os_mirror_type}-${this.form.kuboardspray_os_mirror_name}`)
            } else {
              this.$router.push(`/settings/mirrors/${this.form.kuboardspray_os_mirror_type}-${this.form.kuboardspray_os_mirror_name}?mode=edit`)
            }
          }).catch(e => {
            this.$message.error(this.$t('msg.save_failed', {msg: e.response.data.message}))
          })
          this.saving = false
        }
      })
    },

  }
}
</script>

<style scoped lang="css">
.create_kind {
  padding: 10px;
  background-color: var(--el-color-primary-light-9);
  margin-top: 10px;
}
.dialog-footer {
  padding: 20px 0;
  text-align: right;
}
</style>
