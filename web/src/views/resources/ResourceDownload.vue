<i18n>
en:
  title: "Load resource package {name}."
  selectSource: "Select a source"
  download: Download Resource Package
zh:
  title: "加载资源包 {name}"
  selectSource: "选择一个源"
  download: '加载资源包'
</i18n>

<template>
  <ExecuteTask v-if="resource" :history="resource.history" :startTask="startTask" :label="$t('download')"
    :title="$t('title', {name: resource.package.metadata.version})" :loading="loading" @refresh="$emit('refresh')">
    <el-form @submit.prevent.stop :model="form" ref="form" label-position="left" label-width="120px">
      <el-form-item :label="$t('selectSource')" prop="downloadFrom" :rules="sourceRules">
        <el-radio-group v-model="form.downloadFrom">
          <div style="line-height: 28px; padding-top: 5px;">
            <div v-for="(source, index) in resource.package.metadata.available_at" :key="'source' + index">
              <el-radio :label="source">
                <span class="app_text_mono">{{source}}</span>
              </el-radio>
            </div>
          </div>
        </el-radio-group>
      </el-form-item>
    </el-form>
  </ExecuteTask>
</template>

<script>
import clone from 'clone'
import ExecuteTask from '../common/task/ExecuteTask.vue'

export default {
  props: {
    resource: { type: Object, required: true },
    loading: { type: Boolean, required: false },
    action: { type: String, required: true },
  },
  data() {
    return {
      form: {
        downloadFrom: undefined,
      },
      sourceRules: [
        { required: true, message: this.$t('selectSource'), trigger: 'change' }
      ]
    }
  },
  components: { ExecuteTask },
  emits: ['refresh'],
  mounted () {
    this.form.downloadFrom = this.resource.package.metadata.available_at[0]
  },
  methods: {
    startTask () {
      return new Promise((resolve, reject) => {
        this.$refs.form.validate(async flag => {
          if (flag) {
            let request = {
              package: clone(this.resource.package),
              downloadFrom: this.form.downloadFrom,
            }
            this.kuboardSprayApi.post(`/resources/${request.package.metadata.version}/${this.action}`, request).then(resp => {
              this.$router.replace(`/settings/resources/${request.package.metadata.version}`)
              resolve(resp.data.data.pid)
            }).catch(e => {
              reject(e)
            })
          } else {
            resolve()
          }
        })
      })
    }
  }
}
</script>

<style scoped lang="scss">
.confirmText {
  font-size: 15px;
  color: $--color-danger;
  font-weight: bold;
}
</style>
