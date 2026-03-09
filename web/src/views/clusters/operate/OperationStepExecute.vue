<template>
  <div>
    <ClusterTask :cluster="cluster" :operation="cluster.resourcePackage.data.operations[currentOperation].name"
      :step="cluster.resourcePackage.data.operations[currentOperation].steps[currentStep].name"
      :title="cluster.resourcePackage.data.operations[currentOperation].steps[currentStep].title[locale]"
      :populateRequest="populateRequest" @refresh="$emit('refresh')" :width="600">
      <!-- <el-form-item :label="cluster.resourcePackage.data.operations[currentOperation].steps[currentStep].title[locale]">
        <div>参数</div>
      </el-form-item> -->
    </ClusterTask>
  </div>
</template>

<script>
import ClusterTask from '@/views/common/task/ClusterTask.vue';

export default {
  props: {
    cluster: { type: Object, required: true },
    currentOperation: { type: Number, required: true },
    currentStep: { type: Number, required: true },

  },
  data() {
    return {

    }
  },
  computed: {
  },
  components: { ClusterTask },
  emits: ['refresh'],
  mounted() {
  },
  methods: {
    async populateRequest(form) {
      console.log(form)
      let req = {}

      return req
    }
  }
}
</script>