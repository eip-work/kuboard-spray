<i18n>
en:
  addMirror: Add OS Mirror
  name: Name
  conflict: conflict with existing one.
  provision_msg: Requires one server, and at least 200G disk space.
  existing_msg: If you already have a existing OS mirror source that is available to all the machines used to install K8S.
zh:
  addMirror: 添加 OS 软件源
  name: 名称
  conflict: 与已有的重复 {name}
  provision_msg: 需要您提供一台机器，包含至少 200G 磁盘空间
  existing_msg: 如果您已经有一个安装集群所用机器都可以访问到的 OS 软件源
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
          <div class="create_kind">
            <div v-if="form.kuboardspray_os_mirror_kind === 'provision'">
              <el-alert :closable="false" type="warning" effect="dark" :title="$t('msg.prompt')">{{$t('provision_msg')}}</el-alert>
            </div>
            <div v-else>
              <FieldString :holder="form" fieldName="kuboardspray_os_mirror_url" required></FieldString>
              <el-alert :closable="false" type="warning" :title="$t('msg.prompt')">{{$t('existing_msg')}}</el-alert>
            </div>
          </div>
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
              console.log(e.response)
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
        { label: 'ubuntu', value: 'ubuntu' },
        { label: 'centos', value: 'centos', disabled: true },
        { label: 'docker-ce_ubuntu', value: 'docker-ce_ubuntu' },
        { label: 'docker-ce_centos', value: 'docker-ce_centos', disabled: true },
      ]
      return result
    },
    save () {
      this.$refs.form.validate(async flag => {
        if (flag) {
          this.saving = true
          await this.kuboardSprayApi.post('/mirrors', this.form).then(() => {
            this.$message.success(this.$t('msg.save_succeeded'))
            this.$router.push(`/settings/mirrors/${this.form.kuboardspray_os_mirror_type}-${this.form.kuboardspray_os_mirror_name}?mode=edit`)
          }).catch(e => {
            this.$message.error(this.$t('msg.save_failed', {msg: e.response.data.message}))
          })
          this.saving = false
        }
      })
    }
  }
}
</script>

<style scoped lang="scss">
.create_kind {
  padding: 10px;
  background-color: $--color-primary-light-9;
  margin-top: 10px;
}
.dialog-footer {
  padding: 20px 0;
  text-align: right;
}
</style>
