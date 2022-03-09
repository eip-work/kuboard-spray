<i18n>
en: 

zh:

</i18n>

<template>
  <div :class="klass" v-if="klass.indexOf('error') > 0 || !hideSuccess">
    <div class="title" @click="expand = !expand">
      <div style="width: 150px">{{condition.type}}</div>
      <div style="width: 60px">{{condition.status}}</div>
      <div style="flex-grow: 1">{{condition.reason}}</div>
    </div>
    <transition name="el-zoom-in-top">
      <div class="expand" v-if="expand">
        <div style="display: flex;">
          <div class="label">message:</div>
          {{condition.message}}
        </div>
        <div style="display: flex;">
          <div class="label">lastHeartbeatTime:</div>
          {{condition.lastHeartbeatTime}}
        </div>
        <div style="display: flex;">
          <div class="label">lastHeartbeatTime:</div>
          {{condition.lastTransitionTime}}
        </div>
      </div>
    </transition>
  </div>
</template>

<script>
export default {
  props: {
    condition: { type: Object, required: true },
    hideSuccess: { type: Boolean, required: false, default: false }
  },
  data() {
    return {
      expand: false
    }
  },
  computed: {
    klass () {
      let result = "condition app_text_mono"
      if (this.condition.type === 'Ready') {
        if (this.condition.status !== 'True') {
          result += ' error'
        }
      } else {
        if (this.condition.status !== 'False') {
          result += ' error'
        }
      }
      return result
    }
  },
  components: { },
  mounted () {
  },
  methods: {

  }
}
</script>

<style scoped lang="scss">
.condition {
  margin: 5px 0;
  font-size: 13px;
  line-height: 18px;
  .title {
    padding: 5px 10px;
    background-color: var(--el-color-success-light-9);
    cursor: pointer;
    display: flex;
  }
  .expand {
    border: solid 1px var(--el-color-success-light);
    border-top: none;
    padding: 5px 10px;
    .label {
      width: 150px;
      color: var(--el-text-color-secondary);
    }
  }
}
.condition.error {
  .title {
    background-color: var(--el-color-danger);
    color: white;
  }
  .expand {
    border-color: var(--el-color-danger);
  }
}
</style>
