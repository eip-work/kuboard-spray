<i18n locale="en" lang="yaml">
  hi: kuboard
</i18n>

<i18n locale="cn" lang="yaml">
  hi: kuboard
</i18n>

<template>
  <div class="command_wrapper" @click="send">
    <div class="name">{{command.name}}</div>
    <div class="app_text_mono command nowrap">{{command.cmd}}</div>
  </div>
</template>

<script>
export default {
  props: {
    command: { type: Object, required: true },
    socket: { type: Object, required: true },
    visible: { type: Boolean, required: true },
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
    send() {
      this.socket.send('0' + this.command.cmd + '\n')
      this.$emit('update:visible', false)
    }
  }
}
</script>

<style scoped lang="css">
.command_wrapper {
  display: flex;
  margin: 0 0 10px 0;
  padding: 5px 10px;
  cursor: pointer;
  background-color: var(--el-color-primary-light-9);
  .name {
    font-weight: bold;
    min-width: 200px;
  }
  .command {
    flex-grow: 1;
    margin-left: 20px;
    color: var(--el-text-color-secondary);
  }
}
</style>
