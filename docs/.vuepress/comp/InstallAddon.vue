<template>
  <div>
    本文档描述了可选组件 {{ addonName }} 的安装参数，所有可选组件的安装过程相同，请参考文档：<router-link to="/guide/addons/install_addon.html">安装可选组件</router-link>
  </div>
</template>

<script>
export default {
  props: {
    addonName: { type: String, required: false, default: 'Demo' },
  },
  data() {
    return {

    }
  },
  computed: {
  },
  components: { },
  mounted () {
  },
  methods: {

  }
}
</script>

<style scoped lang="css">

</style>
