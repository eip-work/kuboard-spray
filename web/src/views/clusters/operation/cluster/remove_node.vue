<i18n>
en:
  removeOneEtcdNodeOnce: If you want to remove etcd node, can only remove them one by one.

  resetNodes: Reset nodes(non-applicable to offline nodes)
  resetNodesNo: Do not reset nodes
  allow_ungraceful_removal: Allow remove nodes failed to drain
  allow_ungraceful_removal_desc: 'kuboard-spray does "kubectl drain" before removing nodes, continue on failure?'
  drain_timeout: drain_timeout
  drain_retries: drain_retries
  drain_grace_period: drain_grace_period
  drain_retry_delay_seconds: drain_retry_delay_seconds
  nodes_to_remove_required: Please select nodes to remove.
zh:
  removeOneEtcdNodeOnce: 如果要删除 etcd 节点，一次只能选择一个节点删除。

  resetNodes: 重置节点（不能选择已停机的节点）
  resetNodesNo: 不重置节点
  allow_ungraceful_removal: 允许删除未排空的节点
  allow_ungraceful_removal_desc: 删除节点时，kuboard-spray 会自动执行如下指令以排空节点，如果失败，是否继续删除节点？
  drain_timeout: 排空节点超时时间
  drain_retries: 排空节点重试次数
  drain_grace_period: 应用停止等候时间
  drain_retry_delay_seconds: 两次排空尝试间隔
  nodes_to_remove_required: 请选择要删除的节点
</i18n>

<template>
  <el-form-item prop="action" :rules="rules" v-if="form.remove_node">
    <div style="margin-left: 0px; margin-top: 10px;" class="app_form_mini">
      <el-form-item prop="remove_node.reset_nodes" required>
        <el-radio-group v-model="reset_nodes">
          <el-radio-button :label="true">{{ $t('resetNodes') }}</el-radio-button>
          <el-radio-button :label="false">{{ $t('resetNodesNo') }}</el-radio-button>
        </el-radio-group>
      </el-form-item>
      <el-form-item prop="remove_node.nodes_to_remove" :rules="nodes_to_remove_rules">
        <el-checkbox-group v-model="formRef.remove_node.nodes_to_remove">
          <el-checkbox v-for="(item, key) in pendingRemoveNodes" :key="'h' + key" style="margin-top: 10px; margin-right: 10px;" :label="item.name"
            :disabled="(pingpong[item.name] === undefined || pingpong[item.name].status !== 'SUCCESS') && formRef.remove_node.reset_nodes">
            <el-tooltip v-if="pingpong[item.name] && pingpong[item.name].status !== 'SUCCESS'" :content="pingpong[item.name].message" placement="top-start">
              <span class="app_text_mono">{{item.name}}</span>
            </el-tooltip>
            <span v-else class="app_text_mono">{{item.name}}</span>
          </el-checkbox>
        </el-checkbox-group>
      </el-form-item>
      <el-form-item label-width="150px" :label="$t('allow_ungraceful_removal')" style="margin-top: 10px;">
        <el-switch v-model="formRef.remove_node.allow_ungraceful_removal"></el-switch>
        <div style="line-height: 20px; display: inline-block; vertical-align: top; margin-left: 10px;">
          <div class="form_description">{{$t('allow_ungraceful_removal_desc')}}</div>
          <div class="form_description app_text_mono" style="color: #666; font-weight: bold;">
            kubectl drain --force --ignore-daemonsets --grace-period {{ formRef.remove_node.drain_grace_period * 60}} 
              --timeout {{ formRef.remove_node.drain_timeout * 60 }}s
          </div>
        </div>
      </el-form-item>
      <el-form-item label-width="150px" :label="$t('drain_grace_period')" required prop="remove_node.drain_grace_period">
        <el-input-number v-model="formRef.remove_node.drain_grace_period" :min="1" :max="formRef.remove_node.drain_timeout - 1"></el-input-number> 分钟
        <span style="margin-left: 20px;" class="form_description">kubectl drain --grace-period</span>
      </el-form-item>
      <el-form-item label-width="150px" :label="$t('drain_timeout')" required prop="remove_node.drain_timeout">
        <el-input-number v-model="formRef.remove_node.drain_timeout" :min="formRef.remove_node.drain_grace_period + 1"></el-input-number> 分钟
        <span style="margin-left: 20px;" class="form_description">kubectl drain --timeout</span>
      </el-form-item>
      <el-form-item label-width="150px" :label="$t('drain_retries')" required prop="remove_node.drain_retries">
        <el-input-number v-model="formRef.remove_node.drain_retries" :min="1"></el-input-number> 次
      </el-form-item>
      <el-form-item label-width="150px" :label="$t('drain_retry_delay_seconds')" required prop="remove_node.drain_retry_delay_seconds">
        <el-input-number v-model="formRef.remove_node.drain_retry_delay_seconds" :min="5" :step="5"></el-input-number> 秒
      </el-form-item>
    </div>
  </el-form-item>
