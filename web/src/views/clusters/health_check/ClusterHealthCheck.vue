<i18n>
en:

zh:

</i18n>

<template>
  <div>
    <el-scrollbar height="calc(100vh - 220px)">
      <ConnectivityCheck :cluster="cluster"></ConnectivityCheck>
      
    </el-scrollbar>
  </div>
</template>

<script>
import ConnectivityCheck from './ConnectivityCheck.vue'

export default {
  props: {
    cluster: { type: Object, required: true },
  },
  data() {
    return {

    }
  },
  computed: {
  },
  components: { ConnectivityCheck },
  mounted () {

  },
  methods: {

  }
}
</script>

<style scoped lang="css">

</style>
