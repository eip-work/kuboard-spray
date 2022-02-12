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

  os_type: OS Type
  mirror_type: Source Type
  mirror_type_os: OS source
  mirror_type_os_desc: OS Source
  mirror_type_docker: Docker source
  mirror_type_docker_desc: contains installation packages for docker-ce / docker-cli / containerd.
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

  os_type: 操作系统
  mirror_type: 软件源类型
  mirror_type_os: 操作系统软件源
  mirror_type_os_desc: 操作系统提供的软件源
  mirror_type_docker: docker 软件源
  mirror_type_docker_desc: 安装 docker-ce / docker-cli / containerd 等软件包的源
</i18n>


<template>
  <div>
    <el-dialog v-model="dialogVisible" :close-on-click-modal="false" :modal="true" top="20vh"
      :title="$t('addMirror')" width="45%">
      <el-form :model="form" label-position="left" label-width="120px" v-if="dialogVisible" ref="form">
        <el-form-item :label="$t('mirror_type')">
          <el-radio-group v-model="mirror_type">
            <el-radio-button label="os">{{ $t('mirror_type_os') }}</el-radio-button>
            <el-radio-button label="docker">{{ $t('mirror_type_docker') }}</el-radio-button>
          </el-radio-group>
          <span style="margin-left: 20px; vertical-align: bottom; color: var(--el-text-color-secondary);">{{$t('mirror_type_' + mirror_type + '_desc')}}</span>
        </el-form-item>
        <el-form-item :label="$t('os_type')" prop="kuboardspray_os_mirror_type" :rules="[{ required: true, message: 'required', trigger: 'change' }]">
          <el-radio-group v-model="os_type">
            <el-radio-button label="anolis">Anolis</el-radio-button>
            <el-radio-button label="centos">Centos</el-radio-button>
            <el-radio-button label="openeuler">openEuler</el-radio-button>
            <el-radio-button label="oraclelinux">OracleLinux</el-radio-button>
            <el-radio-button label="redhat">Redhat</el-radio-button>
            <el-radio-button label="rocky">Rocky</el-radio-button>
            <el-radio-button label="ubuntu">Ubuntu</el-radio-button>
          </el-radio-group>
        </el-form-item>
        <FieldString :holder="form" fieldName="kuboardspray_os_mirror_name" required :rules="nameRules" :disabled="!form.kuboardspray_os_mirror_type"></FieldString>
        <FieldRadio :holder="form" fieldName="kuboardspray_os_mirror_kind" required :options="['existing']">
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
import mirror_options from './mirror_options.js'

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
    mirror_type: {
      get () {
        if (this.form.kuboardspray_os_mirror_type && this.form.kuboardspray_os_mirror_type.indexOf('docker_') === 0) {
          return 'docker'
        } else {
          return 'os'
        }
      },
      set (v) {
        if (v === 'docker') {
          this.form.kuboardspray_os_mirror_type = 'docker_' + this.os_type
        } else {
          this.form.kuboardspray_os_mirror_type = this.os_type
        }
        this.form.kuboardspray_os_mirror_url = undefined
      }
    },
    os_type: {
      get () {
        let os_mirror_type = this.form.kuboardspray_os_mirror_type
        if ( os_mirror_type && os_mirror_type.indexOf('docker_') === 0) {
          return os_mirror_type.split('_')[1]
        } else {
          return os_mirror_type.split('_')[0]
        }
      },
      set (v) {
        if (this.mirror_type === 'docker') {
          this.form.kuboardspray_os_mirror_type = 'docker_' + v
        } else {
          this.form.kuboardspray_os_mirror_type = v
        }
        this.form.kuboardspray_os_mirror_url = undefined
      }
    }
  },
  components: { },
  mounted () {
  },
  methods: {
    show () {
      this.dialogVisible = true
    },
    async loadMirrorList () {
      let result = []
      for (let item of mirror_options[this.form.kuboardspray_os_mirror_type]) {
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
