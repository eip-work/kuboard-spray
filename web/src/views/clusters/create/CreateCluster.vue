<i18n>
en:
  addCluster: Add Cluster Installation Plan
  createResource: "Add Resource Package"
  name: Name
  requiresName: Please input name.
  conflict: conflict with existing one.
  goToResourcePage: This is going to open Resource package management page in a new window, do you confirm?
zh:
  addCluster: 添加集群安装计划
  createResource: '添加资源包'
  name: 名称
  requiresName: 请填写名称
  conflict: 与已有的重复 {name}
  goToResourcePage: 此操作将要在新窗口打开资源包加载页面，完成设置后，您可以切换窗口回到当前页面，是否继续？
</i18n>


<template>
  <div>
    <el-dialog v-model="dialogVisible" :close-on-click-modal="false" :modal="true" top="20vh"
      :title="$t('addCluster')" width="45%">
      <el-form :model="form" label-position="left" label-width="120px" v-if="dialogVisible" ref="form">
        <FieldString :holder="form" fieldName="cluster_name" required :placeholder="$t('requiresName')" :rules="nameRules"></FieldString>
        <FieldSelect :holder="form" fieldName="kuboardspray_resource_package" :loadOptions="loadResourceList" required>
          <template #edit>
            <ConfirmButton buttonStyle="margin-left: 10px;" icon="el-icon-plus"
              @confirm="openUrlInBlank('/#/settings/resources')"
              :text="$t('createResource')" :message="$t('goToResourcePage')" width="420"></ConfirmButton>
          </template>
        </FieldSelect>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false" icon="el-icon-close" type="info" plain>{{$t('msg.cancel')}}</el-button>
          <el-button @click="save" icon="el-icon-check" type="primary" :loading="saving">{{$t('msg.ok')}}</el-button>
        </span>
      </template>
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
        create: {
          cluster_name: '',
          kuboardspray_resource_package: '',
        }
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
              return callback('必须以字母开头，可以包含数字和字母')
            }
            this.kuboardSprayApi.get(`/clusters/${value}`).then(() => {
              callback(this.$t('conflict', {name: value}))
            }).catch(e => {
              // console.log(e.response)
              if (e.response && e.response.data.code === 500) {
                callback()
              }
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
  watch: {
  },
  methods: {
    show () {
      this.dialogVisible = true
    },
    async loadResourceList () {
      let result = []
      await this.kuboardSprayApi.get('/resources').then(resp => {
        for (let res of resp.data.data) {
          result.push({ label: res, value: res })
        }
      }).catch(e => {
        console.log(e)
      })
      return result
    },
    save () {
      this.$refs.form.validate(async flag => {
        if (flag) {
          this.saving = true
          let req = {
            name: this.form.cluster_name,
            resource_package: this.form.kuboardspray_resource_package,
          }
          await this.kuboardSprayApi.post('/clusters', req).then(resp => {
            console.log(resp.data.data)
            this.$message.success(this.$t('msg.save_succeeded'))
            this.$router.push(`/clusters/${this.form.cluster_name}?mode=edit`)
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

<style scoped lang="css">

</style>
