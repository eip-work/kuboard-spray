<i18n>
en:
  startTime: Start Time
  duration: Duration
  status: Status
  viewLog: View Log
  operation: Action
  statusInHistory: Result

  killed: Killed
  processing: Processing
  success: Success
  failed: Failed
zh:
  startTime: 开始时间
  duration: 耗时
  status: 状态
  viewLog: 查看日志
  operation: 操作
  statusInHistory: 执行结果

  killed: 强制终止
  processing: 正在执行
  success: 成功
  failed: 失败
</i18n>

<template>
  <div class="operation_step_history">
    <el-table v-loading="loading" :data="history" :header-cell-style="{ 'text-align': 'center' }"
      :cell-style="{ 'text-align': 'center' }">
      <el-table-column :label="t('startTime')" min-width="150px">
        <template #default="scope">
          {{ formatTime(scope.row.time) }}
        </template>
      </el-table-column>
      <el-table-column :label="t('duration')" min-width="50px">
        <template #default="scope">
          <div class="duration">
            {{ Math.round(scope.row.duration / 1000) }}s
          </div>
        </template>
      </el-table-column>
      <el-table-column :label="t('status')">
        <template #default="scope">
          <el-tag v-if="scope.row.status == 'killed'" type="warning">{{ t(scope.row.status) }}</el-tag>
          <el-tag v-else-if="scope.row.status == 'processing'" type="primary">{{ t(scope.row.status) }}</el-tag>
          <el-tag v-else-if="scope.row.status == 'success'" type="success">{{ t(scope.row.status) }}</el-tag>
          <el-tag v-else-if="scope.row.status == 'failed'" type="danger">{{ t(scope.row.status) }}</el-tag>
          <el-tag v-else> {{ scope.row.status }} </el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="t('operation')" min-width="80px">
        <template #default="scope">
          <div class="operation">
            <el-button type="primary" link icon="el-icon-pointer" @click="viewTaskLogs(scope.row)">
              {{ t("viewLog") }}
            </el-button>
            <el-button type="primary" link icon="el-icon-notification" style="margin: 0"
              @click="viewTaskResult(scope.row)">
              {{ t("statusInHistory") }}
            </el-button>
          </div>
        </template>
      </el-table-column>
    </el-table>
    <div>
      <OperationStepStatusDialog ref="statusDialog"></OperationStepStatusDialog>
    </div>
  </div>
</template>


<script>
import dayjs from 'dayjs';
import OperationStepStatusDialog from './OperationStepStatusDialog.vue';

export default {
  props: {
    cluster: { type: Object, required: true },
    currentOperation: { type: Number, required: true },
    currentStep: { type: Number, required: true },
  },
  emits: ['loadingStatusChanged'],
  data() {
    return {
      history: [],
      loading: false
    }
  },
  computed: {
    apiPath() {
      let path = "/clusters/" + this.cluster.name;
      if (this.cluster.resourcePackage.data.operations[this.currentOperation] == undefined) {
        return "";
      }
      path += "/history/" + this.cluster.resourcePackage.data.operations[this.currentOperation]?.name;

      if (this.cluster.resourcePackage.data.operations[this.currentOperation]?.steps[this.currentStep] == undefined) {
        return "";
      }
      path += "/step/" + this.cluster.resourcePackage.data.operations[this.currentOperation]?.steps[this.currentStep]?.name;
      return path;
    }
  },
  components: { OperationStepStatusDialog },
  mounted() {
    this.refresh()
  },
  watch: {
    apiPath() {
      this.refresh()
    }
  },
  methods: {
    formatTime(time) {
      let temp = time.padEnd(23, "0")
      return dayjs(temp, "YYYY-MM-DD_HH-mm-ss.SSS").format("YYYY-MM-DD HH:mm:ss")
    },
    /**
     * 这是一个异步刷新函数。如果存在apiPath，则通过pangeeClusterApi获取数据，并将返回的数据赋值给history变量，同时在控制台打印出返回的数据。
     *
     * @async
     * @function refresh
     * @return {Promise}  无返回值，但会改变this.history的值，并在控制台打印出返回的数据。
     */
    async refresh() {
      if (this.apiPath) {
        this.loading = true;
        this.$emit('loadingStatusChanged', true);
        this.history = [];
        this.pangeeClusterApi.get(this.apiPath).then(resp => {
          this.history = resp.data.data.history
          this.loading = false;
          this.$emit('loadingStatusChanged', false);
        }).catch(e => {
          console.log(e);
          this.loading = false;
          this.$emit('loadingStatusChanged', false);
        })
      }
    },
    viewTaskLogs(item) {
      let path = `/#/tail/cluster/${this.cluster.name}/history/`
      path = `${path}${this.cluster.resourcePackage.data.operations[this.currentOperation]?.name}/`
      path = `${path}${this.cluster.resourcePackage.data.operations[this.currentOperation]?.steps[this.currentStep]?.name}/`
      path = `${path}${item.time}/execute.log`
      this.openUrlInBlank(path)
    },
    viewTaskResult(item) {
      let path = "/clusters/" + this.cluster.name + "/history/" + this.cluster.resourcePackage.data.operations[this.currentOperation].name;
      path += "/step/" + this.cluster.resourcePackage.data.operations[this.currentOperation].steps[this.currentStep].name;
      path += "/" + item.time;
      console.log(path)
      this.pangeeClusterApi.get(path).then(resp => {
        let status = resp.data;
        this.$refs.statusDialog.show(status);
      }).catch(e => {
        if (e.status == 500) {
          this.$message.error(e.response.data.message);
        }
        console.log(e);
      })
    }
  }
}
</script>

<style lang="scss" scoped>
.operation_step_history {

  .duration {
    text-align: right;
    width: 100%;
  }

  .operation {
    display: flex;
    flex-wrap: wrap;
    width: 100%;
    gap: 5px;
  }
}
</style>