</template>

<script>
import { trimMark } from './utils'

export default {
  props: {
    cluster: { type: Object, required: true },
    name: { type: String, required: true },
    form: { type: Object, required: true },
    pingpong: { type: Object, required: false, default: () => { return {} } },
    pingpong_loading: { type: Boolean, required: false, default: false },
    offlineNodes: { type: Array, required: false, default: () => { return [] } },
  },
  data() {
    let _this = this
    return {
      nodes_to_remove_rules: [
        { required: true, message: _this.$t('nodes_to_remove_required') }
      ],
      rules: [{
        validator (rule, value, callback) {
          for (let node of _this.formRef.remove_node.nodes_to_remove) {
            if (_this.cluster.inventory.all.children.target.children.etcd.hosts[node] !== undefined && _this.formRef.remove_node.nodes_to_remove.length > 1) {
              return callback(_this.$t('removeOneEtcdNodeOnce'))
            }
          }
          return callback()
        },
        trigger: 'change',
      }]
    }
  },
  inject: ['isClusterInstalled', 'isClusterOnline', 'pendingRemoveNodes', 'pendingAddNodes'],
  computed: {
    formRef: {
      get () { return this.form },
      set () {},
    },
    reset_nodes: {
      get () { return this.form.remove_node ? this.form.remove_node.reset_nodes : true },
      set (v) {
        let temp = []
        for (let node of this.pendingRemoveNodes) {
          let online = false
          if (this.pingpong[node.name]) {
            online = this.pingpong[node.name].status === 'SUCCESS'
          }
          if ((online) === v) {
            temp.push(node.name)
          }
        }
        this.formRef.remove_node.nodes_to_remove = temp
        this.formRef.remove_node.reset_nodes = v
      }
    },
  },
  components: { },
  mounted () {
    this.updateForm()
  },
  methods: {
    updateForm() {
      this.formRef.remove_node = {}
      this.formRef.remove_node.nodes_to_remove = []
      this.formRef.remove_node.reset_nodes = true
      this.formRef.remove_node.allow_ungraceful_removal = false
      this.formRef.remove_node.drain_grace_period = 5
      this.formRef.remove_node.drain_timeout = 6
      this.formRef.remove_node.drain_retries = 3
      this.formRef.remove_node.drain_retry_delay_seconds = 10
    },
    populateRequest() {
      let temp = ''
      for (let node of this.form.remove_node.nodes_to_remove) {
          temp += node + ','
      }
      temp = trimMark(temp)
      let req = {}
      req.nodes_to_remove = temp
      req.reset_nodes = this.form.remove_node.reset_nodes
      req.allow_ungraceful_removal = this.form.remove_node.allow_ungraceful_removal
      req.drain_grace_period = this.form.remove_node.drain_grace_period * 60 + ''
      req.drain_timeout = this.form.remove_node.drain_timeout * 60 + 's'
      req.drain_retries = this.form.remove_node.drain_retries + ''
      req.drain_retry_delay_seconds = this.form.remove_node.drain_retry_delay_seconds + ''
      return req
    },
  }
}
</script>

<style scoped lang="scss">
.form_description {
  font-size: 12px;
  // color: var(--el-text-color-placeholder);
  color: #aaa;
  max-width: 700px;
}
</style>
