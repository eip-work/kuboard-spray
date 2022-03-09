<i18n>
en:
  sync_params: Sync params to container engine
  sync_params_desc: Sync params to container engine.
  sync_params_desc_2: Will result in outage of service for several minutes.
  desc_label: Task description
zh:
  sync_params: 更新容器引擎参数
  sync_params_desc: 将更新所有节点的容器引擎参数，并重启容器引擎。
  sync_params_desc_2: 将会导致服务短暂不可用，请认真考虑后再执行操作。
  desc_label: 任务描述
</i18n>

<template>
  <ClusterTask v-if="!cluster.history.processing" :history="cluster.history" :title="title" action="sync_container_engine_params" :cluster="cluster"
    placement="top-end" @refresh="$emit('refresh')" :populateRequest="populateRequest" :width="600">
    <div style="width: 450px;">
      <el-form-item :label="$t('desc_label')">
        <li style="font-weight: bold;">{{ $t('sync_params_desc') }}</li>
        <li style="font-weight: bold;">{{ $t('sync_params_desc_2') }}</li>
      </el-form-item>
    </div>
  </ClusterTask>
</template>

<script>
import ClusterTask from '../../../common/task/ClusterTask.vue'

export default {
  props: {
    cluster: { type: Object, required: true },
  },
  data() {
    return {
    }
  },
  computed: {
    title () {
      return this.$t('sync_params')
    },
  },
  components: { ClusterTask },
  mounted () {
  },
  methods: {
    async populateRequest () {
      return {}
    }
  }
}
</script>

<style scoped lang="css">
</style>
