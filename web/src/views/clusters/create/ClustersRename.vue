<i18n>
en:
  conflict: conflict with existing one {name}.
  rename: Rename
  renameCluster: Rename cluster {name} to
  input_newName: Please input new cluster name here.
zh:
  conflict: 与已有的重复 {name}
  rename: 重命名
  renameCluster: 将集群 {name} 重命名为
  input_newName: 请输入集群的新名称
</i18n>

<template>
  <el-popover trigger="manual" v-model:visible="show" placement="top-start" :title="$t('rename')" width="420px">
    <template #reference>
      <el-button icon="el-icon-edit" circle @click="show = true"></el-button>
    </template>
    <el-form :model="form" ref="form">
      <div class="app_text_mono" style="text-align: left; line-height: 28px;">
        <span>{{ $t('renameCluster', { name: clusterName}) }}</span>
        <el-form-item prop="name" :rules="nameRules">
          <el-input v-model.trim="form.name" :placeholder="$t('input_newName')"></el-input>
        </el-form-item>
      </div>
      <div style="text-align: right;">
        <el-button icon="el-icon-close" type="default" @click="show = false">{{ $t('msg.cancel') }}</el-button>
        <el-button icon="el-icon-check" type="primary" @click="renameCluster">{{ $t('msg.ok') }}</el-button>
      </div>
    </el-form>
  </el-popover>
</template>

<script>
export default {
  props: {
    clusterName: { type: String, required: true }
  },
  data() {
    return {
      show: false,
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
      ],
      form: {
        name: undefined,
      }
    }
  },
  computed: {
  },
  components: { },
  mounted () {
  },
  methods: {
    renameCluster() {
      this.$refs.form.validate(flag => {
        if (flag) {
          this.kuboardSprayApi.patch(`/clusters/${this.clusterName}?newName=${this.form.name}`).then(resp => {
            console.log(resp.data)
            this.show = false
            this.$message.success(this.$t('msg.save_succeeded'))
            this.$emit('success')
          }).catch(e => {
            console.error(e)
            this.$message.error('修改集群名称失败 : ' + e)
          })
        }
      })
    }
  }
}
</script>

<style scoped lang="scss">

</style>
