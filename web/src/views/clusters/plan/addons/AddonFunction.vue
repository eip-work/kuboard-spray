<template>
  <div class="addon_function">
    {{ message }}
    <slot></slot>
  </div>
</template>

<script>
export default {
  props: {
    message: { type: String, required: false }
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

<style scoped lang="scss">
.addon_function {
  font-size: 13px;
  padding: 2px 10px;
  // border-bottom: var(--el-border-color-base) 1px solid;
  // border-radius: var(--el-border-radius-base);
  margin-bottom: 10px;
  background-color: var(--el-background-color-base);
  color: var(--el-color-info);
  margin: -10px -15px 10px -15px;
}
</style>
