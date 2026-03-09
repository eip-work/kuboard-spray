<template>
  <div style="flex-grow: 1; text-align: right;">
    <el-button v-if="loading" :loading="true" type="primary" plain> Loading </el-button>
    <el-tag type="danger" effect="dark" size="" v-else-if="error">
      {{ error.response.status || error.status }}
    </el-tag>
    <span v-else>
      <el-tag :type="computedStatus.success == computedStatus.total ? 'success' : 'danger'">
        {{ computedStatus.success }}
      </el-tag>
      /
      <el-tag style="margin-right: 5px;">
        {{ computedStatus.total }}
      </el-tag>
      <el-button @click="$refs.dialog.show(status)" :disabled="loading || status == undefined" link
        icon="el-icon-pointer">查看状态明细</el-button>
    </span>
    <el-button icon="el-icon-refresh" style="margin-left: 5px" @click="() => { checkStatus(); $emit('refresh') }"
      :loading="loading"></el-button>
    <div style="text-align: left;">
      <OperationStepStatusDialog ref="dialog"></OperationStepStatusDialog>
    </div>
  </div>
</template>

<script>
import OperationStepStatusDialog from './OperationStepStatusDialog.vue';

export default {
  props: {
    cluster: { type: Object, required: true },
    currentOperation: { type: Number, required: true },
    currentStep: { type: Number, required: true },
  },
  data() {
    return {
      status: undefined,
      error: undefined,
      loading: false,
      axiosController: undefined
    }
  },
  computed: {
    computedStatus() {
      if (this.status === undefined) {
        return { loading: true };
      }
      let total = 0;
      let success = 0;
      for (let i in this.status.data.result) {
        let node = this.status.data.result[i];
        for (let j in node) {
          let action = node[j];
            if (!action.skipped) {
            total++;
            if (action.stdout == "0") {
              success++;
            }
          }
        }
      }
      return {
        loading: false,
        success, total
      };
    },
  },
  components: { OperationStepStatusDialog },
  mounted() {
    this.checkStatus()
  },
  watch: {
    currentOperation: function () {
      this.checkStatus()
    },
    currentStep: function () {
      this.checkStatus()
    }
  },
  methods: {
    async checkStatus() {
      if (this.axiosController != undefined) {
        try {
          this.axiosController.abort();
        } catch (e) {
          console.log(e)
        }
      }
      this.loading = true;
      this.axiosController = new AbortController();
      this.status = undefined;
      this.error = undefined;
      let path = "/clusters/" + this.cluster.name + "/operation/" + this.cluster.resourcePackage.data.operations[this.currentOperation].name;
      path += "/step/" + this.cluster.resourcePackage.data.operations[this.currentOperation].steps[this.currentStep].name;
      this.pangeeClusterApi.get(path, {
        signal: this.axiosController.signal
      }).then(resp => {
        this.status = resp.data;
        this.loading = false;
      }).catch(e => {
        if (e.code != "ERR_CANCELED") {
          console.log(e);
          this.error = e;
          this.loading = false;
        }
      })
    }
  }
}
</script>

<style lang="scss" scoped></style>